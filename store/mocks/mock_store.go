// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/figment-networks/skale-indexer/store (interfaces: DataStore)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	common "github.com/ethereum/go-ethereum/common"
	structs "github.com/figment-networks/skale-indexer/scraper/structs"
	gomock "github.com/golang/mock/gomock"
	big "math/big"
	reflect "reflect"
	time "time"
)

// MockDataStore is a mock of DataStore interface
type MockDataStore struct {
	ctrl     *gomock.Controller
	recorder *MockDataStoreMockRecorder
}

// MockDataStoreMockRecorder is the mock recorder for MockDataStore
type MockDataStoreMockRecorder struct {
	mock *MockDataStore
}

// NewMockDataStore creates a new mock instance
func NewMockDataStore(ctrl *gomock.Controller) *MockDataStore {
	mock := &MockDataStore{ctrl: ctrl}
	mock.recorder = &MockDataStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDataStore) EXPECT() *MockDataStoreMockRecorder {
	return m.recorder
}

// GetAccounts mocks base method
func (m *MockDataStore) GetAccounts(arg0 context.Context, arg1 structs.AccountParams) ([]structs.Account, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", arg0, arg1)
	ret0, _ := ret[0].([]structs.Account)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccounts indicates an expected call of GetAccounts
func (mr *MockDataStoreMockRecorder) GetAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockDataStore)(nil).GetAccounts), arg0, arg1)
}

