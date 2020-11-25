package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type eventContractor interface {
	SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error
	GetEvents(ctx context.Context, params structs.QueryParams) (events []structs.Event, err error)
}
