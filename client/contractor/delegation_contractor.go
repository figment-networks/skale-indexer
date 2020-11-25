package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type delegationContractor interface {
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error)
}