// GetContractEvents mocks base method
func (m *MockDataStore) GetContractEvents(arg0 context.Context, arg1 structs.EventParams) ([]structs.ContractEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractEvents", arg0, arg1)
	ret0, _ := ret[0].([]structs.ContractEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractEvents indicates an expected call of GetContractEvents
func (mr *MockDataStoreMockRecorder) GetContractEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractEvents", reflect.TypeOf((*MockDataStore)(nil).GetContractEvents), arg0, arg1)
}

// GetDelegationTimeline mocks base method
func (m *MockDataStore) GetDelegationTimeline(arg0 context.Context, arg1 structs.DelegationParams) ([]structs.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegationTimeline", arg0, arg1)
	ret0, _ := ret[0].([]structs.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegationTimeline indicates an expected call of GetDelegationTimeline
func (mr *MockDataStoreMockRecorder) GetDelegationTimeline(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegationTimeline", reflect.TypeOf((*MockDataStore)(nil).GetDelegationTimeline), arg0, arg1)
}

// GetDelegations mocks base method
func (m *MockDataStore) GetDelegations(arg0 context.Context, arg1 structs.DelegationParams) ([]structs.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegations", arg0, arg1)
	ret0, _ := ret[0].([]structs.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegations indicates an expected call of GetDelegations
func (mr *MockDataStoreMockRecorder) GetDelegations(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegations", reflect.TypeOf((*MockDataStore)(nil).GetDelegations), arg0, arg1)
}

// GetNodes mocks base method
func (m *MockDataStore) GetNodes(arg0 context.Context, arg1 structs.NodeParams) ([]structs.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes", arg0, arg1)
	ret0, _ := ret[0].([]structs.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockDataStoreMockRecorder) GetNodes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockDataStore)(nil).GetNodes), arg0, arg1)
}

// GetSystemEvents mocks base method
func (m *MockDataStore) GetSystemEvents(arg0 context.Context, arg1 structs.SystemEventParams) ([]structs.SystemEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSystemEvents", arg0, arg1)
	ret0, _ := ret[0].([]structs.SystemEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSystemEvents indicates an expected call of GetSystemEvents
func (mr *MockDataStoreMockRecorder) GetSystemEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSystemEvents", reflect.TypeOf((*MockDataStore)(nil).GetSystemEvents), arg0, arg1)
}

// GetValidatorStatistics mocks base method
func (m *MockDataStore) GetValidatorStatistics(arg0 context.Context, arg1 structs.ValidatorStatisticsParams) ([]structs.ValidatorStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorStatistics", arg0, arg1)
	ret0, _ := ret[0].([]structs.ValidatorStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorStatistics indicates an expected call of GetValidatorStatistics
func (mr *MockDataStoreMockRecorder) GetValidatorStatistics(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorStatistics", reflect.TypeOf((*MockDataStore)(nil).GetValidatorStatistics), arg0, arg1)
}

// GetValidatorStatisticsTimeline mocks base method
func (m *MockDataStore) GetValidatorStatisticsTimeline(arg0 context.Context, arg1 structs.ValidatorStatisticsParams) ([]structs.ValidatorStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorStatisticsTimeline", arg0, arg1)
	ret0, _ := ret[0].([]structs.ValidatorStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorStatisticsTimeline indicates an expected call of GetValidatorStatisticsTimeline
func (mr *MockDataStoreMockRecorder) GetValidatorStatisticsTimeline(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorStatisticsTimeline", reflect.TypeOf((*MockDataStore)(nil).GetValidatorStatisticsTimeline), arg0, arg1)
}

// GetValidators mocks base method
func (m *MockDataStore) GetValidators(arg0 context.Context, arg1 structs.ValidatorParams) ([]structs.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", arg0, arg1)
	ret0, _ := ret[0].([]structs.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators
func (mr *MockDataStoreMockRecorder) GetValidators(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockDataStore)(nil).GetValidators), arg0, arg1)
}

// SaveAccount mocks base method
func (m *MockDataStore) SaveAccount(arg0 context.Context, arg1 structs.Account) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAccount indicates an expected call of SaveAccount
func (mr *MockDataStoreMockRecorder) SaveAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAccount", reflect.TypeOf((*MockDataStore)(nil).SaveAccount), arg0, arg1)
}

// SaveContractEvent mocks base method
func (m *MockDataStore) SaveContractEvent(arg0 context.Context, arg1 structs.ContractEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveContractEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveContractEvent indicates an expected call of SaveContractEvent
func (mr *MockDataStoreMockRecorder) SaveContractEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContractEvent", reflect.TypeOf((*MockDataStore)(nil).SaveContractEvent), arg0, arg1)
}

// SaveDelegation mocks base method
func (m *MockDataStore) SaveDelegation(arg0 context.Context, arg1 structs.Delegation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDelegation", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDelegation indicates an expected call of SaveDelegation
func (mr *MockDataStoreMockRecorder) SaveDelegation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDelegation", reflect.TypeOf((*MockDataStore)(nil).SaveDelegation), arg0, arg1)
}

// SaveNodes mocks base method
func (m *MockDataStore) SaveNodes(arg0 context.Context, arg1 []structs.Node, arg2 common.Address) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNodes", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveNodes indicates an expected call of SaveNodes
func (mr *MockDataStoreMockRecorder) SaveNodes(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNodes", reflect.TypeOf((*MockDataStore)(nil).SaveNodes), arg0, arg1, arg2)
}

// SaveSystemEvent mocks base method
func (m *MockDataStore) SaveSystemEvent(arg0 context.Context, arg1 structs.SystemEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveSystemEvent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveSystemEvent indicates an expected call of SaveSystemEvent
func (mr *MockDataStoreMockRecorder) SaveSystemEvent(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveSystemEvent", reflect.TypeOf((*MockDataStore)(nil).SaveSystemEvent), arg0, arg1)
}

// SaveValidator mocks base method
func (m *MockDataStore) SaveValidator(arg0 context.Context, arg1 structs.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveValidator", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveValidator indicates an expected call of SaveValidator
func (mr *MockDataStoreMockRecorder) SaveValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValidator", reflect.TypeOf((*MockDataStore)(nil).SaveValidator), arg0, arg1)
}

// SaveValidatorStatistic mocks base method
func (m *MockDataStore) SaveValidatorStatistic(arg0 context.Context, arg1 *big.Int, arg2 uint64, arg3 time.Time, arg4 structs.StatisticTypeVS, arg5 *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveValidatorStatistic", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveValidatorStatistic indicates an expected call of SaveValidatorStatistic
func (mr *MockDataStoreMockRecorder) SaveValidatorStatistic(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValidatorStatistic", reflect.TypeOf((*MockDataStore)(nil).SaveValidatorStatistic), arg0, arg1, arg2, arg3, arg4, arg5)
}

// UpdateCountsOfValidator mocks base method
func (m *MockDataStore) UpdateCountsOfValidator(arg0 context.Context, arg1 *big.Int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCountsOfValidator", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCountsOfValidator indicates an expected call of UpdateCountsOfValidator
func (mr *MockDataStoreMockRecorder) UpdateCountsOfValidator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCountsOfValidator", reflect.TypeOf((*MockDataStore)(nil).UpdateCountsOfValidator), arg0, arg1)
}
