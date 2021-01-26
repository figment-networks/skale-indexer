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

func (c *Caller) GetValidator(ctx context.Context, bc *bind.BoundContract, blockNumber uint64, validatorID *big.Int) (v structs.Validator, err error) {

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

	err = bc.Call(co, &results, "getValidator", validatorID)

	if err != nil {
		return v, fmt.Errorf("error calling getValidator function %w", err)
	}

	if len(results) == 0 {
		return v, errors.New("empty result")
	}

	vr := &ValidatorRaw{}
	vraw := *abi.ConvertType(results[0], vr).(*ValidatorRaw)

	return structs.Validator{
		ValidatorID:             validatorID,
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

	err = bc.Call(&bind.CallOpts{
		Pending: false,
		Context: ctxT,
	}, &results, "isAuthorizedValidator", validatorID)

	if err != nil {
		return false, fmt.Errorf("error calling getValidator function %w", err)
	}

	var ok bool
	isAuthorized, ok = results[0].(bool)
	if !ok {
		return false, errors.New("earned is not bool type ")
	}

	return isAuthorized, nil
}

/* gets 10 validators based on ind parameter
 * to be used for synchronization
 *
 * example: if ind is 5, then it will fetch validators for validator_id between 41 and 50
 */
func (c *Caller) FetchNextRoundValidators(ctx context.Context, bc *bind.BoundContract, ind int64, currentBlock uint64) (validators []structs.Validator, err error) {
	validators = []structs.Validator{}
	length := int64(10)
	vldID := (ind-1)*length + 1
	for i := 0; i < int(length); i++ {
		vldIDBig := big.NewInt(vldID)

		v, err := c.GetValidator(ctx, bc, currentBlock, vldIDBig)
		if err != nil {
			break
		}
		v.Authorized, err = c.IsAuthorizedValidator(ctx, bc, currentBlock, vldIDBig)
		if err != nil {
			break
		}

		v.BlockHeight = currentBlock
		vldID++
		validators = append(validators, v)
	}

	return validators, nil
}
