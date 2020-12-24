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

	FUNCTIONS
	totalSupply()
	balanceOf(account)
	transfer(recipient, amount)
	allowance(owner, spender)
	approve(spender, amount)
	transferFrom(sender, recipient, amount)
*/

/*
	EVENTS
	Transfer(from, to, value)
	Approval(owner, spender, value)
*/

type ERC20Caller struct{}

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
