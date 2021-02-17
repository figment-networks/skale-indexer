package scraper

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	lru "github.com/hashicorp/golang-lru"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"go.uber.org/zap"
)

const (
	workerCount            = 5
	backCheckSlidingWindow = 100
	cacheSize              = 1000
)

type ActionManager interface {
	GetImplementedContractNames() []string
	GetBlockHeader(ctx context.Context, height big.Int) (h *types.Header, err error)
	AfterEventLog(ctx context.Context, c contract.ContractsContents, ce structs.ContractEvent) error
	SyncForBeginningOfEpoch(ctx context.Context, contractVersion string, currentBlock uint64, blockTime time.Time) error
}

type EthereumAPI struct {
	log                   *zap.Logger
	transport             transport.EthereumTransport
	AM                    ActionManager
	blockLRU              *lru.Cache
	smallestPossibleBlock types.Header

	slock sync.Mutex
}

func NewEthereumAPI(log *zap.Logger, transport transport.EthereumTransport, spb types.Header, am ActionManager) *EthereumAPI {
	cache, _ := lru.New(cacheSize)
	return &EthereumAPI{
		log:                   log,
		transport:             transport,
		AM:                    am,
		smallestPossibleBlock: spb,
		blockLRU:              cache,
	}
}

func (eAPI *EthereumAPI) getLastBlockTimeBefore(ctx context.Context, taskID string, fromBlockID uint64, window uint64, addr []common.Address) (blockTime time.Time, err error) {

	f, t := fromBlockID-window, fromBlockID
	for {
		eAPI.log.Debug("Running back check ", zap.Uint64("from", f), zap.Uint64("to", t))

		if f < eAPI.smallestPossibleBlock.Number.Uint64() {
			blockTime = time.Unix(int64(eAPI.smallestPossibleBlock.Time), 0)
			eAPI.blockLRU.Add(taskID, eAPI.smallestPossibleBlock)
			return
		}
		eAPI.log.Debug("Running GetLogs", zap.Uint64("from", f), zap.Uint64("to", t))
		logsBackwards, err := eAPI.transport.GetLogs(ctx, *new(big.Int).SetUint64(f), *new(big.Int).SetUint64(t), addr)
		if err != nil {
			return blockTime, fmt.Errorf("error on getting logs for last block before :%w", err)
		}

		if len(logsBackwards) > 0 {
			height := logsBackwards[len(logsBackwards)-1].BlockNumber
			lastLoggedBlockHeader, err := eAPI.AM.GetBlockHeader(ctx, *new(big.Int).SetUint64(height))
			if err != nil {
				return blockTime, fmt.Errorf("error on getting block header for last block before :%w", err)
			}
			blockTime = time.Unix(int64(lastLoggedBlockHeader.Time), 0)
			eAPI.blockLRU.Add(taskID, *lastLoggedBlockHeader)
			return blockTime, nil
		}
		f, t = f-window, t-window
	}
}

func (eAPI *EthereumAPI) GetLatestBlockHeight(ctx context.Context) (uint64, error) {
	return eAPI.transport.GetLatestBlockHeight(ctx)
}

func (eAPI *EthereumAPI) ParseLogs(ctx context.Context, ccs *contract.Contracts, taskID string, from, to big.Int) (err error) {
	defer eAPI.log.Sync()

	addr := ccs.GetAddresses()
	logs, err := eAPI.transport.GetLogs(ctx, from, to, addr)
	if err != nil {
		return fmt.Errorf("error in GetLogs request: %w", err)
	}

	eAPI.log.Debug("[EthTransport] GetLogs ", zap.Int("len", len(logs)), zap.String("taskID", taskID), zap.Uint64("from", from.Uint64()), zap.Uint64("to", to.Uint64()))

	if len(logs) == 0 { // spot tx block crossing month
		var lastLoggedBlockTime types.Header
		llbt, ok := eAPI.blockLRU.Get(taskID)
		if ok {
			lastLoggedBlockTime = llbt.(types.Header)
			if lastLoggedBlockTime.Time == 0 {
				lbtB, err := eAPI.getLastBlockTimeBefore(ctx, taskID, from.Uint64(), backCheckSlidingWindow, addr)
				if err != nil {
					return err
				}
				lastLoggedBlockTime.Time = uint64(lbtB.Unix())

			}
		}
		h, err := eAPI.AM.GetBlockHeader(ctx, to)
		if err != nil {
			return err
		}

		hTime := time.Unix(int64(h.Time), 0)
		a := time.Unix(int64(lastLoggedBlockTime.Time), 0)
		if isInRange(a, hTime) {
			if err = eAPI.AM.SyncForBeginningOfEpoch(ctx, "1.7.2", h.Number.Uint64(), hTime); err != nil { // latest version?
				return err
			}
			eAPI.blockLRU.Add(taskID, *h)
			return nil
		}

		eAPI.blockLRU.Add(taskID, *h)
		return nil
	}

	var lastLoggedBlockTime types.Header
	llbt, _ := eAPI.blockLRU.Get(taskID)
	if llbt != nil {
		lastLoggedBlockTime = llbt.(types.Header)
	}

	if lastLoggedBlockTime.Time == 0 {
		lbtB, err := eAPI.getLastBlockTimeBefore(ctx, taskID, from.Uint64(), backCheckSlidingWindow, addr)
		if err != nil {
			return err
		}
		lastLoggedBlockTime.Time = uint64(lbtB.Unix())
	}

	input := make(chan ProcInput, workerCount)
	output := make(chan ProcOutput, workerCount)
	defer close(output)

	cCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	wg := &sync.WaitGroup{}

	go eAPI.populateToWorkers(cCtx, logs, input, time.Unix(int64(lastLoggedBlockTime.Time), 0))
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go eAPI.processLogAsync(cCtx, ccs, wg, input, output)
	}

	processed := make(map[uint64][]ProcOutput, len(logs))
	var gotResponses int

