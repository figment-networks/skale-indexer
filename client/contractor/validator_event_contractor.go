package contractor

import (
	"../../structs"
	"context"
)

type validatorEventStore interface {
	SaveOrUpdateValidatorEvents(ctx context.Context, validatorEvents []structs.ValidatorEvent) error
	GetValidatorEventById(ctx context.Context, id string) (res structs.ValidatorEvent, err error)
	GetValidatorEventsByValidatorId(ctx context.Context, validatorId string) (validatorEvents []structs.ValidatorEvent, err error)
	GetAllValidatorEvents(ctx context.Context) (validatorEvents []structs.ValidatorEvent, err error)
}
