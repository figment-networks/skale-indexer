package store

import (
	"context"
	"github.com/ethereum/go-ethereum/common"
	"math/big"

	"github.com/figment-networks/skale-indexer/scraper/structs"
)

//go:generate mockgen -destination=./mocks/mock_store.go  -package=mocks github.com/figment-networks/skale-indexer/store DataStore

type DBDriver interface {
	ContractEventStore
	SystemEventStore
	SkaleStore
}

type DataStore interface {
	ContractEventStore
	SystemEventStore
	SkaleStore
}

type SkaleStore interface {
	SaveNodes(ctx context.Context, nodes []structs.Node, removedNodeAddress common.Address) error
	GetNodes(ctx context.Context, params structs.NodeParams) (nodes []structs.Node, err error)

	SaveValidator(ctx context.Context, validator structs.Validator) error
	GetValidators(ctx context.Context, params structs.ValidatorParams) (validators []structs.Validator, err error)

	SaveDelegations(ctx context.Context, delegation []structs.Delegation) error
	GetDelegations(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)
	GetDelegationTimeline(ctx context.Context, params structs.DelegationParams) (delegations []structs.Delegation, err error)

	SaveAccount(ctx context.Context, account structs.Account) error
	GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error)

	SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error)

	GetValidatorStatistics(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error)
	CalculateTotalStake(ctx context.Context, params structs.ValidatorStatisticsParams) error
	CalculateActiveNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error
	CalculateLinkedNodes(ctx context.Context, params structs.ValidatorStatisticsParams) error
}

type ContractEventStore interface {
	SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error
	GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error)
}

type SystemEventStore interface {
	SaveSystemEvent(ctx context.Context, event structs.SystemEvent) error
	GetSystemEvents(ctx context.Context, params structs.SystemEventParams) (events []structs.SystemEvent, err error)
}

type Store struct {
	driver DBDriver
}

func New(driver DBDriver) *Store {
	return &Store{driver: driver}
}

// Skale objects

func (s *Store) SaveAccount(ctx context.Context, account structs.Account) error {
	return s.driver.SaveAccount(ctx, account)
}

func (s *Store) GetAccounts(ctx context.Context, params structs.AccountParams) (accounts []structs.Account, err error) {
	return s.driver.GetAccounts(ctx, params)
}

func (s *Store) SaveNodes(ctx context.Context, nodes []structs.Node, removedNodeAddress common.Address) error {
	return s.driver.SaveNodes(ctx, nodes, removedNodeAddress)
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

func (s *Store) SaveDelegations(ctx context.Context, delegations []structs.Delegation) error {
	return s.driver.SaveDelegations(ctx, delegations)
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

func (s *Store) GetValidatorStatisticsTimeline(ctx context.Context, params structs.ValidatorStatisticsParams) (validatorStatistics []structs.ValidatorStatistics, err error) {
	return s.driver.GetValidatorStatisticsTimeline(ctx, params)
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

func (s *Store) SaveValidatorStatistic(ctx context.Context, validatorID *big.Int, blockHeight uint64, statisticsType structs.StatisticTypeVS, amount *big.Int) (err error) {
	return s.driver.SaveValidatorStatistic(ctx, validatorID, blockHeight, statisticsType, amount)
}

// Contract events

func (s *Store) SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error {
	return s.driver.SaveContractEvent(ctx, contractEvent)
}

func (s *Store) GetContractEvents(ctx context.Context, params structs.EventParams) (contractEvents []structs.ContractEvent, err error) {
	return s.driver.GetContractEvents(ctx, params)
}

// System events

func (s *Store) SaveSystemEvent(ctx context.Context, event structs.SystemEvent) error {
	return s.driver.SaveSystemEvent(ctx, event)
}

func (s *Store) GetSystemEvents(ctx context.Context, params structs.SystemEventParams) (event []structs.SystemEvent, err error) {
	return s.driver.GetSystemEvents(ctx, params)
}
