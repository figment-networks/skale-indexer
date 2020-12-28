package erc20

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

/*
	EVENTS
	Transfer(from, to, value)
	Approval(owner, spender, value)
*/

type ERC20Call interface {
	TotalSupply(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (ts big.Int, err error)
	BalanceOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, tokenHolder common.Address) (balance big.Int, err error)
	Transfer(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, recipient common.Address, amount *big.Int) (successful bool, err error)
	Allowance(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, owner, spender common.Address) (res big.Int, err error)
	Approve(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, spender common.Address, amount *big.Int) (successful bool, err error)
	TransferFrom(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, sender, recipient common.Address, amount *big.Int) (successful bool, err error)
}

type ERC20Caller struct{}

func (c *ERC20Caller) TotalSupply(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (ts big.Int, err error) {
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
	err = bc.Call(co, &results, "totalSupply")

	if err != nil {
		return ts, fmt.Errorf("error calling totalSupply function %w", err)
	}

	if len(results) == 0 {
		return ts, errors.New("empty result")
	}

	b, ok := results[0].(*big.Int)
	if !ok {
		return ts, errors.New("total supply is not *big.Int type")
	}

	return *b, nil
}

func (c *ERC20Caller) BalanceOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, tokenHolder common.Address) (balance big.Int, err error) {
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
	err = bc.Call(co, &results, "balanceOf", tokenHolder)

	if err != nil {
		return balance, fmt.Errorf("error calling balanceOf function %w", err)
	}

	if len(results) == 0 {
		return balance, errors.New("empty result")
	}

	b, ok := results[0].(*big.Int)
	if !ok {
		return balance, errors.New("balance is not *big.Int type")
	}

	return *b, nil
}

func (c *ERC20Caller) Transfer(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, recipient common.Address, amount *big.Int) (successful bool, err error) {
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
	err = bc.Call(co, &results, "transfer", recipient, amount)

	if err != nil {
		return successful, fmt.Errorf("error calling transfer function %w", err)
	}

	if len(results) == 0 {
		return successful, errors.New("empty result")
	}

	successful, ok := results[0].(bool)
	if !ok {
		return successful, errors.New("result is not boolean type")
	}

	return successful, nil
}

func (c *ERC20Caller) Allowance(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, owner, spender common.Address) (res big.Int, err error) {
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
	err = bc.Call(co, &results, "allowance", owner, spender)

	if err != nil {
		return res, fmt.Errorf("error calling allowance function %w", err)
	}

	if len(results) == 0 {
		return res, errors.New("empty result")
	}

	a, ok := results[0].(*big.Int)
	if !ok {
		return res, errors.New("balance is not *big.Int type")
	}

	return *a, nil
}

func (c *ERC20Caller) Approve(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, spender common.Address, amount *big.Int) (successful bool, err error) {
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
	err = bc.Call(co, &results, "approve", spender, amount)

	if err != nil {
		return successful, fmt.Errorf("error calling approve function %w", err)
	}

	if len(results) == 0 {
		return successful, errors.New("empty result")
	}

	successful, ok := results[0].(bool)
	if !ok {
		return successful, errors.New("result is not boolean type")
	}

	return successful, nil
}

func (c *ERC20Caller) TransferFrom(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, sender, recipient common.Address, amount *big.Int) (successful bool, err error) {
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
	err = bc.Call(co, &results, "transferFrom", sender, recipient, amount)

	if err != nil {
		return successful, fmt.Errorf("error calling transferFrom function %w", err)
	}

	if len(results) == 0 {
		return successful, errors.New("empty result")
	}

	successful, ok := results[0].(bool)
	if !ok {
		return successful, errors.New("result is not boolean type")
	}

	return successful, nil
}
