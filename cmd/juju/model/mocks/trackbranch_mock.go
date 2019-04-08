// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/juju/juju/cmd/juju/model (interfaces: TrackBranchCommandAPI)

// Package mocks is a generated GoMock package.
package mocks

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockTrackBranchCommandAPI is a mock of TrackBranchCommandAPI interface
type MockTrackBranchCommandAPI struct {
	ctrl     *gomock.Controller
	recorder *MockTrackBranchCommandAPIMockRecorder
}

// MockTrackBranchCommandAPIMockRecorder is the mock recorder for MockTrackBranchCommandAPI
type MockTrackBranchCommandAPIMockRecorder struct {
	mock *MockTrackBranchCommandAPI
}

// NewMockTrackBranchCommandAPI creates a new mock instance
func NewMockTrackBranchCommandAPI(ctrl *gomock.Controller) *MockTrackBranchCommandAPI {
	mock := &MockTrackBranchCommandAPI{ctrl: ctrl}
	mock.recorder = &MockTrackBranchCommandAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTrackBranchCommandAPI) EXPECT() *MockTrackBranchCommandAPIMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockTrackBranchCommandAPI) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockTrackBranchCommandAPIMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockTrackBranchCommandAPI)(nil).Close))
}

// TrackBranch mocks base method
func (m *MockTrackBranchCommandAPI) TrackBranch(arg0, arg1 string, arg2 []string) error {
	ret := m.ctrl.Call(m, "TrackBranch", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// TrackBranch indicates an expected call of TrackBranch
func (mr *MockTrackBranchCommandAPIMockRecorder) TrackBranch(arg0, arg1, arg2 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TrackBranch", reflect.TypeOf((*MockTrackBranchCommandAPI)(nil).TrackBranch), arg0, arg1, arg2)
}
