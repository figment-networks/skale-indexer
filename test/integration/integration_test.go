package integration

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/figment-networks/skale-indexer/api/actions"
	"github.com/figment-networks/skale-indexer/api/skale"
	"github.com/figment-networks/skale-indexer/client"
	clientStructures "github.com/figment-networks/skale-indexer/client/structs"
	"github.com/figment-networks/skale-indexer/client/transport/eth"
	"github.com/figment-networks/skale-indexer/client/transport/eth/contract"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestGetLogs(t *testing.T) {
	type args struct {
		address string
		from    big.Int
		to      big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantDelegations []clientStructures.Delegation
		wantErr         bool
	}{
		{
			name: "test1",
			args: args{
				address: "http://localhost:8545",
				from:    *big.NewInt(10880000),
				to:      *big.NewInt(10890000),
			},
		},
		{
			name: "test2",
			args: args{
				address: "http://localhost:8545",
				from:    *big.NewInt(10910000),
				to:      *big.NewInt(10920000),
			},
		},
		{
			name: "test3",
			args: args{
				address: "http://localhost:8545",
				from:    *big.NewInt(11010000),
				to:      *big.NewInt(11020000),
			},
		},
		{
			name: "test4",
			args: args{
				address: "http://localhost:8545",
				from:    *big.NewInt(11110000),
				to:      *big.NewInt(11120000),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			tr := eth.NewEthTransport(tt.args.address)
			if err := tr.Dial(ctx); err != nil {
				t.Errorf("Error dialing %s : %w", tt.args.address, err)
				return
			}
			defer tr.Close(ctx)

			zl := zaptest.NewLogger(t)

			cm := contract.NewManager()
			if err := cm.LoadContractsFromDir("./testFiles"); err != nil {
				t.Error(err)
				return
			}
			caller := &skale.Caller{}
			slm := &StoreLogMock{zl}
			clm := &CalculatorLogMock{zl}
			am := actions.NewManager(caller, slm, clm, tr, cm)
			eAPI := client.NewEthereumAPI(zl, tr, am)

			ccs := cm.GetContractsByNames(am.GetImplementedEventsNames())
			if err := eAPI.ParseLogs(ctx, ccs, tt.args.from, tt.args.to); err != nil {
				t.Error(err)
				return
			}
		})
	}
}

type StoreLogMock struct {
	logger *zap.Logger
}

func (slm *StoreLogMock) StoreEvent(ctx context.Context, v clientStructures.ContractEvent) error {
	slm.logger.Info("Storing event: ", zap.Any("event", v))
	slm.logger.Sync()
	return nil
}

func (slm *StoreLogMock) StoreValidator(ctx context.Context, height uint64, t time.Time, v clientStructures.Validator) error {
	slm.logger.Info("Storing validator: ", zap.Any("validator", v))
	slm.logger.Sync()
	return nil
}

func (slm *StoreLogMock) StoreDelegation(ctx context.Context, height uint64, t time.Time, d clientStructures.Delegation) error {
	slm.logger.Info("Storing delegation: ", zap.Any("delegation", d))
	slm.logger.Sync()
	return nil
}

func (slm *StoreLogMock) StoreNode(ctx context.Context, height uint64, t time.Time, v clientStructures.Node) error {
	slm.logger.Info("Storing node: ", zap.Any("node", v))
	slm.logger.Sync()
	return nil
}

func (slm *StoreLogMock) StoreValidatorNodes(ctx context.Context, height uint64, t time.Time, nodes []clientStructures.Node) error {
	slm.logger.Info("Storing validator nodes: ", zap.Any("nodes", nodes))
	slm.logger.Sync()
	return nil
}

type CalculatorLogMock struct {
	logger *zap.Logger
}

func (clm *CalculatorLogMock) ValidatorParams(ctx context.Context, height uint64, vID *big.Int) error {
	clm.logger.Info("Calculating validator params: ", zap.Any("validatorID", vID))
	clm.logger.Sync()
	return nil
}

func (clm *CalculatorLogMock) DelegationParams(ctx context.Context, height uint64, dID *big.Int) error {
	clm.logger.Info("Calculating delegation params: ", zap.Any("delegationID", dID))
	clm.logger.Sync()
	return nil
}
