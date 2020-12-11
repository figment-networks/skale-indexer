package store

import (
	"context"

	"github.com/figment-networks/skale-indexer/structs"
)

type DBDriver interface {
	ContractEventStore
	NodeStore
	ValidatorStore
}

type DataStore interface {
	ContractEventStore
	NodeStore
	ValidatorStore
}

type ContractEventStore interface {
	SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error
	GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error)
}

type NodeStore interface {
	SaveNode(ctx context.Context, node structs.Node) error
	GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error)
}

type ValidatorStore interface {
	SaveValidator(ctx context.Context, validator structs.Validator) error
	GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error)
}

type Store struct {
	driver DBDriver
}

func New(driver DBDriver) *Store {
	return &Store{driver: driver}
}

func (s *Store) SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error {
	return s.driver.SaveContractEvent(ctx, contractEvent)
}

func (s *Store) GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error) {
	return s.driver.GetContractEvents(ctx, params)
}

func (s *Store) SaveNode(ctx context.Context, node structs.Node) error {
	return s.driver.SaveNode(ctx, node)
}

func (s *Store) GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error) {
	return s.driver.GetNodes(ctx, params)
}

func (s *Store) SaveValidator(ctx context.Context, validator structs.Validator) error {
	return s.driver.SaveValidator(ctx, validator)
}

func (s *Store) GetValidators(ctx context.Context, params structs.QueryParams) (validators []structs.Validator, err error) {
	return s.driver.GetValidators(ctx, params)
}
