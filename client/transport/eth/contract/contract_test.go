package contract

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/require"
)

func TestManager_LoadContract(t *testing.T) {
	type args struct {
		name        string
		addr        string
		version     string
		abiContents []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{name: "first", args: args{name: "name", addr: "0x06dD71dAb27C1A3e0B172d53735f00Bf1a66Eb79", abiContents: []byte(`[{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"DelegationAccepted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"DelegationProposed","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"DelegationRequestCanceledByUser","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"role","type":"bytes32"},{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"RoleGranted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"bytes32","name":"role","type":"bytes32"},{"indexed":true,"internalType":"address","name":"account","type":"address"},{"indexed":true,"internalType":"address","name":"sender","type":"address"}],"name":"RoleRevoked","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"UndelegationRequested","type":"event"},{"inputs":[],"name":"DEFAULT_ADMIN_ROLE","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"contractManager","outputs":[{"internalType":"contract ContractManager","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"}],"name":"delegations","outputs":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"delegationPeriod","type":"uint256"},{"internalType":"uint256","name":"created","type":"uint256"},{"internalType":"uint256","name":"started","type":"uint256"},{"internalType":"uint256","name":"finished","type":"uint256"},{"internalType":"string","name":"info","type":"string"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"","type":"address"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"delegationsByHolder","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"","type":"uint256"},{"internalType":"uint256","name":"","type":"uint256"}],"name":"delegationsByValidator","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"}],"name":"getRoleAdmin","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"uint256","name":"index","type":"uint256"}],"name":"getRoleMember","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"}],"name":"getRoleMemberCount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"grantRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"hasRole","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"renounceRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes32","name":"role","type":"bytes32"},{"internalType":"address","name":"account","type":"address"}],"name":"revokeRole","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"}],"name":"getAndUpdateDelegatedToValidatorNow","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"}],"name":"getAndUpdateDelegatedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"month","type":"uint256"}],"name":"getAndUpdateEffectiveDelegatedByHolderToValidator","outputs":[{"internalType":"uint256","name":"effectiveDelegated","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"delegationPeriod","type":"uint256"},{"internalType":"string","name":"info","type":"string"}],"name":"delegate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"wallet","type":"address"}],"name":"getAndUpdateLockedAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"wallet","type":"address"}],"name":"getAndUpdateForbiddenForDelegationAmount","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"cancelPendingDelegation","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"acceptPendingDelegation","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"requestUndelegation","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"confiscate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"month","type":"uint256"}],"name":"getAndUpdateEffectiveDelegatedToValidator","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"validatorId","type":"uint256"}],"name":"getAndUpdateDelegatedByHolderToValidatorNow","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"getDelegation","outputs":[{"components":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint256","name":"delegationPeriod","type":"uint256"},{"internalType":"uint256","name":"created","type":"uint256"},{"internalType":"uint256","name":"started","type":"uint256"},{"internalType":"uint256","name":"finished","type":"uint256"},{"internalType":"string","name":"info","type":"string"}],"internalType":"struct DelegationController.Delegation","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"validatorId","type":"uint256"}],"name":"getFirstDelegationMonth","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"}],"name":"getDelegationsByValidatorLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"}],"name":"getDelegationsByHolderLength","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"contractsAddress","type":"address"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"validatorId","type":"uint256"},{"internalType":"uint256","name":"month","type":"uint256"}],"name":"getAndUpdateDelegatedToValidator","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"},{"internalType":"uint256","name":"limit","type":"uint256"}],"name":"processSlashes","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"}],"name":"processAllSlashes","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"delegationId","type":"uint256"}],"name":"getState","outputs":[{"internalType":"enum DelegationController.State","name":"state","type":"uint8"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"}],"name":"getLockedInPendingDelegations","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"address","name":"holder","type":"address"}],"name":"hasUnprocessedSlashes","outputs":[{"internalType":"bool","name":"","type":"bool"}],"stateMutability":"view","type":"function"}]
`)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := NewManager()
			a := abi.ABI{}
			json.Unmarshal(tt.args.abiContents, &a)
			if err := m.LoadContract(tt.args.name, tt.args.addr, tt.args.version, a); (err != nil) != tt.wantErr {
				t.Errorf("Manager.LoadContracts() error = %v, wantErr %v", err, tt.wantErr)
			}
			addr := ConvertToAddress([]byte(tt.args.addr))
			cc, ok := m.GetContract(addr)

			require.True(t, ok)
			require.Equal(t, cc.Addr, addr)
			require.Equal(t, cc.Name, tt.args.name)

			for k, v := range cc.Abi.Events {
				t.Log("ev", k, v)
			}

			for k, v := range cc.Abi.Methods {
				t.Log("m", k, v)
			}
		})
	}
}

func TestManager_LoadContracts(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		m := NewManager()
		if err := m.LoadContractsFromDir("./testFiles"); err != nil {
			t.Error(err)
		}

		ccs := m.GetContractsByNames([]string{"delegation_controller", "validator_service"})

		require.Len(t, ccs, 2)

	})
}

func TestManager_LoadContracts1(t *testing.T) {

	t.Run("ok", func(t *testing.T) {
		ctx := context.Background()
		client, err := ethclient.DialContext(ctx, "http://localhost:8545")
		if err != nil {
			t.Error(err)
			return
		}
		defer client.Close()

		fq := ethereum.FilterQuery{
			FromBlock: big.NewInt(10887703),
			ToBlock:   big.NewInt(10899993),
		}

		m := NewManager()
		if err := m.LoadContractsFromDir("./testFiles"); err != nil {
			t.Error(err)
		}

		ccs := m.GetContractsByNames([]string{"delegation_controller", "validator_service"})
		for k := range ccs {
			fq.Addresses = append(fq.Addresses, k)
		}

		logs, err := client.FilterLogs(ctx, fq)

		t.Logf("Got  %d logs (%+v)", len(logs), fq)
		//	for _, l := range logs {
		/*
			c, ok := ccs[l.Address]
			if !ok {
				t.Error(err)
				return
			}
			if len(l.Topics) == 0 {
				t.Error("topics are empty")
				return
			}
			t.Logf("Got contract %s, height (%d)", c.Name, l.BlockNumber)
			event, err := c.Abi.EventByID(l.Topics[0])
			if !ok {
				t.Error(err)
				return
			}
			t.Logf("Got event %s", event.Name)
			mapped := make(map[string]interface{}, len(event.Inputs))
			event.Inputs.UnpackIntoMap(mapped, l.Data)
			t.Logf("Got values %s %+v ", event.Name, mapped)

			switch c.Name {
			case "validator_service":
				switch event.Name {
				case "ValidatorRegistered":
					ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
					defer cancel()
					bc := bind.NewBoundContract(c.Addr, c.Abi, client, nil, nil)

					t.Logf("Calling getValidator %+v", mapped["validatorId"])

					vID := mapped["validatorId"]
					results := []interface{}{}
					err = bc.Call(&bind.CallOpts{
						Pending: false,
						Context: ctxT,
					}, &results, "getValidator", vID)

					if err != nil {
						t.Error(err)
					}
					t.Logf("Got getValidator response    %T", results[0])

					vr := &ValidatorRaw{}
					v := *abi.ConvertType(results[0], vr).(*ValidatorRaw)
					t.Logf("Got getValidator response %T %+v  %T", v, v, results[0])

				}
			}
		*/
		//	}
		/*
			header, err := client.HeaderByNumber(context.Background(), nil)
			if err != nil {
				log.Fatal(err)
			}
		*/
		//fmt.Println(header.Number.String())
	})
}
