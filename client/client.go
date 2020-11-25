package client

import (
	"context"
	"errors"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
)

var (
	InvalidId = errors.New("invalid id")
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
