package contractor

import (
	"../../structs"
	"../../types"
	"context"
)

type delegationContractor interface {
	SaveOrUpdateDelegation(ctx context.Context, delegation structs.Delegation) error
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegationById(ctx context.Context, id *types.ID) (res structs.Delegation, err error)
	GetDelegationsByHolder(ctx context.Context, holder *string) (delegations []structs.Delegation, err error)
	GetDelegationsByValidatorId(ctx context.Context, validatorId *uint64) (delegations []structs.Delegation, err error)
}
