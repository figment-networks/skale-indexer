package client

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/figment-networks/skale-indexer/client/actions"
	"github.com/figment-networks/skale-indexer/client/structures"
	"github.com/figment-networks/skale-indexer/client/transport"
	"github.com/figment-networks/skale-indexer/client/transport/eth/contract"
	"github.com/figment-networks/skale-indexer/cmd/skale-indexer/logger"
	"go.uber.org/zap"
)

const workerCount = 5

type EthereumClient interface {
}

type EthereumAPI struct {
	transport transport.EthereumTransport
	log       *zap.Logger
	CM        *contract.Manager
	AM        *actions.Manager
}

/*
	m := NewManager()
	if err := m.LoadContractsFromDir("./testFiles"); err != nil {
		t.Error(err)
	}
*/

func (eAPI *EthereumAPI) ParseLogs(ctx context.Context, from, to big.Int) error {

	err := eAPI.transport.Dial(ctx)
	if err != nil {
		return fmt.Errorf("error dialing ethereum in ParseLogs: %w", err)
	}
	defer eAPI.transport.Close(ctx)

	ccs := eAPI.CM.GetContractsByNames([]string{"delegation_controller", // get it from somewhere else
		"validator_service",
		"nodes",
		"distributor",
		"punisher",
		"skale_manager",
		"bounty",
		"bounty_v2"})
	var addr []common.Address
	for k := range ccs {
		addr = append(addr, k)
	}
	logs, err := eAPI.transport.GetLogs(ctx, from, to)
	if err != nil {
		return fmt.Errorf("error in GetLogs request: %w", err)
	}

	eAPI.log.Debug("[EthTransport] GetLogs  ", zap.Int("len", len(logs)), zap.Any("request", logs))

	// TODO(lukanus): Make it configurable
	input := make(chan ProcInput, workerCount)
	output := make(chan ProcOutput, workerCount)

	go populateToWokers(logs, input)

	for i := 0; i < workerCount; i++ {
		go processLogAsync(ctx, eAPI.log, ccs, input, output)
	}

	processed := make(map[uint64][]ProcOutput, len(logs))

	// Process it first, but calculate
OutputLoop:
	for {
		select {
		case <-ctx.Done():
			break OutputLoop
		case o := <-output:
			logger.Debug("Process contract Event", zap.Any("ContractEvent", o.CE))
			eAPI.AM.StoreEvent(ctx, o.CE)
			p, ok := processed[o.CE.Height]
			if !ok {
				p = []ProcOutput{}
			}
			p = append(p, o)
			processed[o.CE.Height] = p
		}
	}

	return nil
}

type ProcInput struct {
	InID int
	Log  types.Log
}

type ProcOutput struct {
	InID  int
	CE    structures.ContractEvent
	Error error
}

//caller ContractCaller // to be changed to interface in future

func populateToWokers(logs []types.Log, populateCh chan ProcInput) {
	for i, l := range logs {
		populateCh <- ProcInput{i, l}
	}

	close(populateCh)
}

func processLogAsync(ctx context.Context, logger *zap.Logger, t transport.EthereumTransport, ccs map[common.Address]contract.ContractsContents, in chan ProcInput, out chan ProcOutput) {
	for {
		select {
		case <-ctx.Done():
			return
		case inp, ok := <-in:
			if !ok {
				return
			}
			ce, err := processLog(logger, inp.Log, ccs)

			c, ok := ccs[inp.Log.Address]
			bc := bind.NewBoundContract(c.Addr, c.Abi, t, nil, nil)

			eAPI.AM.AfterEventLog(ctx, bc, ce)

			out <- ProcOutput{inp.InID, ce, err}
		}
	}
}

func processLog(logger *zap.Logger, l types.Log, ccs map[common.Address]contract.ContractsContents) (ce structures.ContractEvent, err error) {
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

	return structures.ContractEvent{
		ContractName: c.Name,
		Type:         event.Name,
		Address:      c.Addr,
		Height:       l.BlockNumber,
		TxHash:       l.TxHash,
		Params:       mapped,
	}, nil
}
