package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type validatorContractor interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error)
}
