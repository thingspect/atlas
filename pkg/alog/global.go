package alog

import (
	"fmt"
	"sync"
)

// Since logger is global and may be replaced, locking is required for all
// operations.
var (
	logger   Logger
	loggerMu sync.Mutex
)

// NewConsole creates a new Logger with console formatting at the debug level.
func NewConsole() Logger {
	return newZlogConsole()
}

// NewJSON creates a new Logger with JSON formatting at the debug level.
func NewJSON() Logger {
	return newZlogJSON()
}

// SetGlobal sets a new global logger.
func SetGlobal(l Logger) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger = l
}

// WithStr creates a derived Logger from the global Logger with a string field.
func WithStr(key, val string) Logger {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	return logger.WithStr(key, val)
}

// WithFields creates a derived Logger from the global Logger using a map to set
// fields.
func WithFields(fields map[string]interface{}) Logger {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	return logger.WithFields(fields)
}

// Debug logs a new message with debug level.
func Debug(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Debug(fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func Debugf(format string, v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Debugf(format, v...)
}

// Info logs a new message with info level.
func Info(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Info(fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func Infof(format string, v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Infof(format, v...)
}

// Error logs a new message with error level.
func Error(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Error(fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func Errorf(format string, v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Errorf(format, v...)
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
//lint:ignore U1001 call to os.Exit(1).
func Fatal(v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Fatal(fmt.Sprint(v...))
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func Fatalf(format string, v ...interface{}) {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	logger.Fatalf(format, v...)
}
