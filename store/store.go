package store

import (
	"../structs"
	"context"
)

type DBDriver interface {
	DelegationStore
	DelegationEventStore
	ValidatorStore
	ValidatorEventStore
}

type DataStore interface {
	DelegationStore
	DelegationEventStore
	ValidatorStore
	ValidatorEventStore
}

type DelegationStore interface {
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error)
	GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error)
	GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error)
}

type DelegationEventStore interface {
	SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error
	GetDelegationEventById(ctx context.Context, id string) (res structs.DelegationEvent, err error)
	GetDelegationEventsByDelegationId(ctx context.Context, delegationId string) (delegationEvents []structs.DelegationEvent, err error)
	GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error)
}

type ValidatorStore interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error)
	GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error)
	GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress string) (validators []structs.Validator, err error)
}

type ValidatorEventStore interface {
	SaveOrUpdateValidatorEvents(ctx context.Context, validatorEvents []structs.ValidatorEvent) error
	GetValidatorEventById(ctx context.Context, id string) (res structs.ValidatorEvent, err error)
	GetValidatorEventsByValidatorId(ctx context.Context, validatorId string) (validatorEvents []structs.ValidatorEvent, err error)
	GetAllValidatorEvents(ctx context.Context) (validatorEvents []structs.ValidatorEvent, err error)
}

type Store struct {
	driver DBDriver
}

func New(driver DBDriver) *Store {
	return &Store{driver: driver}
}

func (s *Store) SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error {
	return s.driver.SaveOrUpdateDelegations(ctx, delegations)
}

func (s *Store) GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error) {
	return s.driver.GetDelegationById(ctx, id)
}

func (s *Store) GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegationsByHolder(ctx, holder)
}

func (s *Store) GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegationsByValidatorId(ctx, validatorId)
}

func (s *Store) SaveOrUpdateDelegationEvents(ctx context.Context, delegationEvents []structs.DelegationEvent) error {
	return s.driver.SaveOrUpdateDelegationEvents(ctx, delegationEvents)
}

func (s *Store) GetDelegationEventById(ctx context.Context, id string) (res structs.DelegationEvent, err error) {
	return s.driver.GetDelegationEventById(ctx, id)
}

func (s *Store) GetDelegationEventsByDelegationId(ctx context.Context, delegationId string) (delegationEvents []structs.DelegationEvent, err error) {
	return s.driver.GetDelegationEventsByDelegationId(ctx, delegationId)
}

func (s *Store) GetAllDelegationEvents(ctx context.Context) (delegationEvents []structs.DelegationEvent, err error) {
	return s.driver.GetAllDelegationEvents(ctx)
}

func (s *Store) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	return s.driver.SaveOrUpdateValidators(ctx, validators)
}

func (s *Store) GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error) {
	return s.driver.GetValidatorById(ctx, id)
}

func (s *Store) GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error) {
	return s.driver.GetValidatorsByAddress(ctx, validatorAddress)
}

func (s *Store) GetValidatorsByRequestedAddress(ctx context.Context, validatorId string) (validators []structs.Validator, err error) {
	return s.driver.GetValidatorsByRequestedAddress(ctx, validatorId)
}

func (s *Store) SaveOrUpdateValidatorEvents(ctx context.Context, validatorEvents []structs.ValidatorEvent) error {
	return s.driver.SaveOrUpdateValidatorEvents(ctx, validatorEvents)
}

func (s *Store) GetValidatorEventById(ctx context.Context, id string) (res structs.ValidatorEvent, err error) {
	return s.driver.GetValidatorEventById(ctx, id)
}

func (s *Store) GetValidatorEventsByValidatorId(ctx context.Context, validatorId string) (delegationEvents []structs.ValidatorEvent, err error) {
	return s.driver.GetValidatorEventsByValidatorId(ctx, validatorId)
}

func (s *Store) GetAllValidatorEvents(ctx context.Context) (validatorEvents []structs.ValidatorEvent, err error) {
	return s.driver.GetAllValidatorEvents(ctx)
}
