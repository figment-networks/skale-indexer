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

/*
EVENTS
Sent(operator, from, to, amount, data, operatorData)
Minted(operator, to, amount, data, operatorData)
Burned(operator, from, amount, data, operatorData)
AuthorizedOperator(operator, tokenHolder)
RevokedOperator(operator, tokenHolder)
*/

type ERC777Call interface {
	Name(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (n string, err error)
	Symbol(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (s string, err error)
	Granularity(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (g big.Int, err error)
	TotalSupply(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (ts big.Int, err error)
	BalanceOf(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, tokenHolder common.Address) (balance big.Int, err error)
	Send(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, recipient common.Address, amount *big.Int, data []byte) (err error)
	Burn(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, amount *big.Int, data []byte) (err error)
	IsOperatorFor(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator, tokenHolder common.Address) (res bool, err error)
	AuthorizeOperator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator common.Address) (err error)
	RevokeOperator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator common.Address) (err error)
	DefaultOperators(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (operators []common.Address, err error)
	OperatorSend(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, sender, recipient common.Address, amount *big.Int, data []byte, operatorData []byte) (err error)
	OperatorBurn(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, account common.Address, amount *big.Int, data []byte, operatorData []byte) (err error)
}

type ERC777Caller struct{}

func (c *ERC777Caller) Name(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (n string, err error) {
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
	err = bc.Call(co, &results, "name")

	if err != nil {
		return n, fmt.Errorf("error calling name function %w", err)
	}

	if len(results) == 0 {
		return n, errors.New("empty result")
	}

	n, ok := results[0].(string)
	if !ok {
		return n, errors.New("total supply is not string type")
	}

	return n, nil
}

func (c *ERC777Caller) Symbol(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (s string, err error) {
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
	err = bc.Call(co, &results, "symbol")

	if err != nil {
		return s, fmt.Errorf("error calling symbol function %w", err)
	}

	if len(results) == 0 {
		return s, errors.New("empty result")
	}

	s, ok := results[0].(string)
	if !ok {
		return s, errors.New("total supply is not string type")
	}

	return s, nil
}

func (c *ERC777Caller) Granularity(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (g big.Int, err error) {
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
	err = bc.Call(co, &results, "granularity")

	if err != nil {
		return g, fmt.Errorf("error calling granularity function %w", err)
	}

	if len(results) == 0 {
		return g, errors.New("empty result")
	}

	grn, ok := results[0].(*big.Int)
	if !ok {
		return g, errors.New("total supply is not *big.Int type")
	}

	return *grn, nil
}

func (c *ERC777Caller) TotalSupply(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (ts big.Int, err error) {
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

func (c *ERC777Caller) Send(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, recipient common.Address, amount *big.Int, data []byte) (err error) {
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
	err = bc.Call(co, &results, "send", recipient, amount, data)
	if err != nil {
		return fmt.Errorf("error calling send function %w", err)
	}

	return nil
}

func (c *ERC777Caller) Burn(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, amount *big.Int, data []byte) (err error) {
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
	err = bc.Call(co, &results, "burn", amount, data)
	if err != nil {
		return fmt.Errorf("error calling burn function %w", err)
	}

	return nil
}

func (c *ERC777Caller) IsOperatorFor(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator, tokenHolder common.Address) (res bool, err error) {
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
	err = bc.Call(co, &results, "isOperatorFor", operator, tokenHolder)

	if err != nil {
		return res, fmt.Errorf("error calling isOperatorFor function %w", err)
	}

	if len(results) == 0 {
		return res, errors.New("empty result")
	}

	res, ok := results[0].(bool)
	if !ok {
		return res, errors.New("balance is not bool type")
	}

	return res, nil
}

func (c *ERC777Caller) AuthorizeOperator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator common.Address) (err error) {
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
	err = bc.Call(co, &results, "authorizeOperator", operator)
	if err != nil {
		return fmt.Errorf("error calling authorizeOperator function %w", err)
	}

	return nil
}

func (c *ERC777Caller) RevokeOperator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, operator common.Address) (err error) {
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
	err = bc.Call(co, &results, "revokeOperator", operator)
	if err != nil {
		return fmt.Errorf("error calling revokeOperator function %w", err)
	}

	return nil
}

func (c *ERC777Caller) DefaultOperators(ctx context.Context, bc *bind.BoundContract, blockNumber uint64) (operators []common.Address, err error) {

	ctxTA, cancelA := context.WithTimeout(ctx, time.Second*30)
	co := &bind.CallOpts{
		Context: ctxTA,
	}

	if blockNumber > 0 { // (lukanus): 0 = latest
		co.BlockNumber = new(big.Int).SetUint64(blockNumber)
		co.Pending = true
	}
	resultsA := []interface{}{}

	err = bc.Call(co, &resultsA, "defaultOperators")

	cancelA()
	if err != nil {
		return nil, fmt.Errorf("error calling defaultOperators function %w", err)
	}

	if len(resultsA) == 0 {
		return nil, errors.New("empty result")
	}

	operators, ok := resultsA[0].([]common.Address)
	if !ok {
		return operators, errors.New("operators is not array of address type")
	}
	return operators, nil
}

func (c *ERC777Caller) OperatorSend(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, sender, recipient common.Address, amount *big.Int, data []byte, operatorData []byte) (err error) {
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
	err = bc.Call(co, &results, "operatorSend", sender, recipient, amount, data, operatorData)
	if err != nil {
		return fmt.Errorf("error calling operatorSend function %w", err)
	}

	return nil
}

func (c *ERC777Caller) OperatorBurn(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, account common.Address, amount *big.Int, data []byte, operatorData []byte) (err error) {
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
	err = bc.Call(co, &results, "operatorBurn", account, amount, data, operatorData)
	if err != nil {
		return fmt.Errorf("error calling operatorBurn function %w", err)
	}

	return nil
}
