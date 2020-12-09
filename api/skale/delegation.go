package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/figment-networks/skale-indexer/api/structures"
)

// Delegation structure - abi.Convert Types is dumb as f****
// and it is decoding data using... field order. this is why we cannot change field order
type DelegationRaw struct {
	Holder           common.Address `json:"holder"`
	ValidatorID      *big.Int       `json:"validatorId"`
	Amount           *big.Int       `json:"amount"`
	DelegationPeriod *big.Int       `json:"delegationPeriod"`
	Created          *big.Int       `json:"created"`
	Started          *big.Int       `json:"started"`
	Finished         *big.Int       `json:"finished"`
	Info             string         `json:"info"`
}

func (c *Caller) GetDelegation(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structures.Delegation, err error) {
	results := []interface{}{}

	co := &bind.CallOpts{
		Context: ctx,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		co.Pending = true
	}

	err = bc.Call(co, &results, "delegations", delegationID)

	if err != nil {
		return d, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return d, errors.New("empty result")
	}

	createT := results[4].(*big.Int)
	dg := structures.Delegation{
		ID:               delegationID,
		Holder:           results[0].(common.Address),
		ValidatorID:      results[1].(*big.Int),
		Amount:           results[2].(*big.Int),
		DelegationPeriod: results[3].(*big.Int),
		Created:          time.Unix(createT.Int64(), 0),
		Started:          results[5].(*big.Int),
		Finished:         results[6].(*big.Int),
		Info:             results[7].(string),
		State:            structures.DelegationStateUNKNOWN,
	}

	//log.Printf("gotDelegations %+v", dg)

	// BUG(lukanus): recover from panic after format update!
	return dg, nil
}

func (c *Caller) GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structures.DelegationState, err error) {
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
	err = bc.Call(co, &results, "getState", delegationID)

	if err != nil {
		return ds, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return ds, errors.New("empty result")
	}

	state := *abi.ConvertType(results[0], new(uint8)).(*uint8)
	return structures.DelegationState(state), nil
}

func (c *Caller) GetPendingDelegationsTokens(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holderAddress common.Address) (amount *big.Int, err error) {
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
	err = bc.Call(co, &results, "getLockedInPendingDelegations", holderAddress)

	if err != nil {
		return nil, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	var ok bool
	amount, ok = results[0].(*big.Int)
	if !ok {
		return nil, errors.New("amount is not *big.Int type ")
	}

	return amount, nil
}

func (c *Caller) GetValidatorDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (delegations []structures.Delegation, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	co := &bind.CallOpts{
		Context: ctxT,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
	}
	results := []interface{}{}

	err = bc.Call(co, &results, "getDelegationsByValidatorLength", validatorID)

	if err != nil {
		return nil, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	count, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("count is not *big.Int type ")
	}

	delegations = []structures.Delegation{}

	for i := uint64(0); i < count.Uint64(); i++ {

		ctxTA, cancelA := context.WithTimeout(ctx, time.Second*30)
		co := &bind.CallOpts{
			Context: ctxTA,
		}

		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
			co.Pending = true
		}
		resultsA := []interface{}{}

		err = bc.Call(co, &resultsA, "delegationsByValidator", validatorID, new(big.Int).SetUint64(i))

		cancelA()
		if err != nil {
			return nil, fmt.Errorf("error calling delegationsByValidator function %w", err)
		}

		if len(resultsA) == 0 {
			return nil, errors.New("empty result")
		}

		id, ok := resultsA[0].(*big.Int)
		if !ok {
			return nil, errors.New("delegation id is not a bigint")
		}

		d, err := c.GetDelegation(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling delegations function %w", err)
		}

		d.State, err = c.GetDelegationState(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error getting delegation state %w", err)
		}

		delegations = append(delegations, d)
	}

	return delegations, nil
}
