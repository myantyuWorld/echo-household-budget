// Code generated by MockGen. DO NOT EDIT.
// Source: shopping.go
//
// Generated by this command:
//
//	mockgen -source=shopping.go -destination=../mock/domainmodel/mock_shopping.go -package=mock
//

// Package mock is a generated GoMock package.
package mock

import (
	domainmodel "echo-household-budget/internal/domain/model"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockShoppingRepository is a mock of ShoppingRepository interface.
type MockShoppingRepository struct {
	ctrl     *gomock.Controller
	recorder *MockShoppingRepositoryMockRecorder
	isgomock struct{}
}

// MockShoppingRepositoryMockRecorder is the mock recorder for MockShoppingRepository.
type MockShoppingRepositoryMockRecorder struct {
	mock *MockShoppingRepository
}

// NewMockShoppingRepository creates a new mock instance.
func NewMockShoppingRepository(ctrl *gomock.Controller) *MockShoppingRepository {
	mock := &MockShoppingRepository{ctrl: ctrl}
	mock.recorder = &MockShoppingRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShoppingRepository) EXPECT() *MockShoppingRepositoryMockRecorder {
	return m.recorder
}

// DeleteShoppingAmount mocks base method.
func (m *MockShoppingRepository) DeleteShoppingAmount(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShoppingAmount", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShoppingAmount indicates an expected call of DeleteShoppingAmount.
func (mr *MockShoppingRepositoryMockRecorder) DeleteShoppingAmount(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShoppingAmount", reflect.TypeOf((*MockShoppingRepository)(nil).DeleteShoppingAmount), id)
}

// DeleteShoppingMemo mocks base method.
func (m *MockShoppingRepository) DeleteShoppingMemo(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteShoppingMemo", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteShoppingMemo indicates an expected call of DeleteShoppingMemo.
func (mr *MockShoppingRepositoryMockRecorder) DeleteShoppingMemo(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteShoppingMemo", reflect.TypeOf((*MockShoppingRepository)(nil).DeleteShoppingMemo), id)
}

// FetchShoppingAmountItem mocks base method.
func (m *MockShoppingRepository) FetchShoppingAmountItem(id string) (*domainmodel.ShoppingAmount, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchShoppingAmountItem", id)
	ret0, _ := ret[0].(*domainmodel.ShoppingAmount)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchShoppingAmountItem indicates an expected call of FetchShoppingAmountItem.
func (mr *MockShoppingRepositoryMockRecorder) FetchShoppingAmountItem(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchShoppingAmountItem", reflect.TypeOf((*MockShoppingRepository)(nil).FetchShoppingAmountItem), id)
}

// FetchShoppingMemoItem mocks base method.
func (m *MockShoppingRepository) FetchShoppingMemoItem(id string) (*domainmodel.ShoppingMemo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchShoppingMemoItem", id)
	ret0, _ := ret[0].(*domainmodel.ShoppingMemo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchShoppingMemoItem indicates an expected call of FetchShoppingMemoItem.
func (mr *MockShoppingRepositoryMockRecorder) FetchShoppingMemoItem(id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchShoppingMemoItem", reflect.TypeOf((*MockShoppingRepository)(nil).FetchShoppingMemoItem), id)
}

// RegisterShoppingAmount mocks base method.
func (m *MockShoppingRepository) RegisterShoppingAmount(shopping *domainmodel.ShoppingAmount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterShoppingAmount", shopping)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterShoppingAmount indicates an expected call of RegisterShoppingAmount.
func (mr *MockShoppingRepositoryMockRecorder) RegisterShoppingAmount(shopping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterShoppingAmount", reflect.TypeOf((*MockShoppingRepository)(nil).RegisterShoppingAmount), shopping)
}

// RegisterShoppingMemo mocks base method.
func (m *MockShoppingRepository) RegisterShoppingMemo(shopping *domainmodel.ShoppingMemo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterShoppingMemo", shopping)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterShoppingMemo indicates an expected call of RegisterShoppingMemo.
func (mr *MockShoppingRepositoryMockRecorder) RegisterShoppingMemo(shopping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterShoppingMemo", reflect.TypeOf((*MockShoppingRepository)(nil).RegisterShoppingMemo), shopping)
}

// UpdateShoppingAmount mocks base method.
func (m *MockShoppingRepository) UpdateShoppingAmount(shopping *domainmodel.ShoppingAmount) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShoppingAmount", shopping)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShoppingAmount indicates an expected call of UpdateShoppingAmount.
func (mr *MockShoppingRepositoryMockRecorder) UpdateShoppingAmount(shopping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShoppingAmount", reflect.TypeOf((*MockShoppingRepository)(nil).UpdateShoppingAmount), shopping)
}

// UpdateShoppingMemo mocks base method.
func (m *MockShoppingRepository) UpdateShoppingMemo(shopping *domainmodel.ShoppingMemo) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateShoppingMemo", shopping)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateShoppingMemo indicates an expected call of UpdateShoppingMemo.
func (mr *MockShoppingRepositoryMockRecorder) UpdateShoppingMemo(shopping any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateShoppingMemo", reflect.TypeOf((*MockShoppingRepository)(nil).UpdateShoppingMemo), shopping)
}
