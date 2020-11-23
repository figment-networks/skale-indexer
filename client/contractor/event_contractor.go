package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type eventContractor interface {
	SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error
	GetEventById(ctx context.Context, id string) (res structs.Event, err error)
	GetAllEvents(ctx context.Context) (events []structs.Event, err error)
}
