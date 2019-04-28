// Code generated by MockGen. DO NOT EDIT.
// Source: broker.go

// Package mqb is a generated GoMock package.
package mqb

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockBrokerInterface is a mock of BrokerInterface interface
type MockBrokerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBrokerInterfaceMockRecorder
}

// MockBrokerInterfaceMockRecorder is the mock recorder for MockBrokerInterface
type MockBrokerInterfaceMockRecorder struct {
	mock *MockBrokerInterface
}

// NewMockBrokerInterface creates a new mock instance
func NewMockBrokerInterface(ctrl *gomock.Controller) *MockBrokerInterface {
	mock := &MockBrokerInterface{ctrl: ctrl}
	mock.recorder = &MockBrokerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBrokerInterface) EXPECT() *MockBrokerInterfaceMockRecorder {
	return m.recorder
}

// configure mocks base method
func (m *MockBrokerInterface) configure(arg0 brokerOptions) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "configure", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// configure indicates an expected call of configure
func (mr *MockBrokerInterfaceMockRecorder) configure(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "configure", reflect.TypeOf((*MockBrokerInterface)(nil).configure), arg0)
}

// getNotify mocks base method
func (m *MockBrokerInterface) getNotify() chan error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getNotify")
	ret0, _ := ret[0].(chan error)
	return ret0
}

// getNotify indicates an expected call of getNotify
func (mr *MockBrokerInterfaceMockRecorder) getNotify() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getNotify", reflect.TypeOf((*MockBrokerInterface)(nil).getNotify))
}

// serve mocks base method
func (m *MockBrokerInterface) serve() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "serve")
	ret0, _ := ret[0].(error)
	return ret0
}

// serve indicates an expected call of serve
func (mr *MockBrokerInterfaceMockRecorder) serve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "serve", reflect.TypeOf((*MockBrokerInterface)(nil).serve))
}

// stop mocks base method
func (m *MockBrokerInterface) stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// stop indicates an expected call of stop
func (mr *MockBrokerInterfaceMockRecorder) stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "stop", reflect.TypeOf((*MockBrokerInterface)(nil).stop))
}

// newConsumer mocks base method
func (m *MockBrokerInterface) newConsumer(ctx context.Context, iConsumerOptions interface{}, pool poolInterface) (ConsumerInterface, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "newConsumer", ctx, iConsumerOptions, pool)
	ret0, _ := ret[0].(ConsumerInterface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// newConsumer indicates an expected call of newConsumer
func (mr *MockBrokerInterfaceMockRecorder) newConsumer(ctx, iConsumerOptions, pool interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "newConsumer", reflect.TypeOf((*MockBrokerInterface)(nil).newConsumer), ctx, iConsumerOptions, pool)
}

// newProducer mocks base method
func (m *MockBrokerInterface) newProducer(iProducerOptions interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "newProducer", iProducerOptions)
	ret0, _ := ret[0].(error)
	return ret0
}

// newProducer indicates an expected call of newProducer
func (mr *MockBrokerInterfaceMockRecorder) newProducer(iProducerOptions interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "newProducer", reflect.TypeOf((*MockBrokerInterface)(nil).newProducer), iProducerOptions)
}

// Publish mocks base method
func (m *MockBrokerInterface) Publish(iRequest interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", iRequest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockBrokerInterfaceMockRecorder) Publish(iRequest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockBrokerInterface)(nil).Publish), iRequest)
}

// startConsumer mocks base method
func (m *MockBrokerInterface) startConsumer(name string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "startConsumer", name)
	ret0, _ := ret[0].(error)
	return ret0
}

// startConsumer indicates an expected call of startConsumer
func (mr *MockBrokerInterfaceMockRecorder) startConsumer(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "startConsumer", reflect.TypeOf((*MockBrokerInterface)(nil).startConsumer), name)
}

// MockConsumerInterface is a mock of ConsumerInterface interface
type MockConsumerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerInterfaceMockRecorder
}

// MockConsumerInterfaceMockRecorder is the mock recorder for MockConsumerInterface
type MockConsumerInterfaceMockRecorder struct {
	mock *MockConsumerInterface
}

// NewMockConsumerInterface creates a new mock instance
func NewMockConsumerInterface(ctrl *gomock.Controller) *MockConsumerInterface {
	mock := &MockConsumerInterface{ctrl: ctrl}
	mock.recorder = &MockConsumerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConsumerInterface) EXPECT() *MockConsumerInterfaceMockRecorder {
	return m.recorder
}

// Start mocks base method
func (m *MockConsumerInterface) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start
func (mr *MockConsumerInterfaceMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockConsumerInterface)(nil).Start))
}

// Stop mocks base method
func (m *MockConsumerInterface) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop
func (mr *MockConsumerInterfaceMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockConsumerInterface)(nil).Stop))
}

// Destroy mocks base method
func (m *MockConsumerInterface) Destroy() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

// Destroy indicates an expected call of Destroy
func (mr *MockConsumerInterfaceMockRecorder) Destroy() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Destroy", reflect.TypeOf((*MockConsumerInterface)(nil).Destroy))
}

// getState mocks base method
func (m *MockConsumerInterface) getState() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getState")
	ret0, _ := ret[0].(int)
	return ret0
}

// getState indicates an expected call of getState
func (mr *MockConsumerInterfaceMockRecorder) getState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getState", reflect.TypeOf((*MockConsumerInterface)(nil).getState))
}

// setState mocks base method
func (m *MockConsumerInterface) setState(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "setState", arg0)
}

// setState indicates an expected call of setState
func (mr *MockConsumerInterfaceMockRecorder) setState(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setState", reflect.TypeOf((*MockConsumerInterface)(nil).setState), arg0)
}

// setName mocks base method
func (m *MockConsumerInterface) setName(name string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "setName", name)
}

// setName indicates an expected call of setName
func (mr *MockConsumerInterfaceMockRecorder) setName(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "setName", reflect.TypeOf((*MockConsumerInterface)(nil).setName), name)
}

// getName mocks base method
func (m *MockConsumerInterface) getName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getName")
	ret0, _ := ret[0].(string)
	return ret0
}

// getName indicates an expected call of getName
func (mr *MockConsumerInterfaceMockRecorder) getName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getName", reflect.TypeOf((*MockConsumerInterface)(nil).getName))
}

// getPool mocks base method
func (m *MockConsumerInterface) getPool() poolInterface {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getPool")
	ret0, _ := ret[0].(poolInterface)
	return ret0
}

// getPool indicates an expected call of getPool
func (mr *MockConsumerInterfaceMockRecorder) getPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getPool", reflect.TypeOf((*MockConsumerInterface)(nil).getPool))
}

// poolStart mocks base method
func (m *MockConsumerInterface) poolStart() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "poolStart")
	ret0, _ := ret[0].(error)
	return ret0
}

// poolStart indicates an expected call of poolStart
func (mr *MockConsumerInterfaceMockRecorder) poolStart() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "poolStart", reflect.TypeOf((*MockConsumerInterface)(nil).poolStart))
}

// poolStop mocks base method
func (m *MockConsumerInterface) poolStop() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "poolStop")
}

// poolStop indicates an expected call of poolStop
func (mr *MockConsumerInterfaceMockRecorder) poolStop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "poolStop", reflect.TypeOf((*MockConsumerInterface)(nil).poolStop))
}
