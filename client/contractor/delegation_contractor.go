package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type delegationContractor interface {
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error)
	GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error)
	GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error)
}
