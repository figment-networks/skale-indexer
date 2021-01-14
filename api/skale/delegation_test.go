package skale

import (
	"context"
	"math/big"
	"testing"

	"github.com/figment-networks/skale-indexer/scraper/structs"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGetValidatorDelegations(t *testing.T) {
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

			ccs := m.GetContractsByNames([]string{"delegation_controller"})
			c, ok := ccs[common.HexToAddress("0x06dD71dAb27C1A3e0B172d53735f00Bf1a66Eb79")]
			if !ok {
				t.Fail()
			}
			bc := bind.NewBoundContract(c.Addr, c.Abi, client, nil, nil)

			gotDelegations, err := GetValidatorDelegations(context.Background(), bc, 0, tt.args.validatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetValidatorDelegations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			totalStake := big.NewInt(0)
			for _, d := range gotDelegations {
				if d.State == structs.DelegationStateUNDELEGATION_REQUESTED || d.State == structs.DelegationStateDELEGATED {
					totalStake.Add(totalStake, d.Amount)
				}
			}
			t.Logf("total stake = %+v", totalStake.String())

		})
	}
}
