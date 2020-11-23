package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type validatorEventStore interface {
	SaveOrUpdateValidatorEvents(ctx context.Context, validatorEvents []structs.ValidatorEvent) error
	GetValidatorEventById(ctx context.Context, id string) (res structs.ValidatorEvent, err error)
	GetValidatorEventsByValidatorId(ctx context.Context, validatorId string) (validatorEvents []structs.ValidatorEvent, err error)
	GetAllValidatorEvents(ctx context.Context) (validatorEvents []structs.ValidatorEvent, err error)
}
