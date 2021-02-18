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

// Validator structure - to be used with abi.ConvertType method
// It is decoding data using... field order. this is why we cannot change field order
type ValidatorRaw struct {
	Name                    string         `json:"name"`
	ValidatorAddress        common.Address `json:"validatorAddress"`
	RequestedAddress        common.Address `json:"requestedAddress"`
	Description             string         `json:"description"`
	FeeRate                 *big.Int       `json:"feeRate"`
	RegistrationTime        *big.Int       `json:"registrationTime"`
	MinimumDelegationAmount *big.Int       `json:"minimumDelegationAmount"`
	AcceptNewRequests       bool           `json:"acceptNewRequests"`
}

func (c *Caller) GetValidator(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error) {

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

	now := time.Now()
	if err = contr.Call(co, &results, "getValidator", validatorID); err != nil {

		_, err2 := bc.RawCall(ctx, co, "getValidator", validatorID)
		if err2 == transport.ErrEmptyResponse {
			rawRequestDuration.WithLabels("getValidator", "empty").Observe(time.Since(now).Seconds())
			return v, err2
		}
		rawRequestDuration.WithLabels("getValidator", "err").Observe(time.Since(now).Seconds())
		return v, fmt.Errorf("error calling getValidator  %w ", err)
	}
	rawRequestDuration.WithLabels("getValidator", "ok").Observe(time.Since(now).Seconds())

	if len(results) == 0 {
		return v, errors.New("empty result")
	}

	vr := &ValidatorRaw{}
	vraw := *abi.ConvertType(results[0], vr).(*ValidatorRaw)
	return structs.Validator{
		ValidatorID:             new(big.Int).Set(validatorID),
		Name:                    vraw.Name,
		ValidatorAddress:        vraw.ValidatorAddress,
		RequestedAddress:        vraw.RequestedAddress,
		Description:             vraw.Description,
		FeeRate:                 vraw.FeeRate,
		RegistrationTime:        time.Unix(vraw.RegistrationTime.Int64(), 0),
		MinimumDelegationAmount: vraw.MinimumDelegationAmount,
		AcceptNewRequests:       vraw.AcceptNewRequests,
	}, nil
}

func (c *Caller) IsAuthorizedValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (isAuthorized bool, err error) {

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
	now := time.Now()
	if err = bc.Call(co, &results, "isAuthorizedValidator", validatorID); err != nil {
		rawRequestDuration.WithLabels("isAuthorizedValidator", "error").Observe(time.Since(now).Seconds())
		return false, fmt.Errorf("error calling getValidator function %w", err)
	}

	rawRequestDuration.WithLabels("getValidator", "ok").Observe(time.Since(now).Seconds())

	var ok bool
	isAuthorized, ok = results[0].(bool)
	if !ok {
		return false, errors.New("earned is not bool type ")
	}

	return isAuthorized, nil
}

func (c *Caller) GetValidatorWithInfo(ctx context.Context, bc transport.BoundContractCaller, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error) {
	validator, err := c.GetValidator(ctx, bc, blockNumber, validatorID)
	if err != nil {
		return validator, err
	}
	validator.Authorized, err = c.IsAuthorizedValidator(ctx, bc.GetContract(), blockNumber, validatorID)
	if err != nil {
		return validator, err
	}

	return validator, nil
}
