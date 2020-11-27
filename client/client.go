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

func (c *ClientContractor) SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error {
	return c.storeEng.SaveOrUpdateDelegations(ctx, delegations)
}

func (c *ClientContractor) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegations(ctx, params)
}

func (c *ClientContractor) SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error {
	return c.storeEng.SaveOrUpdateEvents(ctx, events)
}

func (c *ClientContractor) GetEvents(ctx context.Context, params structs.QueryParams) (events []structs.Event, err error) {
	return c.storeEng.GetEvents(ctx, params)
}

func (c *ClientContractor) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	return c.storeEng.SaveOrUpdateValidators(ctx, validators)
}

func (c *ClientContractor) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidators(ctx, params)
}

func (c *ClientContractor) SaveOrUpdateNodes(ctx context.Context, nodes []structs.Node) error {
	return c.storeEng.SaveOrUpdateNodes(ctx, nodes)
}

func (c *ClientContractor) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	return c.storeEng.GetNodes(ctx, params)
}

func (c *ClientContractor) GetDelegationStateStatistics(ctx context.Context, params structs.QueryParams) (delegationStateStatistics []structs.DelegationStateStatistics, err error) {
	return c.storeEng.GetDelegationStateStatistics(ctx, params)
}
