package skale

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (c *Caller) GetEarnedFeeAmountOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (earned, endMonth *big.Int, err error) {

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, nil, err
	}

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

	n := time.Now()
	err = bc.Call(co, &results, "getEarnedFeeAmountOf", validatorID)
	if err != nil {
		rawRequestDuration.WithLabels("getEarnedFeeAmountOf", "err").Observe(time.Since(n).Seconds())
		return earned, endMonth, fmt.Errorf("error calling getValidator function %w", err)
	}
	rawRequestDuration.WithLabels("getEarnedFeeAmountOf", "ok").Observe(time.Since(n).Seconds())

	if len(results) < 2 {
		return earned, endMonth, errors.New("empty result")
	}

	var ok bool
	earned, ok = results[0].(*big.Int)
	if !ok {
		return earned, endMonth, errors.New("earned is not *big.Int type ")
	}
	endMonth, ok = results[1].(*big.Int)
	if !ok {
		return earned, endMonth, errors.New("endMonth is not *big.Int type ")
	}

	return earned, endMonth, nil
}

func (c *Caller) GetAndUpdateEarnedBountyAmountOf(ctx context.Context, bc *bind.BoundContract, validatorID *big.Int, wallet common.Address, blockNumber uint64) (earned, endMonth *big.Int, err error) {

	if err := c.rateLimiter.Wait(ctx); err != nil {
		return nil, nil, err
	}

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

	n := time.Now()
	err = bc.Call(co, &results, "getAndUpdateEarnedBountyAmountOf", wallet, validatorID)
	if err != nil {
		rawRequestDuration.WithLabels("getAndUpdateEarnedBountyAmountOf", "err").Observe(time.Since(n).Seconds())
		return earned, endMonth, fmt.Errorf("error calling getValidator function %w", err)
	}
	rawRequestDuration.WithLabels("getAndUpdateEarnedBountyAmountOf", "ok").Observe(time.Since(n).Seconds())

	if len(results) < 2 {
		return earned, endMonth, errors.New("empty result")
	}

	var ok bool
	earned, ok = results[0].(*big.Int)
	if !ok {
		return earned, endMonth, errors.New("earned is not *big.Int type ")
	}
	endMonth, ok = results[1].(*big.Int)
	if !ok {
		return earned, endMonth, errors.New("endMonth is not *big.Int type ")
	}

	return earned, endMonth, nil
}
