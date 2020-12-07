package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type validatorStatisticsContractor interface {
	GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error)
}
