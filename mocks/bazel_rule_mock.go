// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ewhauser/bazel-differ/internal (interfaces: BazelRule)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBazelRule is a mock of BazelRule interface.
type MockBazelRule struct {
	ctrl     *gomock.Controller
	recorder *MockBazelRuleMockRecorder
}

// MockBazelRuleMockRecorder is the mock recorder for MockBazelRule.
type MockBazelRuleMockRecorder struct {
	mock *MockBazelRule
}

// NewMockBazelRule creates a new mock instance.
func NewMockBazelRule(ctrl *gomock.Controller) *MockBazelRule {
	mock := &MockBazelRule{ctrl: ctrl}
	mock.recorder = &MockBazelRuleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBazelRule) EXPECT() *MockBazelRuleMockRecorder {
	return m.recorder
}

// Digest mocks base method.
func (m *MockBazelRule) Digest() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Digest")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Digest indicates an expected call of Digest.
func (mr *MockBazelRuleMockRecorder) Digest() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Digest", reflect.TypeOf((*MockBazelRule)(nil).Digest))
}

// Name mocks base method.
func (m *MockBazelRule) Name() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Name")
	ret0, _ := ret[0].(string)
	return ret0
}

// Name indicates an expected call of Name.
func (mr *MockBazelRuleMockRecorder) Name() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Name", reflect.TypeOf((*MockBazelRule)(nil).Name))
}

// RuleInputList mocks base method.
func (m *MockBazelRule) RuleInputList() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RuleInputList")
	ret0, _ := ret[0].([]string)
	return ret0
}

// RuleInputList indicates an expected call of RuleInputList.
func (mr *MockBazelRuleMockRecorder) RuleInputList() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RuleInputList", reflect.TypeOf((*MockBazelRule)(nil).RuleInputList))
}
