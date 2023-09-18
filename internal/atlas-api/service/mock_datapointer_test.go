// Code generated by MockGen. DO NOT EDIT.
// Source: datapoint.go
//
// Generated by this command:
//
//	mockgen -source datapoint.go -destination mock_datapointer_test.go -package service
//
// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"
	time "time"

	common "github.com/thingspect/api/go/common"
	gomock "go.uber.org/mock/gomock"
)

// MockDataPointer is a mock of DataPointer interface.
type MockDataPointer struct {
	ctrl     *gomock.Controller
	recorder *MockDataPointerMockRecorder
}

// MockDataPointerMockRecorder is the mock recorder for MockDataPointer.
type MockDataPointerMockRecorder struct {
	mock *MockDataPointer
}

// NewMockDataPointer creates a new mock instance.
func NewMockDataPointer(ctrl *gomock.Controller) *MockDataPointer {
	mock := &MockDataPointer{ctrl: ctrl}
	mock.recorder = &MockDataPointerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataPointer) EXPECT() *MockDataPointerMockRecorder {
	return m.recorder
}

// Latest mocks base method.
func (m *MockDataPointer) Latest(ctx context.Context, orgID, uniqID, devID string, start time.Time) ([]*common.DataPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Latest", ctx, orgID, uniqID, devID, start)
	ret0, _ := ret[0].([]*common.DataPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Latest indicates an expected call of Latest.
func (mr *MockDataPointerMockRecorder) Latest(ctx, orgID, uniqID, devID, start any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Latest", reflect.TypeOf((*MockDataPointer)(nil).Latest), ctx, orgID, uniqID, devID, start)
}

// List mocks base method.
func (m *MockDataPointer) List(ctx context.Context, orgID, uniqID, devID, attr string, end, start time.Time) ([]*common.DataPoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, orgID, uniqID, devID, attr, end, start)
	ret0, _ := ret[0].([]*common.DataPoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockDataPointerMockRecorder) List(ctx, orgID, uniqID, devID, attr, end, start any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDataPointer)(nil).List), ctx, orgID, uniqID, devID, attr, end, start)
}
//lint:file-ignore ST1000 Mockgen package comment