OutputLoop:
	for {
		select {
		case <-ctx.Done():
			break OutputLoop
		case <-cCtx.Done():
			break OutputLoop
		case o := <-output:
			gotResponses++
			if o.Error != nil {
				err = o.Error
				cancel()
				eAPI.log.Error("[ScraperCLient] ParseLogs Error", zap.Error(o.Error))
				continue
			}
			eAPI.log.Debug("[ScraperCLient] Process contract Event", zap.Any("ContractEvent", o.CE))
			p, ok := processed[o.CE.BlockHeight]
			if !ok {
				p = []ProcOutput{}
			}
			p = append(p, o)
			processed[o.CE.BlockHeight] = p

			if gotResponses == len(logs) {
				break OutputLoop
			}
		}
	}

	if err == nil {
		lastBlock := logs[len(logs)-1]

		a := processed[lastBlock.BlockNumber]
		eAPI.blockLRU.Add(taskID, types.Header{
			Number: new(big.Int).SetUint64(lastBlock.BlockNumber),
			Time:   uint64(a[0].CE.Time.Unix()),
		})
		eAPI.log.Debug("[ScraperCLient] Processed", zap.Uint64("last_height", lastBlock.BlockNumber))
	}

	wg.Wait()

	return err
}

func isInRange(prvTime, crnTime time.Time) bool {
	return (crnTime.Year() > prvTime.Year()) || (crnTime.Month() > prvTime.Month())
}

type ProcInput struct {
	Order             int
	Log               types.Log
	Header            types.Header
	PreviousBlockTime time.Time

	Error error
}

type ProcOutput struct {
	InID  int
	CE    structs.ContractEvent
	Error error
}

func (eAPI *EthereumAPI) populateToWorkers(ctx context.Context, logs []types.Log, populateCh chan ProcInput, lastLoggedBlockTime time.Time) {

	defer close(populateCh)
	previousBlockTime := lastLoggedBlockTime
	for i, l := range logs {
		select {
		case <-ctx.Done():
			return
		default:
		}
		h, err := eAPI.AM.GetBlockHeader(ctx, *new(big.Int).SetUint64(l.BlockNumber))
		if err != nil {
			populateCh <- ProcInput{Error: err}
			break
		}

		populateCh <- ProcInput{i, l, *h, previousBlockTime, nil}
		previousBlockTime = time.Unix(int64(h.Time), 0)
	}

}

func (eAPI *EthereumAPI) processLogAsync(ctx context.Context, ccs *contract.Contracts, wg *sync.WaitGroup, in chan ProcInput, out chan ProcOutput) {
	defer eAPI.log.Sync()
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case inp, ok := <-in:
			if !ok {
				return
			}
			if inp.Error != nil {
				if eAPI.sendIfPossible(ctx, out, ProcOutput{Error: inp.Error}) {
					return
				}
				continue
			}

			ce, err := processLog(eAPI.log, inp.Log, inp.Header, ccs)
			if err != nil {
				if eAPI.sendIfPossible(ctx, out, ProcOutput{Error: err}) {
					return
				}
				continue
			}
			c, ok := ccs.GetByAddress(inp.Log.Address)
			if err = eAPI.AM.AfterEventLog(ctx, c, ce); err != nil {
				if eAPI.sendIfPossible(ctx, out, ProcOutput{Error: err}) {
					return
				}
				continue
			}

			hTime := time.Unix(int64(inp.Header.Time), 0)
			if isInRange(inp.PreviousBlockTime, hTime) {
				if err = eAPI.AM.SyncForBeginningOfEpoch(ctx, c.Version, inp.Log.BlockNumber, hTime); err != nil {
					eAPI.log.Error("error occurred on synchronization ", zap.Error(err))
					if eAPI.sendIfPossible(ctx, out, ProcOutput{Error: err}) {
						return
					}
					continue
				}
			}
			if eAPI.sendIfPossible(ctx, out, ProcOutput{inp.Order, ce, err}) {
				return
			}
		}
	}
}

func (eAPI *EthereumAPI) sendIfPossible(ctx context.Context, out chan ProcOutput, po ProcOutput) bool {
	select {
	case <-ctx.Done():
		return true
	case out <- po:
	}
	return false
}

