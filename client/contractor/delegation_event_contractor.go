package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type delegationEventContractor interface {
	SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error
	GetDelegationEventById(ctx context.Context, id string) (res structs.DelegationEvent, err error)
	GetDelegationEventsByDelegationId(ctx context.Context, delegationId uint64) (delegationEvents []structs.DelegationEvent, err error)
	GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error)
}
