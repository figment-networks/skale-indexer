package client

import (
	"context"
	"errors"
	"github.com/figment-networks/skale-indexer/store"
	"github.com/figment-networks/skale-indexer/structs"
)

var (
	InvalidId = errors.New("invalid id")
)

type ClientContractor struct {
	storeEng store.DataStore
}

func NewClientContractor(storeEng store.DataStore) *ClientContractor {
	return &ClientContractor{
		storeEng: storeEng,
	}
}

func (c *ClientContractor) SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error {
	return c.storeEng.SaveOrUpdateDelegations(ctx, delegations)
}

func (c *ClientContractor) GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error) {
	if !(id != "") {
		return res, InvalidId
	}
	return c.storeEng.GetDelegationById(ctx, id)
}

func (c *ClientContractor) GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegationsByHolder(ctx, holder)
}

func (c *ClientContractor) GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error) {
	return c.storeEng.GetDelegationsByValidatorId(ctx, validatorId)
}

func (c *ClientContractor) SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error {
	return c.storeEng.SaveOrUpdateDelegationEvents(ctx, delegationEvents)
}

func (c *ClientContractor) GetDelegationEventById(ctx context.Context, id string) (res structs.DelegationEvent, err error) {
	return c.storeEng.GetDelegationEventById(ctx, id)
}

func (c *ClientContractor) GetDelegationEventsByDelegationId(ctx context.Context, delegationId string) (delegationEvents []structs.DelegationEvent, err error) {
	return c.storeEng.GetDelegationEventsByDelegationId(ctx, delegationId)
}

func (c *ClientContractor) GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error) {
	return c.storeEng.GetAllDelegationEvents(ctx)
}

func (c *ClientContractor) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	return c.storeEng.SaveOrUpdateValidators(ctx, validators)
}

func (c *ClientContractor) GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error) {
	if !(id != "") {
		return res, InvalidId
	}
	return c.storeEng.GetValidatorById(ctx, id)
}

func (c *ClientContractor) GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidatorsByAddress(ctx, validatorAddress)
}

func (c *ClientContractor) GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress string) (validators []structs.Validator, err error) {
	return c.storeEng.GetValidatorsByRequestedAddress(ctx, requestedAddress)
}
