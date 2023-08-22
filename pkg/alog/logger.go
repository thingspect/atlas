// Package alog provides functions to write logs with a friendly API. An
// interface and implementation was chosen over simpler package wrappers to
// support changing implementations.
package alog

// Logger defines the methods provided by a Log.
type Logger interface {
	// WithField returns a derived Logger with a string field.
	WithField(key, val string) Logger

	// Debug logs a new message with debug level.
	Debug(v ...any)
	// Debugf logs a new formatted message with debug level.
	Debugf(format string, v ...any)
	// Info logs a new message with info level.
	Info(v ...any)
	// Infof logs a new formatted message with info level.
	Infof(format string, v ...any)
	// Error logs a new message with error level.
	Error(v ...any)
	// Errorf logs a new formatted message with error level.
	Errorf(format string, v ...any)
	// Fatal logs a new message with fatal level followed by a call to
	// os.Exit(1).
	Fatal(v ...any)
	// Fatalf logs a new formatted message with fatal level followed by a call
	// to os.Exit(1).
	Fatalf(format string, v ...any)
}
