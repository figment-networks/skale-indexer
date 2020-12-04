package store

import (
	"context"
	"github.com/figment-networks/skale-indexer/structs"
)

type DBDriver interface {
	DelegationStore
	EventStore
	ValidatorStore
	NodeStore
	DelegationStatisticsStore
}

type DataStore interface {
	DelegationStore
	EventStore
	ValidatorStore
	NodeStore
	DelegationStatisticsStore
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

type NodeStore interface {
	SaveOrUpdateNodes(ctx context.Context, nodes []structs.Node) error
	GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error)
}

type DelegationStatisticsStore interface {
	GetDelegationStatistics(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error)
	CalculateLatestDelegationStatesStatistics(ctx context.Context, params structs.QueryParams) error
	GetLatestDelegationStates(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error)
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

func (s *Store) SaveOrUpdateNodes(ctx context.Context, nodes []structs.Node) error {
	return s.driver.SaveOrUpdateNodes(ctx, nodes)
}

func (s *Store) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	return s.driver.GetNodes(ctx, params)
}

func (s *Store) GetDelegationStatistics(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error) {
	return s.driver.GetDelegationStatistics(ctx, params)
}

func (s *Store) CalculateLatestDelegationStatesStatistics(ctx context.Context, params structs.QueryParams) error {
	return s.driver.CalculateLatestDelegationStatesStatistics(ctx, params)
}

func (s *Store) GetLatestDelegationStates(ctx context.Context, params structs.QueryParams) (delegationStatistics []structs.DelegationStatistics, err error) {
	return s.driver.GetLatestDelegationStates(ctx, params)
}
