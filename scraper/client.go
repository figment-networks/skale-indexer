package scraper

import (
	"context"
	"fmt"
	"math/big"
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
}

type EthereumAPI struct {
	log *zap.Logger

	transport transport.EthereumTransport
	AM        ActionManager
}

func NewEthereumAPI(log *zap.Logger, transport transport.EthereumTransport, am ActionManager) *EthereumAPI {
	return &EthereumAPI{
		log:       log,
		transport: transport,
		AM:        am,
	}
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
		return nil
	}

	// TODO(lukanus): Make it configurable
	input := make(chan ProcInput, workerCount)
	output := make(chan ProcOutput, workerCount)

	go populateToWokers(logs, input)

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

type ProcInput struct {
	Order int
	Log   types.Log
}

type ProcOutput struct {
	InID  int
	CE    structs.ContractEvent
	Error error
}

func populateToWokers(logs []types.Log, populateCh chan ProcInput) {
	for i, l := range logs {
		populateCh <- ProcInput{i, l}
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
