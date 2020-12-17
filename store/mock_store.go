// Code generated by MockGen. DO NOT EDIT.
// Source: store/store.go

// Package store is a generated GoMock package.
package store

import (
	context "context"
	structs "github.com/figment-networks/skale-indexer/scraper/structs"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDBDriver is a mock of DBDriver interface
type MockDBDriver struct {
	ctrl     *gomock.Controller
	recorder *MockDBDriverMockRecorder
}

// MockDBDriverMockRecorder is the mock recorder for MockDBDriver
type MockDBDriverMockRecorder struct {
	mock *MockDBDriver
}

// NewMockDBDriver creates a new mock instance
func NewMockDBDriver(ctrl *gomock.Controller) *MockDBDriver {
	mock := &MockDBDriver{ctrl: ctrl}
	mock.recorder = &MockDBDriverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDBDriver) EXPECT() *MockDBDriverMockRecorder {
	return m.recorder
}

// SaveContractEvent mocks base method
func (m *MockDBDriver) SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveContractEvent", ctx, contractEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveContractEvent indicates an expected call of SaveContractEvent
func (mr *MockDBDriverMockRecorder) SaveContractEvent(ctx, contractEvent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContractEvent", reflect.TypeOf((*MockDBDriver)(nil).SaveContractEvent), ctx, contractEvent)
}

// GetContractEvents mocks base method
func (m *MockDBDriver) GetContractEvents(ctx context.Context, params structs.EventParams) ([]structs.ContractEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractEvents", ctx, params)
	ret0, _ := ret[0].([]structs.ContractEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractEvents indicates an expected call of GetContractEvents
func (mr *MockDBDriverMockRecorder) GetContractEvents(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractEvents", reflect.TypeOf((*MockDBDriver)(nil).GetContractEvents), ctx, params)
}

// SaveNode mocks base method
func (m *MockDBDriver) SaveNode(ctx context.Context, node structs.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNode", ctx, node)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveNode indicates an expected call of SaveNode
func (mr *MockDBDriverMockRecorder) SaveNode(ctx, node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNode", reflect.TypeOf((*MockDBDriver)(nil).SaveNode), ctx, node)
}

// GetNodes mocks base method
func (m *MockDBDriver) GetNodes(ctx context.Context, params structs.NodeParams) ([]structs.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes", ctx, params)
	ret0, _ := ret[0].([]structs.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockDBDriverMockRecorder) GetNodes(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockDBDriver)(nil).GetNodes), ctx, params)
}

// SaveValidator mocks base method
func (m *MockDBDriver) SaveValidator(ctx context.Context, validator structs.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveValidator", ctx, validator)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveValidator indicates an expected call of SaveValidator
func (mr *MockDBDriverMockRecorder) SaveValidator(ctx, validator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValidator", reflect.TypeOf((*MockDBDriver)(nil).SaveValidator), ctx, validator)
}

// GetValidators mocks base method
func (m *MockDBDriver) GetValidators(ctx context.Context, params structs.QueryParams) ([]structs.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, params)
	ret0, _ := ret[0].([]structs.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators
func (mr *MockDBDriverMockRecorder) GetValidators(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockDBDriver)(nil).GetValidators), ctx, params)
}

// SaveDelegation mocks base method
func (m *MockDBDriver) SaveDelegation(ctx context.Context, delegation structs.Delegation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDelegation", ctx, delegation)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDelegation indicates an expected call of SaveDelegation
func (mr *MockDBDriverMockRecorder) SaveDelegation(ctx, delegation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDelegation", reflect.TypeOf((*MockDBDriver)(nil).SaveDelegation), ctx, delegation)
}

// GetDelegations mocks base method
func (m *MockDBDriver) GetDelegations(ctx context.Context, params structs.DelegationParams) ([]structs.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegations", ctx, params)
	ret0, _ := ret[0].([]structs.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegations indicates an expected call of GetDelegations
func (mr *MockDBDriverMockRecorder) GetDelegations(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegations", reflect.TypeOf((*MockDBDriver)(nil).GetDelegations), ctx, params)
}

// GetValidatorStatistics mocks base method
func (m *MockDBDriver) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) ([]structs.ValidatorStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorStatistics", ctx, params)
	ret0, _ := ret[0].([]structs.ValidatorStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorStatistics indicates an expected call of GetValidatorStatistics
func (mr *MockDBDriverMockRecorder) GetValidatorStatistics(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorStatistics", reflect.TypeOf((*MockDBDriver)(nil).GetValidatorStatistics), ctx, params)
}

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

// SaveContractEvent mocks base method
func (m *MockDataStore) SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveContractEvent", ctx, contractEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveContractEvent indicates an expected call of SaveContractEvent
func (mr *MockDataStoreMockRecorder) SaveContractEvent(ctx, contractEvent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContractEvent", reflect.TypeOf((*MockDataStore)(nil).SaveContractEvent), ctx, contractEvent)
}

// GetContractEvents mocks base method
func (m *MockDataStore) GetContractEvents(ctx context.Context, params structs.EventParams) ([]structs.ContractEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractEvents", ctx, params)
	ret0, _ := ret[0].([]structs.ContractEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractEvents indicates an expected call of GetContractEvents
func (mr *MockDataStoreMockRecorder) GetContractEvents(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractEvents", reflect.TypeOf((*MockDataStore)(nil).GetContractEvents), ctx, params)
}

// SaveNode mocks base method
func (m *MockDataStore) SaveNode(ctx context.Context, node structs.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNode", ctx, node)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveNode indicates an expected call of SaveNode
func (mr *MockDataStoreMockRecorder) SaveNode(ctx, node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNode", reflect.TypeOf((*MockDataStore)(nil).SaveNode), ctx, node)
}

// GetNodes mocks base method
func (m *MockDataStore) GetNodes(ctx context.Context, params structs.NodeParams) ([]structs.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes", ctx, params)
	ret0, _ := ret[0].([]structs.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockDataStoreMockRecorder) GetNodes(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockDataStore)(nil).GetNodes), ctx, params)
}

// SaveValidator mocks base method
func (m *MockDataStore) SaveValidator(ctx context.Context, validator structs.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveValidator", ctx, validator)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveValidator indicates an expected call of SaveValidator
func (mr *MockDataStoreMockRecorder) SaveValidator(ctx, validator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValidator", reflect.TypeOf((*MockDataStore)(nil).SaveValidator), ctx, validator)
}

// GetValidators mocks base method
func (m *MockDataStore) GetValidators(ctx context.Context, params structs.QueryParams) ([]structs.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, params)
	ret0, _ := ret[0].([]structs.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators
func (mr *MockDataStoreMockRecorder) GetValidators(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockDataStore)(nil).GetValidators), ctx, params)
}

// SaveDelegation mocks base method
func (m *MockDataStore) SaveDelegation(ctx context.Context, delegation structs.Delegation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDelegation", ctx, delegation)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDelegation indicates an expected call of SaveDelegation
func (mr *MockDataStoreMockRecorder) SaveDelegation(ctx, delegation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDelegation", reflect.TypeOf((*MockDataStore)(nil).SaveDelegation), ctx, delegation)
}

// GetDelegations mocks base method
func (m *MockDataStore) GetDelegations(ctx context.Context, params structs.DelegationParams) ([]structs.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegations", ctx, params)
	ret0, _ := ret[0].([]structs.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegations indicates an expected call of GetDelegations
func (mr *MockDataStoreMockRecorder) GetDelegations(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegations", reflect.TypeOf((*MockDataStore)(nil).GetDelegations), ctx, params)
}

// GetValidatorStatistics mocks base method
func (m *MockDataStore) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) ([]structs.ValidatorStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorStatistics", ctx, params)
	ret0, _ := ret[0].([]structs.ValidatorStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorStatistics indicates an expected call of GetValidatorStatistics
func (mr *MockDataStoreMockRecorder) GetValidatorStatistics(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorStatistics", reflect.TypeOf((*MockDataStore)(nil).GetValidatorStatistics), ctx, params)
}

// MockContractEventStore is a mock of ContractEventStore interface
type MockContractEventStore struct {
	ctrl     *gomock.Controller
	recorder *MockContractEventStoreMockRecorder
}

// MockContractEventStoreMockRecorder is the mock recorder for MockContractEventStore
type MockContractEventStoreMockRecorder struct {
	mock *MockContractEventStore
}

// NewMockContractEventStore creates a new mock instance
func NewMockContractEventStore(ctrl *gomock.Controller) *MockContractEventStore {
	mock := &MockContractEventStore{ctrl: ctrl}
	mock.recorder = &MockContractEventStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockContractEventStore) EXPECT() *MockContractEventStoreMockRecorder {
	return m.recorder
}

// SaveContractEvent mocks base method
func (m *MockContractEventStore) SaveContractEvent(ctx context.Context, contractEvent structs.ContractEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveContractEvent", ctx, contractEvent)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveContractEvent indicates an expected call of SaveContractEvent
func (mr *MockContractEventStoreMockRecorder) SaveContractEvent(ctx, contractEvent interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveContractEvent", reflect.TypeOf((*MockContractEventStore)(nil).SaveContractEvent), ctx, contractEvent)
}

// GetContractEvents mocks base method
func (m *MockContractEventStore) GetContractEvents(ctx context.Context, params structs.EventParams) ([]structs.ContractEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractEvents", ctx, params)
	ret0, _ := ret[0].([]structs.ContractEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractEvents indicates an expected call of GetContractEvents
func (mr *MockContractEventStoreMockRecorder) GetContractEvents(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractEvents", reflect.TypeOf((*MockContractEventStore)(nil).GetContractEvents), ctx, params)
}

// MockNodeStore is a mock of NodeStore interface
type MockNodeStore struct {
	ctrl     *gomock.Controller
	recorder *MockNodeStoreMockRecorder
}

// MockNodeStoreMockRecorder is the mock recorder for MockNodeStore
type MockNodeStoreMockRecorder struct {
	mock *MockNodeStore
}

// NewMockNodeStore creates a new mock instance
func NewMockNodeStore(ctrl *gomock.Controller) *MockNodeStore {
	mock := &MockNodeStore{ctrl: ctrl}
	mock.recorder = &MockNodeStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockNodeStore) EXPECT() *MockNodeStoreMockRecorder {
	return m.recorder
}

// SaveNode mocks base method
func (m *MockNodeStore) SaveNode(ctx context.Context, node structs.Node) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveNode", ctx, node)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveNode indicates an expected call of SaveNode
func (mr *MockNodeStoreMockRecorder) SaveNode(ctx, node interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveNode", reflect.TypeOf((*MockNodeStore)(nil).SaveNode), ctx, node)
}

// GetNodes mocks base method
func (m *MockNodeStore) GetNodes(ctx context.Context, params structs.NodeParams) ([]structs.Node, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodes", ctx, params)
	ret0, _ := ret[0].([]structs.Node)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodes indicates an expected call of GetNodes
func (mr *MockNodeStoreMockRecorder) GetNodes(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodes", reflect.TypeOf((*MockNodeStore)(nil).GetNodes), ctx, params)
}

// MockValidatorStore is a mock of ValidatorStore interface
type MockValidatorStore struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorStoreMockRecorder
}

// MockValidatorStoreMockRecorder is the mock recorder for MockValidatorStore
type MockValidatorStoreMockRecorder struct {
	mock *MockValidatorStore
}

// NewMockValidatorStore creates a new mock instance
func NewMockValidatorStore(ctrl *gomock.Controller) *MockValidatorStore {
	mock := &MockValidatorStore{ctrl: ctrl}
	mock.recorder = &MockValidatorStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidatorStore) EXPECT() *MockValidatorStoreMockRecorder {
	return m.recorder
}

// SaveValidator mocks base method
func (m *MockValidatorStore) SaveValidator(ctx context.Context, validator structs.Validator) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveValidator", ctx, validator)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveValidator indicates an expected call of SaveValidator
func (mr *MockValidatorStoreMockRecorder) SaveValidator(ctx, validator interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveValidator", reflect.TypeOf((*MockValidatorStore)(nil).SaveValidator), ctx, validator)
}

// GetValidators mocks base method
func (m *MockValidatorStore) GetValidators(ctx context.Context, params structs.QueryParams) ([]structs.Validator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidators", ctx, params)
	ret0, _ := ret[0].([]structs.Validator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidators indicates an expected call of GetValidators
func (mr *MockValidatorStoreMockRecorder) GetValidators(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidators", reflect.TypeOf((*MockValidatorStore)(nil).GetValidators), ctx, params)
}

// MockDelegationStore is a mock of DelegationStore interface
type MockDelegationStore struct {
	ctrl     *gomock.Controller
	recorder *MockDelegationStoreMockRecorder
}

// MockDelegationStoreMockRecorder is the mock recorder for MockDelegationStore
type MockDelegationStoreMockRecorder struct {
	mock *MockDelegationStore
}

// NewMockDelegationStore creates a new mock instance
func NewMockDelegationStore(ctrl *gomock.Controller) *MockDelegationStore {
	mock := &MockDelegationStore{ctrl: ctrl}
	mock.recorder = &MockDelegationStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDelegationStore) EXPECT() *MockDelegationStoreMockRecorder {
	return m.recorder
}

// SaveDelegation mocks base method
func (m *MockDelegationStore) SaveDelegation(ctx context.Context, delegation structs.Delegation) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveDelegation", ctx, delegation)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveDelegation indicates an expected call of SaveDelegation
func (mr *MockDelegationStoreMockRecorder) SaveDelegation(ctx, delegation interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveDelegation", reflect.TypeOf((*MockDelegationStore)(nil).SaveDelegation), ctx, delegation)
}

// GetDelegations mocks base method
func (m *MockDelegationStore) GetDelegations(ctx context.Context, params structs.DelegationParams) ([]structs.Delegation, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDelegations", ctx, params)
	ret0, _ := ret[0].([]structs.Delegation)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDelegations indicates an expected call of GetDelegations
func (mr *MockDelegationStoreMockRecorder) GetDelegations(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDelegations", reflect.TypeOf((*MockDelegationStore)(nil).GetDelegations), ctx, params)
}

// MockValidatorStatisticsStore is a mock of ValidatorStatisticsStore interface
type MockValidatorStatisticsStore struct {
	ctrl     *gomock.Controller
	recorder *MockValidatorStatisticsStoreMockRecorder
}

// MockValidatorStatisticsStoreMockRecorder is the mock recorder for MockValidatorStatisticsStore
type MockValidatorStatisticsStoreMockRecorder struct {
	mock *MockValidatorStatisticsStore
}

// NewMockValidatorStatisticsStore creates a new mock instance
func NewMockValidatorStatisticsStore(ctrl *gomock.Controller) *MockValidatorStatisticsStore {
	mock := &MockValidatorStatisticsStore{ctrl: ctrl}
	mock.recorder = &MockValidatorStatisticsStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockValidatorStatisticsStore) EXPECT() *MockValidatorStatisticsStoreMockRecorder {
	return m.recorder
}

// GetValidatorStatistics mocks base method
func (m *MockValidatorStatisticsStore) GetValidatorStatistics(ctx context.Context, params structs.QueryParams) ([]structs.ValidatorStatistics, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetValidatorStatistics", ctx, params)
	ret0, _ := ret[0].([]structs.ValidatorStatistics)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetValidatorStatistics indicates an expected call of GetValidatorStatistics
func (mr *MockValidatorStatisticsStoreMockRecorder) GetValidatorStatistics(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetValidatorStatistics", reflect.TypeOf((*MockValidatorStatisticsStore)(nil).GetValidatorStatistics), ctx, params)
}
