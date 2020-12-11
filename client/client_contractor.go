package client

import (
	"context"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
)

type ClientContractor struct {
	storeEng store.DataStore
}

func NewClientContractor(storeEng store.DataStore) *ClientContractor {
	return &ClientContractor{
		storeEng: storeEng,
	}
}

func (c *ClientContractor) GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error) {
	return c.storeEng.GetContractEvents(ctx, params)
}

func (c *ClientContractor) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	return c.storeEng.GetNodes(ctx, params)
}

func (c *ClientContractor) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidators(ctx, params)
}
