// Code generated by MockGen. DO NOT EDIT.
// Source: accumulator.go
//
// Generated by this command:
//
//	mockgen -source accumulator.go -destination mock_datapointer_test.go -package accumulator
//

// Package accumulator is a generated GoMock package.
package accumulator

import (
	context "context"
	reflect "reflect"

	common "github.com/thingspect/proto/go/common"
	gomock "go.uber.org/mock/gomock"
)

// Mockdatapointer is a mock of datapointer interface.
type Mockdatapointer struct {
	ctrl     *gomock.Controller
	recorder *MockdatapointerMockRecorder
	isgomock struct{}
}

// MockdatapointerMockRecorder is the mock recorder for Mockdatapointer.
type MockdatapointerMockRecorder struct {
	mock *Mockdatapointer
}

// NewMockdatapointer creates a new mock instance.
func NewMockdatapointer(ctrl *gomock.Controller) *Mockdatapointer {
	mock := &Mockdatapointer{ctrl: ctrl}
	mock.recorder = &MockdatapointerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdatapointer) EXPECT() *MockdatapointerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockdatapointer) Create(ctx context.Context, point *common.DataPoint, orgID string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, point, orgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockdatapointerMockRecorder) Create(ctx, point, orgID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockdatapointer)(nil).Create), ctx, point, orgID)
}
