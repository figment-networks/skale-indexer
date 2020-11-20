package contractor

import (
	"../../structs"
	"context"
)

type validatorContractor interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error)
	GetValidatorsByValidatorAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error)
	GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress string) (validators []structs.Validator, err error)
}
