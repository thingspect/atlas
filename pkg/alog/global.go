package alog

import (
	"fmt"
	"sync"
)

// Since logger is global and may be replaced, locking is required.
var (
	logger   Logger
	loggerMu sync.Mutex
)

// NewConsole returns a new Logger with console formatting at the debug level.
func NewConsole() Logger {
	return newZlogConsole()
}

// NewJSON returns a new Logger with JSON formatting at the debug level.
func NewJSON() Logger {
	return newZlogJSON()
}

// Global returns the global logger, which is thread-safe.
func Global() Logger {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	return logger
}

// SetGlobal sets a new global logger.
func SetGlobal(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger = l
}

// WithStr returns a derived Logger from the global Logger with a string field.
func WithStr(key, val string) Logger {
	return Global().WithStr(key, val)
}

// WithFields returns a derived Logger from the global Logger using a map to set
// fields.
func WithFields(fields map[string]interface{}) Logger {
	return Global().WithFields(fields)
}

// Debug logs a new message with debug level.
func Debug(v ...interface{}) {
	Global().Debug(fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func Debugf(format string, v ...interface{}) {
	Global().Debugf(format, v...)
}

// Info logs a new message with info level.
func Info(v ...interface{}) {
	Global().Info(fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func Infof(format string, v ...interface{}) {
	Global().Infof(format, v...)
}

// Error logs a new message with error level.
func Error(v ...interface{}) {
	Global().Error(fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func Errorf(format string, v ...interface{}) {
	Global().Errorf(format, v...)
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Global().Fatal(fmt.Sprint(v...))
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Global().Fatalf(format, v...)
}
