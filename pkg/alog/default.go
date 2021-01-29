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

// Default returns the default logger, which is thread-safe.
func Default() Logger {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	return logger
}

// SetDefault sets a new default logger.
func SetDefault(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger = l
}

// WithStr returns a derived Logger from the default Logger with a string field.
func WithStr(key, val string) Logger {
	return Default().WithStr(key, val)
}

// WithFields returns a derived Logger from the default Logger using a map to
// set fields.
func WithFields(fields map[string]interface{}) Logger {
	return Default().WithFields(fields)
}

// Debug logs a new message with debug level.
func Debug(v ...interface{}) {
	Default().Debug(fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func Debugf(format string, v ...interface{}) {
	Default().Debugf(format, v...)
}

// Info logs a new message with info level.
func Info(v ...interface{}) {
	Default().Info(fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func Infof(format string, v ...interface{}) {
	Default().Infof(format, v...)
}

// Error logs a new message with error level.
func Error(v ...interface{}) {
	Default().Error(fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func Errorf(format string, v ...interface{}) {
	Default().Errorf(format, v...)
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
func Fatal(v ...interface{}) {
	Default().Fatal(fmt.Sprint(v...))
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	Default().Fatalf(format, v...)
}
