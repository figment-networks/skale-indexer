package contractor

import (
	"../../structs"
	"../../types"
	"context"
)

type validatorContractor interface {
	SaveOrUpdateValidator(ctx context.Context, validator structs.Validator) error
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidatorById(ctx context.Context, id *types.ID) (res structs.Validator, err error)
	GetValidatorsByValidatorAddress(ctx context.Context, validatorAddress *string) (validators []structs.Validator, err error)
	GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress *string) (validators []structs.Validator, err error)
}
