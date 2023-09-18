// Code generated by MockGen. DO NOT EDIT.
// Source: alerter.go
//
// Generated by this command:
//
//	mockgen -source alerter.go -destination mock_alarmer_test.go -package alerter
//
// Package alerter is a generated GoMock package.
package alerter

import (
	context "context"
	reflect "reflect"
	time "time"

	api "github.com/thingspect/api/go/api"
	gomock "go.uber.org/mock/gomock"
)

// Mockorger is a mock of orger interface.
type Mockorger struct {
	ctrl     *gomock.Controller
	recorder *MockorgerMockRecorder
}

// MockorgerMockRecorder is the mock recorder for Mockorger.
type MockorgerMockRecorder struct {
	mock *Mockorger
}

// NewMockorger creates a new mock instance.
func NewMockorger(ctrl *gomock.Controller) *Mockorger {
	mock := &Mockorger{ctrl: ctrl}
	mock.recorder = &MockorgerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockorger) EXPECT() *MockorgerMockRecorder {
	return m.recorder
}

// Read mocks base method.
func (m *Mockorger) Read(ctx context.Context, orgID string) (*api.Org, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Read", ctx, orgID)
	ret0, _ := ret[0].(*api.Org)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Read indicates an expected call of Read.
func (mr *MockorgerMockRecorder) Read(ctx, orgID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Read", reflect.TypeOf((*Mockorger)(nil).Read), ctx, orgID)
}

// Mockalarmer is a mock of alarmer interface.
type Mockalarmer struct {
	ctrl     *gomock.Controller
	recorder *MockalarmerMockRecorder
}

// MockalarmerMockRecorder is the mock recorder for Mockalarmer.
type MockalarmerMockRecorder struct {
	mock *Mockalarmer
}

// NewMockalarmer creates a new mock instance.
func NewMockalarmer(ctrl *gomock.Controller) *Mockalarmer {
	mock := &Mockalarmer{ctrl: ctrl}
	mock.recorder = &MockalarmerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockalarmer) EXPECT() *MockalarmerMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *Mockalarmer) List(ctx context.Context, orgID string, lBoundTS time.Time, prevID string, limit int32, ruleID string) ([]*api.Alarm, int32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", ctx, orgID, lBoundTS, prevID, limit, ruleID)
	ret0, _ := ret[0].([]*api.Alarm)
	ret1, _ := ret[1].(int32)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockalarmerMockRecorder) List(ctx, orgID, lBoundTS, prevID, limit, ruleID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*Mockalarmer)(nil).List), ctx, orgID, lBoundTS, prevID, limit, ruleID)
}

// Mockuserer is a mock of userer interface.
type Mockuserer struct {
	ctrl     *gomock.Controller
	recorder *MockusererMockRecorder
}

// MockusererMockRecorder is the mock recorder for Mockuserer.
type MockusererMockRecorder struct {
	mock *Mockuserer
}

// NewMockuserer creates a new mock instance.
func NewMockuserer(ctrl *gomock.Controller) *Mockuserer {
	mock := &Mockuserer{ctrl: ctrl}
	mock.recorder = &MockusererMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockuserer) EXPECT() *MockusererMockRecorder {
	return m.recorder
}

// ListByTags mocks base method.
func (m *Mockuserer) ListByTags(ctx context.Context, orgID string, tags []string) ([]*api.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByTags", ctx, orgID, tags)
	ret0, _ := ret[0].([]*api.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByTags indicates an expected call of ListByTags.
func (mr *MockusererMockRecorder) ListByTags(ctx, orgID, tags any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByTags", reflect.TypeOf((*Mockuserer)(nil).ListByTags), ctx, orgID, tags)
}

// Mockalerter is a mock of alerter interface.
type Mockalerter struct {
	ctrl     *gomock.Controller
	recorder *MockalerterMockRecorder
}

// MockalerterMockRecorder is the mock recorder for Mockalerter.
type MockalerterMockRecorder struct {
	mock *Mockalerter
}

// NewMockalerter creates a new mock instance.
func NewMockalerter(ctrl *gomock.Controller) *Mockalerter {
	mock := &Mockalerter{ctrl: ctrl}
	mock.recorder = &MockalerterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockalerter) EXPECT() *MockalerterMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *Mockalerter) Create(ctx context.Context, alert *api.Alert) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, alert)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockalerterMockRecorder) Create(ctx, alert any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*Mockalerter)(nil).Create), ctx, alert)
}
//lint:file-ignore ST1000 Mockgen package comment
