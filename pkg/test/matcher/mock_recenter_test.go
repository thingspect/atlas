// Code generated by MockGen. DO NOT EDIT.
// Source: recent_test.go

// Package matcher is a generated GoMock package.
package matcher

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// Mockrecenter is a mock of recenter interface.
type Mockrecenter struct {
	ctrl     *gomock.Controller
	recorder *MockrecenterMockRecorder
}

// MockrecenterMockRecorder is the mock recorder for Mockrecenter.
type MockrecenterMockRecorder struct {
	mock *Mockrecenter
}

// NewMockrecenter creates a new mock instance.
func NewMockrecenter(ctrl *gomock.Controller) *Mockrecenter {
	mock := &Mockrecenter{ctrl: ctrl}
	mock.recorder = &MockrecenterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrecenter) EXPECT() *MockrecenterMockRecorder {
	return m.recorder
}

// f mocks base method.
func (m *Mockrecenter) f(t time.Time) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "f", t)
	ret0, _ := ret[0].(error)
	return ret0
}

// f indicates an expected call of f.
func (mr *MockrecenterMockRecorder) f(t interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "f", reflect.TypeOf((*Mockrecenter)(nil).f), t)
}
