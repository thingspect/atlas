// Code generated by MockGen. DO NOT EDIT.
// Source: notifier.go

// Package notify is a generated GoMock package.
package notify

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockNotifier is a mock of Notifier interface.
type MockNotifier struct {
	ctrl     *gomock.Controller
	recorder *MockNotifierMockRecorder
}

// MockNotifierMockRecorder is the mock recorder for MockNotifier.
type MockNotifierMockRecorder struct {
	mock *MockNotifier
}

// NewMockNotifier creates a new mock instance.
func NewMockNotifier(ctrl *gomock.Controller) *MockNotifier {
	mock := &MockNotifier{ctrl: ctrl}
	mock.recorder = &MockNotifierMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotifier) EXPECT() *MockNotifierMockRecorder {
	return m.recorder
}

// App mocks base method.
func (m *MockNotifier) App(ctx context.Context, userKey, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "App", ctx, userKey, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// App indicates an expected call of App.
func (mr *MockNotifierMockRecorder) App(ctx, userKey, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "App", reflect.TypeOf((*MockNotifier)(nil).App), ctx, userKey, subject, body)
}

// Email mocks base method.
func (m *MockNotifier) Email(ctx context.Context, displayName, orgEmail, userEmail, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Email", ctx, displayName, orgEmail, userEmail, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// Email indicates an expected call of Email.
func (mr *MockNotifierMockRecorder) Email(ctx, displayName, orgEmail, userEmail, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Email", reflect.TypeOf((*MockNotifier)(nil).Email), ctx, displayName, orgEmail, userEmail, subject, body)
}

// SMS mocks base method.
func (m *MockNotifier) SMS(ctx context.Context, phone, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SMS", ctx, phone, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SMS indicates an expected call of SMS.
func (mr *MockNotifierMockRecorder) SMS(ctx, phone, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SMS", reflect.TypeOf((*MockNotifier)(nil).SMS), ctx, phone, subject, body)
}
