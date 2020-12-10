package store

import (
	"context"

	"github.com/figment-networks/skale-indexer/structs"
)

type DBDriver interface {
	ContractEventStore
}

type DataStore interface {
	ContractEventStore
	NodeStore
}

type ContractEventStore interface {
	SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error
	GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error)
}

type NodeStore interface {
	SaveNode(ctx context.Context, node structs.Node) error
	GetNodes(ctx context.Context, params structs.QueryParams) (nodes []structs.Node, err error)
}

type Store struct {
	driver DBDriver
}

func New(driver DBDriver) *Store {
	return &Store{driver: driver}
}

func (s *Store) SaveContractEvent(ctx context.Context, ce structs.ContractEvent) error {
	return s.driver.SaveContractEvent(ctx, ce)
}

func (s *Store) GetContractEvents(ctx context.Context, params structs.QueryParams) (contractEvents []structs.ContractEvent, err error) {
	return s.driver.GetContractEvents(ctx, params)
}
