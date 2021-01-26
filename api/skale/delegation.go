package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/figment-networks/skale-indexer/scraper/structs"

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

func (c *Caller) GetDelegation(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error) {
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

	err = bc.Call(co, &results, "delegations", delegationID)

	if err != nil {
		return d, fmt.Errorf("error calling delegations function %w", err)
	}

	if len(results) == 0 {
		return d, errors.New("empty result")
	}

	if len(results) < 8 {
		return d, errors.New("wrong type of result")
	}

	createT := results[4].(*big.Int)
	dID := big.NewInt(delegationID.Int64())
	dg := structs.Delegation{
		DelegationID:     dID,
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
	err = bc.Call(co, &results, "getState", delegationID)

	if err != nil {
		return ds, fmt.Errorf("error calling delegations function %w", err)
	}

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

func (c *Caller) GetValidatorDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (delegations []structs.Delegation, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

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

func (c *Caller) GetHolderDelegations(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, holder common.Address) (delegations []structs.Delegation, err error) {

	ctxT, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

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

	results := []interface{}{}

	err = bc.Call(co, &results, "getDelegationsByHolderLength", holder)

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

		err = bc.Call(co, &resultsA, "delegationsByHolder", holder, new(big.Int).SetUint64(i))

		cancelA()
		if err != nil {
			return nil, fmt.Errorf("error calling delegationsByHolder function %w", err)
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

// GetDelegationInfo delegation info with all parameters
func (c *Caller) GetDelegationWithInfo(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, delegationID *big.Int) (d structs.Delegation, err error) {
	delegation, err := c.GetDelegation(ctx, bc, blockNumber, delegationID)
	if err != nil {
		return delegation, fmt.Errorf("error calling GetDelegation %w", err)
	}
	delegation.State, err = c.GetDelegationState(ctx, bc, blockNumber, delegationID)
	if err != nil {
		return delegation, fmt.Errorf("error calling GetDelegationState %w", err)
	}

	return delegation, nil
}

/* gets 10 delegations based on ind parameter
 * to be used for synchronization
 *
 * example: if ind is 5, then it will fetch delegations for delegation_id between 41 and 50
 */
func (c *Caller) FetchNextRoundDelegations(ctx context.Context, bc *bind.BoundContract, ind int64, currentBlock uint64, cc chan<- []structs.Delegation) {
	delegations := []structs.Delegation{}
	length := int64(10)
	dlgID := (ind-1)*length + 1
	for i := 0; i < int(length); i++ {
		dlgIDBig := big.NewInt(dlgID)
		d, err := c.GetDelegation(ctx, bc, currentBlock, dlgIDBig)
		if err != nil {
			break
		}

		d.State, err = c.GetDelegationState(ctx, bc, currentBlock, dlgIDBig)
		if err != nil {
			break
		}
		d.BlockHeight = currentBlock
		dlgID++
		delegations = append(delegations, d)
	}
	cc <- delegations
}
