package skale

import (
	"context"
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"net"
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

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
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

		adr, err := c.GetNodeAddress(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling GetNodeAddress function %w", err)
		}
		n.Address = adr
	}

	return nodes, nil
}

func (c *Caller) GetNodeNextRewardDate(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (t time.Time, err error) {
	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
	}
	results := []interface{}{}

	err = bc.Call(co, &results, "getNodeNextRewardDate", nodeID)

	if err != nil {
		return t, fmt.Errorf("error calling node function %w", err)
	}

	if len(results) == 0 {
		return t, errors.New("empty result")
	}

	nrDate := results[0].(*big.Int)
	return time.Unix(nrDate.Int64(), 0), nil
}

func (c *Caller) GetNodeAddress(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (adr common.Address, err error) {
	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
	}
	results := []interface{}{}

	err = bc.Call(co, &results, "getNodeAddress", nodeID)

	if err != nil {
		return adr, fmt.Errorf("error calling node function %w", err)
	}

	if len(results) == 0 {
		return adr, errors.New("empty result")
	}

	adr = results[0].(common.Address)
	return adr, nil
}

func (c *Caller) GetNode(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	results := []interface{}{}

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
	}
	err = bc.Call(co, &results, "nodes", nodeID)

	if err != nil {
		return n, fmt.Errorf("error getting nodes function %w", err)
	}

	if len(results) == 0 {
		return n, errors.New("empty result")
	}

	lrDate := results[5].(*big.Int)
	IP := results[1].([4]byte)
	publicIP := results[1].([4]byte)
	return structs.Node{
		NodeID:         nodeID,
		Name:           results[0].(string),
		IP:             net.IPv4(IP[0], IP[1], IP[2], IP[3]),
		PublicIP:       net.IPv4(publicIP[0], publicIP[1], publicIP[2], publicIP[3]),
		Port:           results[3].(uint16),
		StartBlock:     results[4].(*big.Int),
		LastRewardDate: time.Unix(lrDate.Int64(), 0),
		FinishTime:     results[6].(*big.Int),
		Status:         structs.NodeStatus(results[7].(uint8)),
		ValidatorID:    results[8].(*big.Int),

		BlockHeight: blockNumber,
	}, nil
}

func (c *Caller) GetAllCurrentNodes(ctx context.Context, bc *bind.BoundContract) (nodes []structs.Node, err error) {
	nodes = []structs.Node{}
	zeroBlockNumber := uint64(0)
	nodeID := int64(1)

	for {
		nIDBig := big.NewInt(nodeID)

		n, err := c.GetNode(ctx, bc, zeroBlockNumber, nIDBig)
		if err != nil {
			break
		}
		nrd, err := c.GetNodeNextRewardDate(ctx, bc, zeroBlockNumber, nIDBig)
		if err != nil {
			break
		}
		n.NextRewardDate = nrd
		adr, err := c.GetNodeAddress(ctx, bc, zeroBlockNumber, nIDBig)
		if err != nil {
			break
		}
		n.Address = adr

		nodeID++
		nodes = append(nodes, n)
	}

	return nodes, nil
}
