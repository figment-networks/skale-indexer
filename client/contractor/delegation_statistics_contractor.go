package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type delegationStatisticsContractor interface {
	GetDelegationStatistics(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error)
	GetLatestDelegationStates(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error)
}
