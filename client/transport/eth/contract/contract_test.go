package contract

import (
	"encoding/json"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

		ccs := m.GetContractsByContractNames([]string{"delegation_controller", "validator_service"})

		require.Len(t, ccs, 2)

	})
}
