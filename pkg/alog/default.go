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

// NewConsole returns a new Logger with console formatting at the specified
// level.
func NewConsole(level string) Logger {
	return newStlogConsole(level)
}

// NewJSON returns a new Logger with JSON formatting at the specified level.
func NewJSON(level string) Logger {
	return newStlogJSON(level)
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
	logger = l
	loggerMu.Unlock()
}

// WithField returns a derived Logger from the default Logger with a string
// field.
func WithField(key, val string) Logger {
	return Default().WithField(key, val)
}

// Debug logs a new message with debug level.
func Debug(v ...any) {
	Default().Debug(fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func Debugf(format string, v ...any) {
	Default().Debugf(format, v...)
}

// Info logs a new message with info level.
func Info(v ...any) {
	Default().Info(fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func Infof(format string, v ...any) {
	Default().Infof(format, v...)
}

// Error logs a new message with error level.
func Error(v ...any) {
	Default().Error(fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func Errorf(format string, v ...any) {
	Default().Errorf(format, v...)
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
func Fatal(v ...any) {
	Default().Fatal(fmt.Sprint(v...))
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func Fatalf(format string, v ...any) {
	Default().Fatalf(format, v...)
}
