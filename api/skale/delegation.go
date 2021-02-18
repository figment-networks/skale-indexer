package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/figment-networks/skale-indexer/scraper/structs"
	"github.com/figment-networks/skale-indexer/scraper/transport"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Delegation structure - to be used with abi.ConvertType method
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

func (c *Caller) GetDelegation(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error) {
	results := []interface{}{}

	co := &bind.CallOpts{
		Context: ctx,
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
		return d, fmt.Errorf("Contract is nil")
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return d, err
	}
	n := time.Now()
	if err = contr.Call(co, &results, "delegations", delegationID); err != nil {
		_, err2 := bc.RawCall(ctx, co, "delegations", delegationID)
		if err2 == transport.ErrEmptyResponse {
			return d, err2
		}

		rawRequestDuration.WithLabels("delegations", "err").Observe(time.Since(n).Seconds())
		return d, fmt.Errorf("error calling delegations  %w ", err)
	}
	rawRequestDuration.WithLabels("delegations", "ok").Observe(time.Since(n).Seconds())

	if len(results) == 0 {
		return d, errors.New("empty result")
	}

	if len(results) < 8 {
		return d, errors.New("wrong type of result")
	}

	createT := results[4].(*big.Int)
	dg := structs.Delegation{
		DelegationID:     new(big.Int).Set(delegationID),
		Holder:           results[0].(common.Address),
		ValidatorID:      results[1].(*big.Int),
		Amount:           results[2].(*big.Int),
		DelegationPeriod: results[3].(*big.Int),
		Created:          time.Unix(createT.Int64(), 0),
		Started:          results[5].(*big.Int),
		Finished:         results[6].(*big.Int),
		Info:             results[7].(string),
		State:            structs.DelegationStateUNKNOWN,
	}

	// BUG(lukanus): recover from panic after format update!
	return dg, nil
}

func (c *Caller) GetDelegationState(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (ds structs.DelegationState, err error) {
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

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return ds, err
	}

	n := time.Now()
	err = bc.Call(co, &results, "getState", delegationID)
	if err != nil {
		rawRequestDuration.WithLabels("getStateDelegation", "err").Observe(time.Since(n).Seconds())
		return ds, fmt.Errorf("error calling getStateDelegation function %w", err)
	}
	rawRequestDuration.WithLabels("getStateDelegation", "ok").Observe(time.Since(n).Seconds())

	if len(results) == 0 {
		return ds, errors.New("empty result")
	}

	state := *abi.ConvertType(results[0], new(uint8)).(*uint8)
	return structs.DelegationState(state), nil
}

