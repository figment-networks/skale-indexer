package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/figment-networks/skale-indexer/scraper/structs"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Caller) GetValidatorNodes(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (nodes []structs.Node, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		co.Pending = true
	}
	results := []interface{}{}

	err = bc.Call(co, &results, "getValidatorNodeIndexes", validatorID)

	if err != nil {
		return nil, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}
	nodes = []structs.Node{}

	r := results[0].([]*big.Int)
	for _, id := range r {
		n, err := c.GetNode(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling GetNode function %w", err)
		}
		nodes = append(nodes, n)

		nrd, err := c.GetNodeNextRewardDate(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling GetNodeNextRewardDate function %w", err)
		}
		n.NextRewardDate = nrd
	}

	return nodes, nil
}

func (c *Caller) GetNodeNextRewardDate(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (t time.Time, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		co.Pending = true
	}
	results := []interface{}{}

	err = bc.Call(co, &results, "getNodeNextRewardDate", nodeID)

	if err != nil {
		return t, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return t, errors.New("empty result")
	}

	nrDate := results[0].(*big.Int)
	return time.Unix(nrDate.Int64(), 0), nil
}

func (c *Caller) GetNode(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	results := []interface{}{}

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		co.Pending = true
	}

	err = bc.Call(co, &results, "nodes", nodeID)

	if err != nil {
		return n, fmt.Errorf("error getting nodes function %w", err)
	}

	if len(results) == 0 {
		return n, errors.New("empty result")
	}

	lrDate := results[5].(*big.Int)
	return structs.Node{
		NodeID:         nodeID,
		Name:           results[0].(string),
		IP:             results[1].([4]byte),
		PublicIP:       results[2].([4]byte),
		Port:           results[3].(uint16),
		StartBlock:     results[4].(*big.Int),
		LastRewardDate: time.Unix(lrDate.Int64(), 0),
		FinishTime:     results[6].(*big.Int),
		Status:         structs.NodeStatus(results[7].(uint8)),
		ValidatorID:    results[8].(*big.Int),
	}, nil
}
