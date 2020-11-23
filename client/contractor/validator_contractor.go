package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type validatorContractor interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error)
	GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error)
	GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress string) (validators []structs.Validator, err error)
}
