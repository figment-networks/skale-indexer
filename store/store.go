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
	GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error)
}

type EventStore interface {
	SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error
	GetEvents(ctx context.Context, params structs.QueryParams) (events []structs.Event, err error)
}

type ValidatorStore interface {
	SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error
	GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error)
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

func (s *Store) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegations(ctx, params)
}

func (s *Store) SaveOrUpdateEvents(ctx context.Context, events []structs.Event) error {
	return s.driver.SaveOrUpdateEvents(ctx, events)
}

func (s *Store) GetEvents(ctx context.Context, params structs.QueryParams) (events []structs.Event, err error) {
	return s.driver.GetEvents(ctx, params)
}

func (s *Store) SaveOrUpdateValidators(ctx context.Context, validators []structs.Validator) error {
	return s.driver.SaveOrUpdateValidators(ctx, validators)
}

func (s *Store) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	return s.driver.GetValidators(ctx, params)
}
