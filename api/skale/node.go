package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func (c *Caller) GetValidatorNodes(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (nodes []structs.Node, err error) {
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

	contr := bc.GetContract()
	if contr == nil {
		return nil, fmt.Errorf("Contract is nil")
	}

	if err = contr.Call(co, &results, "getValidatorNodeIndexes", validatorID); err != nil {
		return nil, fmt.Errorf("error calling getValidatorNodeIndexes function %w", err)
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

func (c *Caller) GetNodeNextRewardDate(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (t time.Time, err error) {
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

	contr := bc.GetContract()
	if contr == nil {
		return t, fmt.Errorf("Contract is nil")
	}

	if err = contr.Call(co, &results, "getNodeNextRewardDate", nodeID); err != nil {
		return t, fmt.Errorf("error calling node function %w", err)
	}

	if len(results) == 0 {
		return t, errors.New("empty result")
	}

	nrDate := results[0].(*big.Int)
	return time.Unix(nrDate.Int64(), 0), nil
}

func (c *Caller) GetNodeAddress(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (adr common.Address, err error) {
	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (eesmerdag): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
	}
	results := []interface{}{}

	contr := bc.GetContract()
	if contr == nil {
		return adr, fmt.Errorf("Contract is nil")
	}

	if err = contr.Call(co, &results, "getNodeAddress", nodeID); err != nil {
		_, err2 := bc.RawCall(ctx, co, "getNodeAddress", nodeID)
		if err2 == transport.ErrEmptyResponse {
			return adr, nil
		}
		return adr, fmt.Errorf("error calling node function %w ", err)
	}

	if len(results) == 0 {
		return adr, errors.New("empty result")
	}

	adr = results[0].(common.Address)
	return adr, nil
}

func (c *Caller) GetNode(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error) {

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

	contr := bc.GetContract()
	if contr == nil {
		return n, fmt.Errorf("Contract is nil")
	}

	if err = contr.Call(co, &results, "nodes", nodeID); err != nil {
		_, err2 := bc.RawCall(ctx, co, "nodes", nodeID)
		if err2 == transport.ErrEmptyResponse {
			return n, err2
		}
		return n, fmt.Errorf("error calling nodes  %w ", err)
	}

	if len(results) == 0 {
		return n, errors.New("empty result")
	}

	lrDate := results[5].(*big.Int)
	IP := results[1].([4]byte)
	publicIP := results[2].([4]byte)
	return structs.Node{
		NodeID:         new(big.Int).Set(nodeID),
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

func (c *Caller) GetNodeWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, nodeID *big.Int) (n structs.Node, err error) {
	n, err = c.GetNode(ctx, bc, blockNumber, nodeID)
	if err != nil {
		return n, err
	}
	nrd, err := c.GetNodeNextRewardDate(ctx, bc, blockNumber, nodeID)
	if err != nil {
		return n, err
	}
	n.NextRewardDate = nrd

	adr, err := c.GetNodeAddress(ctx, bc, blockNumber, nodeID)
	if err != nil {
		return n, err
	}
	n.Address = adr

	return n, nil
}
