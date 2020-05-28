// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/apiserver/facades/client/machinemanager (interfaces: Leadership)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	params "github.com/juju/juju/apiserver/params"
	names "gopkg.in/juju/names.v3"
	reflect "reflect"
)

// MockLeadership is a mock of Leadership interface
type MockLeadership struct {
	ctrl     *gomock.Controller
	recorder *MockLeadershipMockRecorder
}

// MockLeadershipMockRecorder is the mock recorder for MockLeadership
type MockLeadershipMockRecorder struct {
	mock *MockLeadership
}

// NewMockLeadership creates a new mock instance
func NewMockLeadership(ctrl *gomock.Controller) *MockLeadership {
	mock := &MockLeadership{ctrl: ctrl}
	mock.recorder = &MockLeadershipMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockLeadership) EXPECT() *MockLeadershipMockRecorder {
	return m.recorder
}

// GetMachineApplicationNames mocks base method
func (m *MockLeadership) GetMachineApplicationNames(arg0 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMachineApplicationNames", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMachineApplicationNames indicates an expected call of GetMachineApplicationNames
func (mr *MockLeadershipMockRecorder) GetMachineApplicationNames(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMachineApplicationNames", reflect.TypeOf((*MockLeadership)(nil).GetMachineApplicationNames), arg0)
}

// UnpinApplicationLeadersByName mocks base method
func (m *MockLeadership) UnpinApplicationLeadersByName(arg0 names.Tag, arg1 []string) (params.PinApplicationsResults, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnpinApplicationLeadersByName", arg0, arg1)
	ret0, _ := ret[0].(params.PinApplicationsResults)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UnpinApplicationLeadersByName indicates an expected call of UnpinApplicationLeadersByName
func (mr *MockLeadershipMockRecorder) UnpinApplicationLeadersByName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnpinApplicationLeadersByName", reflect.TypeOf((*MockLeadership)(nil).UnpinApplicationLeadersByName), arg0, arg1)
}
