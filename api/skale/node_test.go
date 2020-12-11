package skale

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/figment-networks/skale-indexer/api/structs"
)

func TestGetNode(t *testing.T) {
	type args struct {
		nodeID *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantDelegations []structs.Delegation
		wantErr         bool
	}{
		{
			name: "test happypath",
			args: args{
				nodeID: new(big.Int).SetInt64(39),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := ethclient.DialContext(ctx, "http://localhost:8545")
			if err != nil {
				t.Error(err)
				return
			}
			defer client.Close()

			m := NewManager()
			if err := m.LoadContractsFromDir("./testFiles"); err != nil {
				t.Error(err)
			}

			ccs := m.GetContractsByNames([]string{"nodes"})
			c, ok := ccs[common.HexToAddress("0xD489665414D051336CE2F2C6e4184De0409e40ba")]
			if !ok {
				t.Fail()
			}
			bc := bind.NewBoundContract(c.Addr, c.Abi, client, nil, nil)

			e, err := GetNode(context.Background(), bc, 0, tt.args.nodeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("earned %+v ", e)
		})
	}
}

func TestGetValidatorNodes(t *testing.T) {
	type args struct {
		validatorID *big.Int
	}
	tests := []struct {
		name            string
		args            args
		wantDelegations []structs.Delegation
		wantErr         bool
	}{
		{
			name: "test happypath",
			args: args{
				validatorID: new(big.Int).SetInt64(15),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			client, err := ethclient.DialContext(ctx, "http://localhost:8545")
			if err != nil {
				t.Error(err)
				return
			}
			defer client.Close()

			m := NewManager()
			if err := m.LoadContractsFromDir("./testFiles"); err != nil {
				t.Error(err)
			}

			ccs := m.GetContractsByNames([]string{"nodes"})
			c, ok := ccs[common.HexToAddress("0xD489665414D051336CE2F2C6e4184De0409e40ba")]
			if !ok {
				t.Fail()
			}
			bc := bind.NewBoundContract(c.Addr, c.Abi, client, nil, nil)

			v, err := GetValidatorNodes(context.Background(), bc, 0, tt.args.validatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("%+v", v)
		})
	}
}
