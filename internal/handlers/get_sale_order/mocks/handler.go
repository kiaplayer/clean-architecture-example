// Code generated by MockGen. DO NOT EDIT.
// Source: handler.go
//
// Generated by this command:
//
//	mockgen -package=get_sale_order -source=handler.go -destination=mocks/handler.go
//

// Package get_sale_order is a generated GoMock package.
package get_sale_order

import (
	context "context"
	reflect "reflect"

	document "github.com/kiaplayer/clean-architecture-example/internal/domain/entity/document"
	gomock "go.uber.org/mock/gomock"
)

// MockuseCase is a mock of useCase interface.
type MockuseCase struct {
	ctrl     *gomock.Controller
	recorder *MockuseCaseMockRecorder
}

// MockuseCaseMockRecorder is the mock recorder for MockuseCase.
type MockuseCaseMockRecorder struct {
	mock *MockuseCase
}

// NewMockuseCase creates a new mock instance.
func NewMockuseCase(ctrl *gomock.Controller) *MockuseCase {
	mock := &MockuseCase{ctrl: ctrl}
	mock.recorder = &MockuseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuseCase) EXPECT() *MockuseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockuseCase) Handle(ctx context.Context, id uint64) (*document.SaleOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, id)
	ret0, _ := ret[0].(*document.SaleOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockuseCaseMockRecorder) Handle(ctx, id any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockuseCase)(nil).Handle), ctx, id)
}
