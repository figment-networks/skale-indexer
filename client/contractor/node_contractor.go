package contractor

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type nodeContractor interface {
	SaveOrUpdateNodes(ctx context.Context, nodes []structs.Node) error
	GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error)
}
