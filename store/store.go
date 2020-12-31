package store

import (
	"context"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

//go:generate mockgen -destination=./mocks/mock_store.go  -package=mocks github.com/figment-networks/skale-indexer/store DataStore

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

type DelegationStore interface {
	SaveDelegation(ctx context.Context, delegation structs.Delegation) error
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
}

type AccountStore interface {
	SaveAccount(ctx context.Context, account structs.Account) error
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)
}

type ValidatorStatisticsStore interface {
	SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error)

	GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetValidatorStatisticsChart(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error
	CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error
	CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error
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

func (s *Store) GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error) {
	return s.driver.GetDelegationTimeline(ctx, params)
}

func (s *Store) GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return s.driver.GetValidatorStatistics(ctx, params)
}

func (s *Store) GetValidatorStatisticsChart(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return s.driver.GetValidatorStatisticsChart(ctx, params)
}

func (s *Store) CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	return s.driver.CalculateTotalStake(ctx, params)
}

func (s *Store) CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	return s.driver.CalculateActiveNodes(ctx, params)
}

func (s *Store) CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error {
	return s.driver.CalculateLinkedNodes(ctx, params)
}

func (s *Store) SaveAccount(ctx context.Context, account structs.Account) error {
	return s.driver.SaveAccount(ctx, account)
}

func (s *Store) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	return s.driver.GetAccounts(ctx, params)
}

func (s *Store) SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error) {
	return s.driver.SaveValidatorStatistic(ctx, validatorID, blockHeight, statisticsType, amount)
}
