package skale

import (
	"context"
	"github.com/figment-networks/skale-indexer/scraper/structs"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGetEarnedFeeAmountOf(t *testing.T) {
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

			ccs := m.GetContractsByNames([]string{"distributor"})
			c, ok := ccs[common.HexToAddress("0x2a42Ccca55FdE8a9CA2D7f3C66fcddE99B4baB90")]
			if !ok {
				t.Fail()
			}
			bc := bind.NewBoundContract(c.Addr, c.Abi, client, nil, nil)

			e, em, err := GetEarnedFeeAmountOf(context.Background(), bc, 0, tt.args.validatorID)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEarnedFeeAmountOf() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("earned %+v, endmonth %+v", e, em)
		})
	}
}
