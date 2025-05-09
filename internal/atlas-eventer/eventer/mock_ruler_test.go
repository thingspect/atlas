// Code generated by MockGen. DO NOT EDIT.
// Source: eventer.go
//
// Generated by this command:
//
//	mockgen -source eventer.go -destination mock_ruler_test.go -package eventer
//

// Package eventer is a generated GoMock package.
package eventer

import (
	context "context"
	reflect "reflect"

	api "github.com/thingspect/proto/go/api"
	gomock "go.uber.org/mock/gomock"
)

// Mockruler is a mock of ruler interface.
type Mockruler struct {
	ctrl     *gomock.Controller
	recorder *MockrulerMockRecorder
	isgomock struct{}
}

// MockrulerMockRecorder is the mock recorder for Mockruler.
type MockrulerMockRecorder struct {
	mock *Mockruler
}

// NewMockruler creates a new mock instance.
func NewMockruler(ctrl *gomock.Controller) *Mockruler {
	mock := &Mockruler{ctrl: ctrl}
	mock.recorder = &MockrulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockruler) EXPECT() *MockrulerMockRecorder {
	return m.recorder
}

// ListByTags mocks base method.
func (m *Mockruler) ListByTags(ctx context.Context, orgID, attr string, deviceTags []string) ([]*api.Rule, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTags", ctx, orgID, attr, deviceTags)
	ret0, _ := ret[0].([]*api.Rule)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByTags indicates an expected call of ListByTags.
func (mr *MockrulerMockRecorder) ListByTags(ctx, orgID, attr, deviceTags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTags", reflect.TypeOf((*Mockruler)(nil).ListByTags), ctx, orgID, attr, deviceTags)
}

// Mockeventer is a mock of eventer interface.
type Mockeventer struct {
	ctrl     *gomock.Controller
	recorder *MockeventerMockRecorder
	isgomock struct{}
}

// MockeventerMockRecorder is the mock recorder for Mockeventer.
type MockeventerMockRecorder struct {
	mock *Mockeventer
}

// NewMockeventer creates a new mock instance.
func NewMockeventer(ctrl *gomock.Controller) *Mockeventer {
	mock := &Mockeventer{ctrl: ctrl}
	mock.recorder = &MockeventerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockeventer) EXPECT() *MockeventerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockeventer) Create(ctx context.Context, event *api.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockeventerMockRecorder) Create(ctx, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockeventer)(nil).Create), ctx, event)
}
