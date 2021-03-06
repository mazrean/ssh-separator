// Code generated by MockGen. DO NOT EDIT.
// Source: user.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/mazrean/separated-webshell/domain"
	values "github.com/mazrean/separated-webshell/domain/values"
)

// MockIUser is a mock of IUser interface.
type MockIUser struct {
	ctrl     *gomock.Controller
	recorder *MockIUserMockRecorder
}

// MockIUserMockRecorder is the mock recorder for MockIUser.
type MockIUserMockRecorder struct {
	mock *MockIUser
}

// NewMockIUser creates a new mock instance.
func NewMockIUser(ctrl *gomock.Controller) *MockIUser {
	mock := &MockIUser{ctrl: ctrl}
	mock.recorder = &MockIUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIUser) EXPECT() *MockIUserMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockIUser) Create(ctx context.Context, user *domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockIUserMockRecorder) Create(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockIUser)(nil).Create), ctx, user)
}

// GetAllUser mocks base method.
func (m *MockIUser) GetAllUser(ctx context.Context) ([]values.UserName, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllUser", ctx)
	ret0, _ := ret[0].([]values.UserName)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUser indicates an expected call of GetAllUser.
func (mr *MockIUserMockRecorder) GetAllUser(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUser", reflect.TypeOf((*MockIUser)(nil).GetAllUser), ctx)
}

// GetPassword mocks base method.
func (m *MockIUser) GetPassword(ctx context.Context, userName values.UserName) (values.HashedPassword, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPassword", ctx, userName)
	ret0, _ := ret[0].(values.HashedPassword)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPassword indicates an expected call of GetPassword.
func (mr *MockIUserMockRecorder) GetPassword(ctx, userName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPassword", reflect.TypeOf((*MockIUser)(nil).GetPassword), ctx, userName)
}