func (c *Caller) GetPendingDelegationsTokens(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holderAddress common.Address) (amount *big.Int, err error) {
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

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	n := time.Now()
	if err = bc.Call(co, &results, "getLockedInPendingDelegations", holderAddress); err != nil {
		rawRequestDuration.WithLabels("getLockedInPendingDelegations", "err").Observe(time.Since(n).Seconds())
		return nil, fmt.Errorf("error calling getLockedInPendingDelegations function %w", err)
	}
	rawRequestDuration.WithLabels("getLockedInPendingDelegations", "ok").Observe(time.Since(n).Seconds())

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

func (c *Caller) GetValidatorDelegations(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (delegations []structs.Delegation, err error) {

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	caller := bc.GetContract()

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
	n := time.Now()
	if err = caller.Call(co, &results, "getDelegationsByValidatorLength", validatorID); err != nil {
		rawRequestDuration.WithLabels("getDelegationsByValidatorLength", "err").Observe(time.Since(n).Seconds())
		return nil, fmt.Errorf("error calling getDelegationsByValidatorLength function %w", err)
	}
	rawRequestDuration.WithLabels("getDelegationsByValidatorLength", "ok").Observe(time.Since(n).Seconds())

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	count, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("count is not *big.Int type ")
	}

	delegations = []structs.Delegation{}

	for i := uint64(0); i < count.Uint64(); i++ {

		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}

		ctxTA, cancelA := context.WithTimeout(ctx, time.Second*30)
		co := &bind.CallOpts{
			Context: ctxTA,
		}

		if c.NodeType == ENTArchive {
			if blockNumber > 0 { // (lukanus): 0 = latest
				co.BlockNumber = new(big.Int).SetUint64(blockNumber)
			} else {
				co.Pending = true
			}
		}

		resultsA := []interface{}{}

		n := time.Now()
		err = caller.Call(co, &resultsA, "delegationsByValidator", validatorID, new(big.Int).SetUint64(i))
		cancelA()
		if err != nil {
			rawRequestDuration.WithLabels("delegationsByValidator", "err").Observe(time.Since(n).Seconds())
			return nil, fmt.Errorf("error calling delegationsByValidator function %w", err)
		}
		rawRequestDuration.WithLabels("delegationsByValidator", "ok").Observe(time.Since(n).Seconds())

		if len(resultsA) == 0 {
			return nil, errors.New("empty result")
		}

		id, ok := resultsA[0].(*big.Int)
		if !ok {
			return nil, errors.New("delegation id is not a bigint")
		}

		d, err := c.GetDelegation(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling delegations function: %w", err)
		}

		d.State, err = c.GetDelegationState(ctx, caller, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error getting delegation state %w", err)
		}

		delegations = append(delegations, d)
	}

	return delegations, nil
}

type IdError struct {
	ID  uint64
	Err error
}

func (c *Caller) GetValidatorDelegationsIDs(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (delegationsIDs []uint64, err error) {

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	caller := bc.GetContract()

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
	n := time.Now()
	if err = caller.Call(co, &results, "getDelegationsByValidatorLength", validatorID); err != nil {
		rawRequestDuration.WithLabels("getDelegationsByValidatorLength", "err").Observe(time.Since(n).Seconds())
		return nil, fmt.Errorf("error calling getDelegationsByValidatorLength function %w", err)
	}
	rawRequestDuration.WithLabels("getDelegationsByValidatorLength", "ok").Observe(time.Since(n).Seconds())

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	count, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("count is not *big.Int type ")
	}

	vdC, ok := c.validatorDelegationsCache[validatorID.Uint64()]
	if ok {
		if count.Uint64() == vdC.Length {
			return vdC.Delegations, nil
		}
	}
	lID := vdC.LastID

	in := make(chan uint64)
	out := make(chan IdError, count.Uint64()-lID)
	defer close(in)
	defer close(out)
	for i := 0; i < 5; i++ {
		go c.delegationsByValidatorAsync(ctx, bc, validatorID, blockNumber, in, out)
	}

	for i := uint64(lID); i < count.Uint64(); i++ {
		in <- i
		vdC.LastID = i
	}

	for j := uint64(lID); j < count.Uint64(); j++ {
		d := <-out
		vdC.Delegations = append(vdC.Delegations, d.ID)
		if d.Err != nil {
			err = d.Err
		}
	}

	vdC.LastID = count.Uint64()
	c.validatorDelegationsCache[validatorID.Uint64()] = vdC
	return vdC.Delegations, err
}

func (c *Caller) delegationsByValidatorAsync(ctx context.Context, bc transport.BoundContractCaller, validatorID *big.Int, blockNumber uint64, in <-chan uint64, out chan<- IdError) {
	for nmbr := range in {
		delegationID, err := c.DelegationsByValidator(ctx, bc, blockNumber, validatorID, nmbr)
		out <- IdError{delegationID.Uint64(), err}
	}
}