func processLog(logger *zap.Logger, l types.Log, h types.Header, ccs *contract.Contracts) (ce structs.ContractEvent, err error) {
	c, ok := ccs.GetByAddress(l.Address)

	if !ok {
		logger.Error("[EthTransport] GetLogs contract not found ", zap.String("txHash", l.TxHash.String()), zap.String("address", l.Address.String()))
		return ce, fmt.Errorf("error in GetLogs, there is no such contract as %s", l.Address.String())
	}
	if len(l.Topics) == 0 {
		logger.Error("[EthTransport] GetLogs list has empty topic list", zap.String("txHash", l.TxHash.String()), zap.String("address", l.Address.String()))
		return ce, fmt.Errorf("getLogs list has empty topic list")
	}
	event, err := c.Abi.EventByID(l.Topics[0])
	if err != nil {
		vers, ok := ccs.GetAllVersions(l.Address)
		if ok {
			for _, v := range vers {
				event, err = v.Abi.EventByID(l.Topics[0])
				if err == nil {
					c = v
					break
				}
			}
		}
		if err != nil {
			logger.Error("[EthTransport] GetLogs abi has no such event",
				zap.String("txHash", l.TxHash.String()),
				zap.String("address", l.Address.String()),
				zap.String("name", c.Name),
				zap.Uint64("request", l.BlockNumber),
			)
			return ce, fmt.Errorf("getLogs list has empty topic list %w", err)
		}
	}

	logger.Debug("[EthTransport] GetLogs got contract", zap.String("name", c.Name), zap.Uint64("block", l.BlockNumber), zap.Time("blockTime", time.Unix(int64(h.Time), 0)))
	mapped := make(map[string]interface{}, len(event.Inputs))
	if len(l.Data) > 0 {
		err = event.Inputs.UnpackIntoMap(mapped, l.Data)
		if err != nil {
			return ce, fmt.Errorf("error unpacking into map %w", err)
		}
	}
	i := 1 // skip first topic, because it's event data
	for _, v := range event.Inputs {
		if v.Indexed == true {
			switch v.Type.String() {
			case "uint256":
				mapped[v.Name] = abi.ReadInteger(v.Type, l.Topics[i].Bytes())
			case "address":
				mapped[v.Name] = common.BytesToAddress(l.Topics[i].Bytes())
			case "bytes32":
				mapped[v.Name] = l.Topics[i].Bytes()
			}
			i++
		}
	}

	return structs.ContractEvent{
		ContractName:    c.Name,
		EventName:       event.Name,
		ContractAddress: c.Addr,
		BlockHeight:     l.BlockNumber,
		Time:            time.Unix(int64(h.Time), 0),
		TransactionHash: l.TxHash,
		Params:          mapped,
		Removed:         l.Removed,
	}, nil
}

type rangeBlockCache struct {
	c map[rangeInfo]types.Header
	l sync.RWMutex
}

type rangeInfo struct {
	from uint64
	to   uint64
}

func newLastBlockCache() rangeBlockCache {
	return rangeBlockCache{
		c: make(map[rangeInfo]types.Header),
	}
}

func (rbc *rangeBlockCache) Set(r rangeInfo, h types.Header) {
	rbc.l.Lock()
	defer rbc.l.Unlock()

	var inMap bool
	for k := range rbc.c {
		//inclusive skip
		if r.from >= k.from && r.to <= k.to {
			inMap = true
			if h.Time != 0 {
				delete(rbc.c, k)
				rbc.c[rangeInfo{from: h.Number.Uint64(), to: k.to}] = h
			}
		} else if k.from == r.to || (k.from == r.to+1) { // left sided - join
			inMap = true
			delete(rbc.c, k)
			if h.Time != 0 {
				rbc.c[rangeInfo{from: h.Number.Uint64(), to: k.to}] = h
			} else {
				rbc.c[rangeInfo{from: r.from, to: k.to}] = h
			}
		} else if k.to == r.from || (k.to == r.from-1) { // right sided - join
			inMap = true
			delete(rbc.c, k)
			if h.Time != 0 {
				rbc.c[rangeInfo{from: h.Number.Uint64(), to: r.to}] = h
			} else {
				rbc.c[rangeInfo{from: k.from, to: r.to}] = h
			}
		}
	}
	if !inMap {
		if h.Time != 0 {
			rbc.c[rangeInfo{from: h.Number.Uint64(), to: r.to}] = h
		} else {
			rbc.c[rangeInfo{from: r.from, to: r.to}] = h
		}
	}

}

func (rbc *rangeBlockCache) Get(bgnBlock uint64) (time.Time, bool) {
	rbc.l.RLock()
	defer rbc.l.RUnlock()

	for k, t := range rbc.c {
		if bgnBlock >= k.from && bgnBlock <= k.to {
			return time.Unix(int64(t.Time), 0), true
		}
	}
	return time.Time{}, false
}

func (rbc *rangeBlockCache) GetRange(start, stop uint64) (time.Time, bool) {
	rbc.l.RLock()
	defer rbc.l.RUnlock()

	for k, t := range rbc.c {
		if start >= k.from && stop <= k.to {
			return time.Unix(int64(t.Time), 0), true
		}
	}
	return time.Time{}, false
}
