// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlecSmith96/dating-api/internal/usecases (interfaces: UserDiscoverer)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../../mocks/userDiscoverer.go . UserDiscoverer
//
// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	entities "github.com/AlecSmith96/dating-api/internal/entities"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockUserDiscoverer is a mock of UserDiscoverer interface.
type MockUserDiscoverer struct {
	ctrl     *gomock.Controller
	recorder *MockUserDiscovererMockRecorder
}

// MockUserDiscovererMockRecorder is the mock recorder for MockUserDiscoverer.
type MockUserDiscovererMockRecorder struct {
	mock *MockUserDiscoverer
}

// NewMockUserDiscoverer creates a new mock instance.
func NewMockUserDiscoverer(ctrl *gomock.Controller) *MockUserDiscoverer {
	mock := &MockUserDiscoverer{ctrl: ctrl}
	mock.recorder = &MockUserDiscovererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserDiscoverer) EXPECT() *MockUserDiscovererMockRecorder {
	return m.recorder
}

// DiscoverNewUsers mocks base method.
func (m *MockUserDiscoverer) DiscoverNewUsers(arg0 uuid.UUID, arg1 entities.PageInfo) ([]entities.UserDiscovery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverNewUsers", arg0, arg1)
	ret0, _ := ret[0].([]entities.UserDiscovery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverNewUsers indicates an expected call of DiscoverNewUsers.
func (mr *MockUserDiscovererMockRecorder) DiscoverNewUsers(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverNewUsers", reflect.TypeOf((*MockUserDiscoverer)(nil).DiscoverNewUsers), arg0, arg1)
}

// GetUsersLocation mocks base method.
func (m *MockUserDiscoverer) GetUsersLocation(arg0 uuid.UUID) (*entities.Location, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUsersLocation", arg0)
	ret0, _ := ret[0].(*entities.Location)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersLocation indicates an expected call of GetUsersLocation.
func (mr *MockUserDiscovererMockRecorder) GetUsersLocation(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersLocation", reflect.TypeOf((*MockUserDiscoverer)(nil).GetUsersLocation), arg0)
}
