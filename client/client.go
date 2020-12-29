package client

import (
	"context"
	"errors"

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

func (c *Client) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidators(ctx, params)
}

func (c *Client) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegations(ctx, params)
}

func (c *Client) GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegationTimeline(ctx, params)
}
func (c *Client) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return c.storeEng.GetValidatorStatistics(ctx, params)
}

func (c *Client) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	return c.storeEng.GetAccounts(ctx, params)
}

func (c *Client) GetAccountDetails(ctx context.Context, params structs.AccountParams) (accountDetails structs.AccountDetails, err error) {
	account, err := c.storeEng.GetAccounts(ctx, params)
	if err != nil {
		return accountDetails, errors.New("error getting account")
	}
	if len(account) == 0 {
		return accountDetails, errors.New("no account for address is found")
	}

	dlgParams := structs.DelegationParams{
		Holder: params.Address,
	}
	delegations, err := c.storeEng.GetDelegations(ctx, dlgParams)
	if err != nil {
		return accountDetails, errors.New("error getting delegations")
	}
	return structs.AccountDetails{
		Account:     account[0],
		Delegations: delegations,
	}, err
}
