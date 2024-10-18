// Code generated by MockGen. DO NOT EDIT.
// Source: validator.go
//
// Generated by this command:
//
//	mockgen -source validator.go -destination mock_devicer_test.go -package validator
//

// Package validator is a generated GoMock package.
package validator

import (
	context "context"
	reflect "reflect"

	api "github.com/thingspect/proto/go/api"
	gomock "go.uber.org/mock/gomock"
)

// Mockdevicer is a mock of devicer interface.
type Mockdevicer struct {
	ctrl     *gomock.Controller
	recorder *MockdevicerMockRecorder
	isgomock struct{}
}

// MockdevicerMockRecorder is the mock recorder for Mockdevicer.
type MockdevicerMockRecorder struct {
	mock *Mockdevicer
}

// NewMockdevicer creates a new mock instance.
func NewMockdevicer(ctrl *gomock.Controller) *Mockdevicer {
	mock := &Mockdevicer{ctrl: ctrl}
	mock.recorder = &MockdevicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockdevicer) EXPECT() *MockdevicerMockRecorder {
	return m.recorder
}

// ReadByUniqID mocks base method.
func (m *Mockdevicer) ReadByUniqID(ctx context.Context, uniqID string) (*api.Device, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadByUniqID", ctx, uniqID)
	ret0, _ := ret[0].(*api.Device)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReadByUniqID indicates an expected call of ReadByUniqID.
func (mr *MockdevicerMockRecorder) ReadByUniqID(ctx, uniqID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadByUniqID", reflect.TypeOf((*Mockdevicer)(nil).ReadByUniqID), ctx, uniqID)
}
