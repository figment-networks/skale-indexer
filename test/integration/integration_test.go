package integration

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/api/skale"
	"github.com/figment-networks/skale-indexer/client/actions"
	"github.com/figment-networks/skale-indexer/scraper"
	clientStructures "github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"github.com/stretchr/testify/require"

	storeMocks "github.com/figment-networks/skale-indexer/store/mocks"
	"github.com/golang/mock/gomock"
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
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()
			mockDB := storeMocks.NewMockDataStore(mockCtrl)

			am := actions.NewManager(caller, mockDB, tr, cm, zl)
			eAPI := scraper.NewEthereumAPI(zl, tr, am)

			ccs := cm.GetContractsByNames(am.GetImplementedContractNames())
			if err := eAPI.ParseLogs(ctx, ccs, tt.args.from, tt.args.to); err != nil {
				t.Error(err)
				return
			}
		})
	}
}

func TestCallWithBlockNumber(t *testing.T) {
	type args struct {
		address  string
		contract common.Address
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
				address:  "http://localhost:8545",
				contract: common.HexToAddress("0x06dD71dAb27C1A3e0B172d53735f00Bf1a66Eb79"),
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

			cm := contract.NewManager()
			if err := cm.LoadContractsFromDir("./testFiles"); err != nil {
				t.Error(err)
				return
			}
			caller := &skale.Caller{}
			ac, _ := cm.GetContract(tt.args.contract)
			bc := tr.GetBoundContractCaller(ctx, ac.Addr, ac.Abi)

			a, err := caller.GetDelegation(ctx, bc, uint64(10814408), big.NewInt(2316))
			require.NoError(t, err)
			require.Equal(t, a.DelegationID, big.NewInt(2316))
			ds, err := caller.GetDelegationState(ctx, bc, uint64(10814408), big.NewInt(2316))
			require.NoError(t, err)
			require.Equal(t, ds, clientStructures.DelegationStateCOMPLETED)

		})
	}
}
