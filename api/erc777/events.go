package erc777

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type ERC777Caller struct{}

func (c *ERC777Caller) BalanceOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, tokenHolder common.Address) (balance big.Int, err error) {
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

/*
FUNCTIONS
name()
symbol()
granularity()
totalSupply()
balanceOf(owner)
send(recipient, amount, data)
burn(amount, data)
isOperatorFor(operator, tokenHolder)
authorizeOperator(operator)
revokeOperator(operator)
defaultOperators()
operatorSend(sender, recipient, amount, data, operatorData)
operatorBurn(account, amount, data, operatorData)
*/

/*
EVENTS
Sent(operator, from, to, amount, data, operatorData)
Minted(operator, to, amount, data, operatorData)
Burned(operator, from, amount, data, operatorData)
AuthorizedOperator(operator, tokenHolder)
RevokedOperator(operator, tokenHolder)
*/
