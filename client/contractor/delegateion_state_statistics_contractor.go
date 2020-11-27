package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type delegationStateStatisticsContractor interface {
	GetDelegationStateStatistics(ctx context.Context, params structs.QueryParams) (delegationStateStatistics []structs.DelegationStateStatistics, err error)
}
