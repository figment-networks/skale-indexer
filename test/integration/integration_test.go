package integration

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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
			eAPI := scraper.NewEthereumAPI(zl, tr, types.Header{Number: big.NewInt(1234), Time: uint64(1234)}, am)

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

			a, err := caller.GetDelegation(ctx, bc.GetContract(), uint64(10814408), big.NewInt(2316))
			require.NoError(t, err)
			require.Equal(t, a.DelegationID, big.NewInt(2316))
			ds, err := caller.GetDelegationState(ctx, bc.GetContract(), uint64(10814408), big.NewInt(2316))
			require.NoError(t, err)
			require.Equal(t, ds, clientStructures.DelegationStateCOMPLETED)

		})
	}
}

func TestGetAndUpdateEarnedBountyAmountOf(t *testing.T) {
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
				contract: common.HexToAddress("0x2a42Ccca55FdE8a9CA2D7f3C66fcddE99B4baB90"),
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

			_, _, err := caller.GetAndUpdateEarnedBountyAmountOf(ctx, bc.GetContract(), big.NewInt(53), common.HexToAddress("0xFFEA818a2a4bF047Af42487A30290cF6F8e80dd2"), 0)
			require.NoError(t, err)

		})
	}
}

func TestGetNodeAddress(t *testing.T) {
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
				contract: common.HexToAddress("0xD489665414D051336CE2F2C6e4184De0409e40ba"),
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

			_, err := caller.GetNodeAddress(ctx, bc, 10930066, big.NewInt(1))
			require.NoError(t, err)

		})
	}
}