func (c *Caller) DelegationsByValidator(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int, delegationsNumber uint64) (delegationID *big.Int, err error) {

	results := []interface{}{}

	contr := bc.GetContract()
	if contr == nil {
		return delegationID, fmt.Errorf("Contract is nil")
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return delegationID, err
	}

	ctxT, cancel := context.WithTimeout(ctx, time.Second*10)
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
	n := time.Now()
	if err = contr.Call(co, &results, "delegationsByValidator", validatorID, new(big.Int).SetUint64(delegationsNumber)); err != nil {
		rawRequestDuration.WithLabels("delegationsByValidator", "err").Observe(time.Since(n).Seconds())
		cancel()
		return delegationID, fmt.Errorf("error calling delegationsByValidator  %w ", err)
	}
	rawRequestDuration.WithLabels("delegationsByValidator", "ok").Observe(time.Since(n).Seconds())
	cancel()

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	id, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("delegation id is not a bigint")
	}

	return id, nil
}

func (c *Caller) GetHolderDelegations(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, holder common.Address) (delegations []structs.Delegation, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	caller := bc.GetContract()
	co := &bind.CallOpts{
		Context: ctxT,
		Pending: true,
	}

	if c.NodeType == ENTArchive {
		if blockNumber > 0 { // (lukanus): 0 = latest
			co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		} else {
			co.Pending = true
		}
	}

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, err
	}
	results := []interface{}{}
	n := time.Now()
	if err = caller.Call(co, &results, "getDelegationsByHolderLength", holder); err != nil {
		rawRequestDuration.WithLabels("getDelegationsByHolderLength", "err").Observe(time.Since(n).Seconds())
		return nil, fmt.Errorf("error calling getDelegationsByHolderLength function %w", err)
	}
	rawRequestDuration.WithLabels("getDelegationsByHolderLength", "ok").Observe(time.Since(n).Seconds())

	if len(results) == 0 {
		return nil, errors.New("empty result")
	}

	count, ok := results[0].(*big.Int)
	if !ok {
		return nil, errors.New("count is not *big.Int type ")
	}

	delegations = []structs.Delegation{}

	for i := uint64(0); i < count.Uint64(); i++ {

		ctxTA, cancelA := context.WithTimeout(ctx, time.Second*30)
		co := &bind.CallOpts{
			Context: ctxTA,
		}

		if c.NodeType == ENTArchive {
			if blockNumber > 0 { // (lukanus): 0 = latest
				co.BlockNumber = new(big.Int).SetUint64(blockNumber)
			} else {
				co.Pending = true
			}
		}

		resultsA := []interface{}{}

		if err := c.rateLimiter.Wait(ctx); err != nil {
			return nil, err
		}
		n := time.Now()
		err = caller.Call(co, &resultsA, "delegationsByHolder", holder, new(big.Int).SetUint64(i))
		cancelA()
		if err != nil {
			rawRequestDuration.WithLabels("delegationsByHolder", "err").Observe(time.Since(n).Seconds())
			return nil, fmt.Errorf("error calling delegationsByHolder function %w", err)
		}

		rawRequestDuration.WithLabels("delegationsByHolder", "ok").Observe(time.Since(n).Seconds())

		if len(resultsA) == 0 {
			return nil, errors.New("empty result")
		}

		id, ok := resultsA[0].(*big.Int)
		if !ok {
			return nil, errors.New("delegation id is not a bigint")
		}

		d, err := c.GetDelegation(ctx, bc, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error calling delegations function: %w", err)
		}

		d.State, err = c.GetDelegationState(ctx, caller, blockNumber, id)
		if err != nil {
			return nil, fmt.Errorf("error getting delegation state %w", err)
		}

		delegations = append(delegations, d)
	}

	return delegations, nil
}

// GetDelegationInfo delegation info with all parameters
func (c *Caller) GetDelegationWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error) {
	delegation, err := c.GetDelegation(ctx, bc, blockNumber, delegationID)
	if err != nil {
		return delegation, err
	}
	delegation.State, err = c.GetDelegationState(ctx, bc.GetContract(), blockNumber, delegationID)
	if err != nil {
		return delegation, err
	}

	return delegation, nil
}
