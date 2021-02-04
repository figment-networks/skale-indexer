package scraper

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"go.uber.org/zap"
)

const workerCount = 5

type ActionManager interface {
	GetImplementedContractNames() []string
	GetBlockHeader(ctx context.Context, height *big.Int) (h *types.Header, err error)
	AfterEventLog(ctx context.Context, c contract.ContractsContents, ce structs.ContractEvent) error
	SyncForBeginningOfEpoch(ctx context.Context, c contract.ContractsContents, currentBlock uint64, blockTime time.Time) error
}

type EthereumAPI struct {
	log *zap.Logger

	transport                 transport.EthereumTransport
	AM                        ActionManager
	rangeBlockCache           *rangeBlockCache
	lowerThresholdForBackward uint64
}

func NewEthereumAPI(log *zap.Logger, transport transport.EthereumTransport, am ActionManager, lowerThresholdForBackward uint64) *EthereumAPI {
	return &EthereumAPI{
		log:                       log,
		transport:                 transport,
		AM:                        am,
		rangeBlockCache:           newLastBlockCache(),
		lowerThresholdForBackward: lowerThresholdForBackward,
	}
}

type rangeBlockCache struct {
	mu              sync.Mutex
	rangeBlockCache map[rangeInfo]uint64
}

type rangeInfo struct {
	from uint64
	to   uint64
}

func newLastBlockCache() *rangeBlockCache {
	return &rangeBlockCache{rangeBlockCache: make(map[rangeInfo]uint64)}
}

func (rbc *rangeBlockCache) add(r rangeInfo) {
	defer rbc.mu.Unlock()

	rbc.mu.Lock()
	rbc.rangeBlockCache[r] = r.from + r.to
}

func (rbc *rangeBlockCache) isInCheckedRange(bgnBlock uint64) bool {
	defer rbc.mu.Unlock()

	rbc.mu.Lock()
	for key, _ := range rbc.rangeBlockCache {
		if bgnBlock >= key.from && bgnBlock <= key.to {
			return true
		}
	}
	return false
}

func (eAPI *EthereumAPI) setBackwards(ctx context.Context, from, to big.Int, addr []common.Address) (block uint64, err error) {
	var logsLength uint64
	r := big.NewInt(50)
	for logsLength == 0 && from.Uint64() >= eAPI.lowerThresholdForBackward {
		from.Sub(&from, r)
		logsBackwards, err := eAPI.transport.GetLogs(ctx, from, to, addr)
		if err != nil {
			return block, errors.New("error on getting logs for backwards")
		}
		logsLength = uint64(len(logsBackwards))
		if logsLength > 0 {
			toRange := logsBackwards[len(logsBackwards)-1].BlockNumber
			return toRange, err
		}
	}
	return eAPI.lowerThresholdForBackward, nil
}

