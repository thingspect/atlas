// Code generated by MockGen. DO NOT EDIT.
// Source: alert.go
//
// Generated by this command:
//
//	mockgen -source alert.go -destination mock_alerter_test.go -package service
//
// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"
	time "time"

	api "github.com/thingspect/api/go/api"
	gomock "go.uber.org/mock/gomock"
)

// MockAlerter is a mock of Alerter interface.
type MockAlerter struct {
	ctrl     *gomock.Controller
	recorder *MockAlerterMockRecorder
}

// MockAlerterMockRecorder is the mock recorder for MockAlerter.
type MockAlerterMockRecorder struct {
	mock *MockAlerter
}

// NewMockAlerter creates a new mock instance.
func NewMockAlerter(ctrl *gomock.Controller) *MockAlerter {
	mock := &MockAlerter{ctrl: ctrl}
	mock.recorder = &MockAlerterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAlerter) EXPECT() *MockAlerterMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockAlerter) List(ctx context.Context, orgID, uniqID, devID, alarmID, userID string, end, start time.Time) ([]*api.Alert, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, orgID, uniqID, devID, alarmID, userID, end, start)
	ret0, _ := ret[0].([]*api.Alert)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockAlerterMockRecorder) List(ctx, orgID, uniqID, devID, alarmID, userID, end, start any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAlerter)(nil).List), ctx, orgID, uniqID, devID, alarmID, userID, end, start)
}
//lint:file-ignore ST1000 Mockgen package comment
