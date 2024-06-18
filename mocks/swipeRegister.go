// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/AlecSmith96/dating-api/internal/usecases (interfaces: SwipeRegister)
//
// Generated by this command:
//
//	mockgen --build_flags=--mod=mod -destination=../../mocks/swipeRegister.go . SwipeRegister
//
// Package mock_usecases is a generated GoMock package.
package mock_usecases

import (
	reflect "reflect"

	entities "github.com/AlecSmith96/dating-api/internal/entities"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockSwipeRegister is a mock of SwipeRegister interface.
type MockSwipeRegister struct {
	ctrl     *gomock.Controller
	recorder *MockSwipeRegisterMockRecorder
}

// MockSwipeRegisterMockRecorder is the mock recorder for MockSwipeRegister.
type MockSwipeRegisterMockRecorder struct {
	mock *MockSwipeRegister
}

// NewMockSwipeRegister creates a new mock instance.
func NewMockSwipeRegister(ctrl *gomock.Controller) *MockSwipeRegister {
	mock := &MockSwipeRegister{ctrl: ctrl}
	mock.recorder = &MockSwipeRegisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwipeRegister) EXPECT() *MockSwipeRegisterMockRecorder {
	return m.recorder
}

// IsMatch mocks base method.
func (m *MockSwipeRegister) IsMatch(arg0, arg1 uuid.UUID) (*entities.Match, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsMatch", arg0, arg1)
	ret0, _ := ret[0].(*entities.Match)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsMatch indicates an expected call of IsMatch.
func (mr *MockSwipeRegisterMockRecorder) IsMatch(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsMatch", reflect.TypeOf((*MockSwipeRegister)(nil).IsMatch), arg0, arg1)
}

// RegisterSwipe mocks base method.
func (m *MockSwipeRegister) RegisterSwipe(arg0, arg1 uuid.UUID, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterSwipe", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterSwipe indicates an expected call of RegisterSwipe.
func (mr *MockSwipeRegisterMockRecorder) RegisterSwipe(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterSwipe", reflect.TypeOf((*MockSwipeRegister)(nil).RegisterSwipe), arg0, arg1, arg2)
}