package store

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type DBDriver interface {
	DelegationStore
	EventStore
	ValidatorStore
}

type DataStore interface {
	DelegationStore
	EventStore
	ValidatorStore
}

type DelegationStore interface {
	SaveOrUpdateDelegations(ctx context.Context, delegations []structs.Delegation) error
	GetDelegationById(ctx context.Context, id string) (res structs.Delegation, err error)
	GetDelegationsByHolder(ctx context.Context, holder string) (delegations []structs.Delegation, err error)
	GetDelegationsByValidatorId(ctx context.Context, validatorId uint64) (delegations []structs.Delegation, err error)
}

type EventStore interface {
	SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error
	GetEventById(ctx context.Context, id string) (res structs.Event, err error)
	GetAllEvents(ctx context.Context) (events []structs.Event, err error)
}

type ValidatorStore interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidatorById(ctx context.Context, id string) (res structs.Validator, err error)
	GetValidatorsByAddress(ctx context.Context, validatorAddress string) (validators []structs.Validator, err error)
	GetValidatorsByRequestedAddress(ctx context.Context, requestedAddress string) (validators []structs.Validator, err error)
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

func (s *Store) SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error {
	return s.driver.SaveOrUpdateEvents(ctx, events)
}

func (s *Store) GetEventById(ctx context.Context, id string) (res structs.Event, err error) {
	return s.driver.GetEventById(ctx, id)
}

func (s *Store) GetAllEvents(ctx context.Context) (events []structs.Event, err error) {
	return s.driver.GetAllEvents(ctx)
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
