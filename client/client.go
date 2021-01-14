package client

import (
	"context"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/store"
	"go.uber.org/zap"
)

type Client struct {
	storeEng store.DataStore
	log      *zap.Logger
}

func NewClient(log *zap.Logger, storeEng store.DataStore) *Client {
	return &Client{
		storeEng: storeEng,
		log:      log,
	}
}

func (c *Client) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {
	ev, err := c.storeEng.GetContractEvents(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetContractEvents", zap.Any("params", params), zap.Error(err))
	}
	return ev, err
}

func (c *Client) GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error) {
	n, err := c.storeEng.GetNodes(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetNodes", zap.Any("params", params), zap.Error(err))
	}
	return n, err
}

func (c *Client) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	v, err := c.storeEng.GetValidators(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetValidators", zap.Any("params", params), zap.Error(err))
	}
	return v, err
}

func (c *Client) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	d, err := c.storeEng.GetDelegations(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetDelegations", zap.Any("params", params), zap.Error(err))
	}
	return d, err
}

func (c *Client) GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	d, err := c.storeEng.GetDelegationTimeline(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetDelegationTimeline", zap.Any("params", params), zap.Error(err))
	}
	return d, err
}

func (c *Client) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	vs, err := c.storeEng.GetValidatorStatistics(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetValidatorStatistics", zap.Any("params", params), zap.Error(err))
	}
	return vs, err
}

func (c *Client) GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	vs, err := c.storeEng.GetValidatorStatisticsTimeline(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetValidatorStatisticsTimeline", zap.Any("params", params), zap.Error(err))
	}
	return vs, err
}

func (c *Client) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	a, err := c.storeEng.GetAccounts(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetAccounts", zap.Any("params", params), zap.Error(err))
	}
	return a, err
}

func (c *Client) GetSystemEvents(ctx context.Context, params structs.SystemEventParams) (systemEvents []structs.SystemEvent, err error) {
	systemEvents, err = c.storeEng.GetSystemEvents(ctx, params)
	if err != nil {
		c.log.Error("[CLIENT] Error in GetContractEvents:", zap.Any("params", params), zap.Error(err))
	}
	return systemEvents, err
}
