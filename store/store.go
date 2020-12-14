package store

import (
	"context"
	"math/big"

	"github.com/figment-networks/skale-indexer/client/structs"
)

type DBDriver interface {
	ContractEventStore
	NodeStore
	ValidatorStore
	DelegationStore
	ValidatorStatisticsStore
}

type DataStore interface {
	ContractEventStore
	NodeStore
	ValidatorStore
	DelegationStore
	ValidatorStatisticsStore
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

type DelegationStore interface {
	SaveDelegation(ctx context.Context, delegation structs.Delegation) error
	GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error)
}

type ValidatorStatisticsStore interface {
	GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	CalculateParams(ctx context.Context, height uint64, vID *big.Int) error
	CalculateTotalStake(ctx context.Context, params structs.QueryParams) error
	CalculateActiveNodes(ctx context.Context, params structs.QueryParams) error
	CalculateLinkedNodes(ctx context.Context, params structs.QueryParams) error
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

func (s *Store) SaveDelegation(ctx context.Context, delegation structs.Delegation) error {
	return s.driver.SaveDelegation(ctx, delegation)
}

func (s *Store) GetDelegations(ctx context.Context, params structs.QueryParams) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegations(ctx, params)
}

func (s *Store) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return s.driver.GetValidatorStatistics(ctx, params)
}

func (s *Store) CalculateParams(ctx context.Context, blockHeight uint64, validatorId *big.Int) error {
	params := structs.QueryParams{
		ValidatorId:    validatorId,
		ETHBlockHeight: blockHeight,
	}
	// TODO: add transactional commit-rollback
	if err := s.driver.CalculateTotalStake(ctx, params); err != nil {
		return err
	}
	if err := s.driver.CalculateActiveNodes(ctx, params); err != nil {
		return err
	}
	if err := s.driver.CalculateLinkedNodes(ctx, params); err != nil {
		return err
	}

	return nil
}

func (s *Store) CalculateTotalStake(ctx context.Context, params structs.QueryParams) error {
	return s.driver.CalculateTotalStake(ctx, params)
}

func (s *Store) CalculateActiveNodes(ctx context.Context, params structs.QueryParams) error {
	return s.driver.CalculateActiveNodes(ctx, params)
}

func (s *Store) CalculateLinkedNodes(ctx context.Context, params structs.QueryParams) error {
	return s.driver.CalculateLinkedNodes(ctx, params)
}
