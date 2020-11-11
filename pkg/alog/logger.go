// Package alog contains functions to write logs with a friendly API. An
// interface and implementation was chosen over simpler package wrappers to
// support Entry-style returns.
package alog

// Logger defines the methods provided by a Log.
type Logger interface {
	WithLevel(level string) Logger

	WithStr(key, val string) Logger
	WithFields(fields map[string]interface{}) Logger

	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	Error(v ...interface{})
	Errorf(format string, v ...interface{})
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
}
