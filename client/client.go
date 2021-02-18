package client

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport/eth/contract"
	"github.com/figment-networks/skale-indexer/store"
	"go.uber.org/zap"
)

type EthereumConnector interface {
	ParseLogs(ctx context.Context, ccs *contract.Contracts, taskID string, from, to big.Int) error
	GetLatestBlockHeight(ctx context.Context) (uint64, error)
}

type Client struct {
	storeEng store.DataStore
	log      *zap.Logger

	ethConn EthereumConnector

	ccs *contract.Contracts

	smallestPossibleHeight uint64
	maxHeightsPerRequest   uint64

	r *Running
}

func NewClient(log *zap.Logger, storeEng store.DataStore, ethConn EthereumConnector, ccs *contract.Contracts, smallestPossibleHeight, maxHeightsPerRequest uint64) *Client {
	return &Client{
		storeEng:               storeEng,
		ethConn:                ethConn,
		log:                    log,
		ccs:                    ccs,
		smallestPossibleHeight: smallestPossibleHeight,
		maxHeightsPerRequest:   maxHeightsPerRequest,
		r:                      NewRunning(),
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

func (c *Client) ParseLogs(ctx context.Context, taskID string, from, to big.Int) error {
	err := c.ethConn.ParseLogs(ctx, c.ccs, taskID, from, to)
	if err != nil {
		c.log.Error("[CLIENT] Error in ParseLogs:", zap.Uint64("from", from.Uint64()), zap.String("task_id", taskID), zap.Uint64("to", to.Uint64()), zap.Error(err))
	}
	return err
}

func (c *Client) GetLatestData(ctx context.Context, taskID string, latest uint64) (latestBlock uint64, isRunning bool, err error) {

	height, err := c.ethConn.GetLatestBlockHeight(ctx)
	if err != nil {
		return latest, false, err
	}

	c.r.lock.RLock()
	p, ok := c.r.Processes[PSig{TaskID: taskID}]
	c.r.lock.RUnlock()
	if ok {
		if !p.Finished {
			if p.Error != nil {
				c.log.Warn("[CLIENT] Last request errored", zap.Error(p.Error), zap.Uint64("height", latest), zap.Duration("since", time.Since(p.Started)))
				delete(c.r.Processes, PSig{TaskID: taskID})
				return latest, false, p.Error
			}
			c.log.Warn("[CLIENT] Last request is still processing", zap.Uint64("height", latest), zap.Duration("since", time.Since(p.Started)))
			return latest, true, nil
		}
		delete(c.r.Processes, PSig{TaskID: taskID})
		return p.EndHeight, false, p.Error
	}

	from, to := &big.Int{}, &big.Int{}
	if latest == 0 || latest <= c.smallestPossibleHeight {
		from = from.SetUint64(c.smallestPossibleHeight)
		to = to.SetUint64(c.smallestPossibleHeight + c.maxHeightsPerRequest)
	} else {
		from = from.SetUint64(latest)
		if height-latest > c.maxHeightsPerRequest {
			to = to.Add(from, new(big.Int).SetUint64(c.maxHeightsPerRequest))
		} else {
			to = to.Add(from, new(big.Int).SetUint64(height))
		}
	}

	psig := PSig{TaskID: taskID}
	c.r.lock.Lock()
	c.r.Processes[psig] = Process{
		Started:   time.Now(),
		EndHeight: to.Uint64(),
	}
	out := make(chan struct{})
	go c.getRange(ctx, taskID, *from, *to, psig, out)
	c.r.lock.Unlock()

	select {
	case <-ctx.Done():
		p, ok = c.r.Processes[psig]
		if ok {
			latest = p.EndHeight
		}
		return latest, true, p.Error
	case <-out:
	}

	c.r.lock.Lock()
	p, ok = c.r.Processes[psig]
	if ok {
		err = p.Error
		latestBlock = p.EndHeight
		delete(c.r.Processes, psig)
	}
	c.r.lock.Unlock()

	return latestBlock, false, err

}

func (c *Client) getRange(ctx context.Context, taskID string, from, to big.Int, sig PSig, out chan struct{}) {

	err := c.ethConn.ParseLogs(context.Background(), c.ccs, taskID, from, to)

	c.r.lock.Lock()
	p, ok := c.r.Processes[sig]
	if ok {
		p.Error = err
		p.Finished = true
	}
	c.r.Processes[sig] = p
	c.r.lock.Unlock()

	select {
	case <-ctx.Done():
	case out <- struct{}{}:
	}
	close(out)

}

type PSig struct {
	TaskID string
}

func NewRunning() *Running {
	return &Running{Processes: make(map[PSig]Process)}
}

type Running struct {
	lock sync.RWMutex

	Processes map[PSig]Process
}

type Process struct {
	Started   time.Time
	Finished  bool
	EndHeight uint64
	Error     error
}
