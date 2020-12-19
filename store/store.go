package store

import (
	"context"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

type DBDriver interface {
	ContractEventStore
	NodeStore
	ValidatorStore
	DelegationStore
	ValidatorStatisticsStore
	AccountStore
}

type DataStore interface {
	ContractEventStore
	NodeStore
	ValidatorStore
	DelegationStore
	ValidatorStatisticsStore
	AccountStore
}

type ContractEventStore interface {
	SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error
	GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error)
}

type NodeStore interface {
	SaveNode(ctx context.Context, node structs.Node) error
	GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error)
}

type ValidatorStore interface {
	SaveValidator(ctx context.Context, validator structs.Validator) error
	GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error)
}

type AccountStore interface {
	SaveAccount(ctx context.Context, account structs.Account) error
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)
}

type ValidatorStatisticsStore interface {
	GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	//	CalculateParams(ctx context.Context, height uint64, vID *big.Int) error
	//  CalculateTotalStake(ctx context.Context, params structs.QueryParams) error
	//	CalculateActiveNodes(ctx context.Context, params structs.QueryParams) error
	//	CalculateLinkedNodes(ctx context.Context, params structs.QueryParams) error
}

type DelegationStore interface {
	SaveDelegation(ctx context.Context, delegation structs.Delegation) error
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
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

func (s *Store) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {
	return s.driver.GetContractEvents(ctx, params)
}

func (s *Store) SaveNode(ctx context.Context, node structs.Node) error {
	return s.driver.SaveNode(ctx, node)
}

func (s *Store) GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error) {
	return s.driver.GetNodes(ctx, params)
}

func (s *Store) SaveValidator(ctx context.Context, validator structs.Validator) error {
	return s.driver.SaveValidator(ctx, validator)
}

func (s *Store) GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error) {
	return s.driver.GetValidators(ctx, params)
}

func (s *Store) SaveDelegation(ctx context.Context, delegation structs.Delegation) error {
	return s.driver.SaveDelegation(ctx, delegation)
}

func (s *Store) GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegations(ctx, params)
}

func (s *Store) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return s.driver.GetValidatorStatistics(ctx, params)
}

/*
func (s *Store) CalculateParams(ctx context.Context, blockHeight uint64, validatorId *big.Int) error {
	params := structs.QueryParams{
		ValidatorId:    validatorId.Uint64(),
		ETHBlockHeight: blockHeight,
	}
	//TODO: add transactional commit-rollback
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
*/

func (s *Store) SaveAccount(ctx context.Context, account structs.Account) error {
	return s.driver.SaveAccount(ctx, account)
}

func (s *Store) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	return s.driver.GetAccounts(ctx, params)
}