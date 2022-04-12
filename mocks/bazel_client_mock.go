// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ewhauser/bazel-differ/internal (interfaces: BazelClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	internal "github.com/ewhauser/bazel-differ/internal"
	gomock "github.com/golang/mock/gomock"
)

// MockBazelClient is a mock of BazelClient interface.
type MockBazelClient struct {
	ctrl     *gomock.Controller
	recorder *MockBazelClientMockRecorder
}

// MockBazelClientMockRecorder is the mock recorder for MockBazelClient.
type MockBazelClientMockRecorder struct {
	mock *MockBazelClient
}

// NewMockBazelClient creates a new mock instance.
func NewMockBazelClient(ctrl *gomock.Controller) *MockBazelClient {
	mock := &MockBazelClient{ctrl: ctrl}
	mock.recorder = &MockBazelClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBazelClient) EXPECT() *MockBazelClientMockRecorder {
	return m.recorder
}

// QueryAllSourceFileTargets mocks base method.
func (m *MockBazelClient) QueryAllSourceFileTargets() (map[string]*internal.BazelSourceFileTarget, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllSourceFileTargets")
	ret0, _ := ret[0].(map[string]*internal.BazelSourceFileTarget)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllSourceFileTargets indicates an expected call of QueryAllSourceFileTargets.
func (mr *MockBazelClientMockRecorder) QueryAllSourceFileTargets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllSourceFileTargets", reflect.TypeOf((*MockBazelClient)(nil).QueryAllSourceFileTargets))
}

// QueryAllTargets mocks base method.
func (m *MockBazelClient) QueryAllTargets() ([]*internal.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryAllTargets")
	ret0, _ := ret[0].([]*internal.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryAllTargets indicates an expected call of QueryAllTargets.
func (mr *MockBazelClientMockRecorder) QueryAllTargets() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryAllTargets", reflect.TypeOf((*MockBazelClient)(nil).QueryAllTargets))
}

// QueryTarget mocks base method.
func (m *MockBazelClient) QueryTarget(arg0 string, arg1 map[string]bool) ([]*internal.Target, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueryTarget", arg0, arg1)
	ret0, _ := ret[0].([]*internal.Target)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryTarget indicates an expected call of QueryTarget.
func (mr *MockBazelClientMockRecorder) QueryTarget(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryTarget", reflect.TypeOf((*MockBazelClient)(nil).QueryTarget), arg0, arg1)
}
