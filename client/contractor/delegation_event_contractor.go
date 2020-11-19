package contractor

import (
	"../../structs"
	"context"
)

type delegationEventContractor interface {
	SaveOrUpdateDelegationEvent(ctx context.Context, dl structs.DelegationEvent) error
	SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error
	GetDelegationEventById(ctx context.Context, id *string) (res structs.DelegationEvent, err error)
	GetDelegationEventsByDelegationId(ctx context.Context, delegationId *uint64) (delegationEvents []structs.DelegationEvent, err error)
	GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error)
}
