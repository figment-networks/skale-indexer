package client

import (
	"context"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/store"
)

type Client struct {
	storeEng store.DataStore
}

func NewClient(storeEng store.DataStore) *Client {
	return &Client{
		storeEng: storeEng,
	}
}

func (c *Client) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {
	return c.storeEng.GetContractEvents(ctx, params)
}

func (c *Client) GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error) {
	return c.storeEng.GetNodes(ctx, params)
}

func (c *Client) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidators(ctx, params)
}

func (c *Client) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegations(ctx, params)
}

func (c *Client) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return c.storeEng.GetValidatorStatistics(ctx, params)
}
