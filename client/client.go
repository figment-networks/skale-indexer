package client

import (
	"../store"
	"../structs"
	"../types"
	"context"
	"errors"

	//"go.uber.org/zap"
)

var (
	InvalidId = errors.New("invalid id")
)

type ClientContractor struct {
	storeEng store.DataStore
	//logger   *zap.Logger
}

//func NewClient(storeEng store.DataStore, logger *zap.Logger) *Client {
func NewClientContractor(storeEng store.DataStore) *ClientContractor {
	return &ClientContractor{
		storeEng: storeEng,
		//logger:   logger,
	}
}

func (c *ClientContractor) SaveOrUpdateDelegation(ctx context.Context, delegation structs.Delegation) error {
	defer c.recoverPanic()
	return c.storeEng.SaveOrUpdateDelegation(ctx, delegation)
}

func (c *ClientContractor) SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error {
	defer c.recoverPanic()
	return c.storeEng.SaveOrUpdateDelegations(ctx, delegations)
}

func (c *ClientContractor) GetDelegationById(ctx context.Context, id *types.ID) (res structs.Delegation, err error) {
	defer c.recoverPanic()
	if !id.Valid() {
		return res, InvalidId
	}
	return c.storeEng.GetDelegationById(ctx, id)
}

func (c *ClientContractor) GetDelegationsByHolder(ctx context.Context, holder *string) (delegations []structs.Delegation, err error) {
	defer c.recoverPanic()
	return c.storeEng.GetDelegationsByHolder(ctx, holder)
}

func (c *ClientContractor) GetDelegationsByValidatorId(ctx context.Context, validatorId *uint64) (delegations []structs.Delegation, err error) {
	defer c.recoverPanic()
	return c.storeEng.GetDelegationsByValidatorId(ctx, validatorId)
}

func (c *ClientContractor) SaveOrUpdateValidator(ctx context.Context, validator structs.Validator) error {
	defer c.recoverPanic()
	return c.storeEng.SaveOrUpdateValidator(ctx, validator)
}

func (c *ClientContractor) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	defer c.recoverPanic()
	return c.storeEng.SaveOrUpdateValidators(ctx, validators)
}

func (c *ClientContractor) GetValidatorById(ctx context.Context, id *types.ID) (res structs.Validator, err error) {
	defer c.recoverPanic()
	if !id.Valid() {
		return res, InvalidId
	}
	return c.storeEng.GetValidatorById(ctx, id)
}

func (c *ClientContractor) GetValidatorsByValidatorAddress(ctx context.Context, validatorAddress *string) (validators []structs.Validator, err error) {
	defer c.recoverPanic()
	return c.storeEng.GetValidatorsByValidatorAddress(ctx, validatorAddress)
}

func (c *ClientContractor) GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress *string) (validators []structs.Validator, err error) {
	defer c.recoverPanic()
	return c.storeEng.GetValidatorsByRequestedAddress(ctx, requestedAddress)
}

func (c *ClientContractor) recoverPanic() {
	if p := recover(); p != nil {
		//c.logger.Error("[Client] Panic ", zap.Any("contents", p))
		//c.logger.Sync()
	}
}