func (eAPI *EthereumAPI) ParseLogs(ctx context.Context, ccs map[common.Address]contract.ContractsContents, from, to big.Int) error {

	addr := make([]common.Address, len(ccs))
	var i int
	for k := range ccs {
		addr[i] = k
		i++
	}
	logs, err := eAPI.transport.GetLogs(ctx, from, to, addr)
	if err != nil {
		return fmt.Errorf("error in GetLogs request: %w", err)
	}

	eAPI.log.Debug("[EthTransport] GetLogs  ", zap.Int("len", len(logs)), zap.Any("request", logs))

	if len(logs) == 0 {
		eAPI.rangeBlockCache.add(rangeInfo{from: from.Uint64(), to: to.Uint64()})
		return nil
	}

	var lastLoggedBlockTime time.Time
	var lastLoggedBlockHeader *types.Header
	if eAPI.rangeBlockCache.isInCheckedRange(logs[0].BlockNumber) {
		lastLoggedBlockHeader, _ = eAPI.AM.GetBlockHeader(ctx, new(big.Int).SetUint64(logs[0].BlockNumber))
		lastLoggedBlockTime = time.Unix(int64(lastLoggedBlockHeader.Time), 0)
	} else {
		lastLoggedBlock, err := eAPI.setBackwards(ctx, *new(big.Int).SetUint64(logs[0].BlockNumber - 1), *new(big.Int).SetUint64(logs[0].BlockNumber - 1), addr)
		if err != nil {
			return err
		}
		lastLoggedBlockHeader, _ = eAPI.AM.GetBlockHeader(ctx, new(big.Int).SetUint64(lastLoggedBlock))
		lastLoggedBlockTime = time.Unix(int64(lastLoggedBlockHeader.Time), 0)
	}

	eAPI.rangeBlockCache.add(rangeInfo{from: from.Uint64(), to: to.Uint64()})

	// TODO(lukanus): Make it configurable
	input := make(chan ProcInput, workerCount)
	output := make(chan ProcOutput, workerCount)

	go populateToWorkers(logs, input, lastLoggedBlockTime)

	for i := 0; i < workerCount; i++ {
		go eAPI.processLogAsync(ctx, ccs, input, output)
	}

	processed := make(map[uint64][]ProcOutput, len(logs))

	var gotResponses int
	// Process it first, but calculate
OutputLoop:
	for {
		select {
		case <-ctx.Done():
			break OutputLoop
		case o := <-output:
			gotResponses++
			if o.Error != nil {
				eAPI.log.Error("Error", zap.Error(o.Error))
				continue
			}
			eAPI.log.Debug("Process contract Event", zap.Any("ContractEvent", o.CE))
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

	return nil
}

func isInRange(crnTime, prvTime time.Time) bool {
	return (crnTime.Year() > prvTime.Year()) || (crnTime.Month() > prvTime.Month())
}

type ProcInput struct {
	Order               int
	Log                 types.Log
	PreviousHeight      uint64
	lastLoggedBlockTime time.Time
}

type ProcOutput struct {
	InID  int
	CE    structs.ContractEvent
	Error error
}

func populateToWorkers(logs []types.Log, populateCh chan ProcInput, lastLoggedBlockTime time.Time) {
	var prevHeight uint64
	for i, l := range logs {
		populateCh <- ProcInput{i, l, prevHeight, lastLoggedBlockTime}
		prevHeight = l.BlockNumber
	}

	close(populateCh)
}

func (eAPI *EthereumAPI) processLogAsync(ctx context.Context, ccs map[common.Address]contract.ContractsContents, in chan ProcInput, out chan ProcOutput) {
	defer eAPI.log.Sync()

	for {
		select {
		case <-ctx.Done():
			return
		case inp, ok := <-in:
			if !ok {
				return
			}
			h, err := eAPI.AM.GetBlockHeader(ctx, new(big.Int).SetUint64(inp.Log.BlockNumber))
			if err != nil {
				out <- ProcOutput{Error: err}
				continue
			}

			ce, err := processLog(eAPI.log, inp.Log, h, ccs)
			if err != nil {
				out <- ProcOutput{Error: err}
				continue
			}
			c, ok := ccs[inp.Log.Address]
			err = eAPI.AM.AfterEventLog(ctx, c, ce)
			if err != nil {
				out <- ProcOutput{Error: err}
				continue
			}

			hTime := time.Unix(int64(h.Time), 0)

			// running sync function for the first log of the new epoch
			if inp.PreviousHeight != 0 && inp.PreviousHeight < inp.Log.BlockNumber {
				prvHeader, _ := eAPI.AM.GetBlockHeader(ctx, new(big.Int).SetUint64(inp.PreviousHeight))
				prvTime := time.Unix(int64(prvHeader.Time), 0)
				if isInRange(hTime, prvTime) {
					err = eAPI.AM.SyncForBeginningOfEpoch(ctx, c, inp.Log.BlockNumber, hTime)
					if err != nil {
						eAPI.log.Error("error occurred on synchronization ", zap.Error(err))
						out <- ProcOutput{Error: err}
						continue
					}
				}
			} else if inp.PreviousHeight == 0 {
				// this is the case for the first block of the current logs' round
				if isInRange(hTime, inp.lastLoggedBlockTime) {
					err = eAPI.AM.SyncForBeginningOfEpoch(ctx, c, inp.Log.BlockNumber, hTime)
					if err != nil {
						eAPI.log.Error("error occurred on synchronization ", zap.Error(err))
						out <- ProcOutput{Error: err}
						continue
					}
				}
			}

			out <- ProcOutput{inp.Order, ce, err}
		}
	}
}

func processLog(logger *zap.Logger, l types.Log, h *types.Header, ccs map[common.Address]contract.ContractsContents) (ce structs.ContractEvent, err error) {
	c, ok := ccs[l.Address]
	if !ok {
		logger.Error("[EthTransport] GetLogs contract not found ", zap.String("txHash", l.TxHash.String()), zap.String("address", l.Address.String()))
		return ce, fmt.Errorf("error in GetLogs, there is no such contract as %s ", l.Address.String())
	}
	if len(l.Topics) == 0 {
		logger.Error("[EthTransport] GetLogs list has empty topic list", zap.String("txHash", l.TxHash.String()), zap.String("address", l.Address.String()))
		return ce, fmt.Errorf("getLogs list has empty topic list")
	}
	logger.Debug("[EthTransport] GetLogs got contract", zap.String("name", c.Name), zap.Uint64("request", l.BlockNumber))
	event, err := c.Abi.EventByID(l.Topics[0])
	if err != nil {
		logger.Error("[EthTransport] GetLogs abi has no such event",
			zap.String("txHash", l.TxHash.String()),
			zap.String("address", l.Address.String()),
			zap.String("name", c.Name),
			zap.Uint64("request", l.BlockNumber),
		)
		return ce, fmt.Errorf("getLogs list has empty topic list %w", err)
	}
	mapped := make(map[string]interface{}, len(event.Inputs))
	err = event.Inputs.UnpackIntoMap(mapped, l.Data)
	if err != nil {
		return ce, fmt.Errorf("error unpacking into map %w", err)
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
