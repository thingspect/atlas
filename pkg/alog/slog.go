package alog

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// levels maps strings to log levels.
var levels = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"ERROR": slog.LevelError,
	"FATAL": slog.LevelError + 4,
}

// stlog contains methods to write logs using slog and implements the Logger
// interface.
type stlog struct {
	sl *slog.Logger
}

// Verify stlog implements Logger.
var _ Logger = &stlog{}

// parseLevel parses a string into a log level.
func parseLevel(level string) slog.Level {
	slevel, ok := levels[strings.ToUpper(level)]
	if !ok {
		slog.LogAttrs(context.Background(), slog.LevelError,
			fmt.Sprintf("parseLevel unknown level, using INFO: %s", level))
		slevel = slog.LevelInfo
	}

	return slevel
}

// newStlogConsole returns a new Logger with console formatting at the specified
// level.
func newStlogConsole(level string) Logger {
	l := new(slog.LevelVar)
	l.Set(slog.LevelDebug)

	s := &stlog{
		sl: slog.New(slog.NewTextHandler(os.Stderr,
			&slog.HandlerOptions{Level: l})),
	}

	l.Set(parseLevel(level))

	return s
}

// newStlogJSON returns a new Logger with JSON formatting at the specified
// level.
func newStlogJSON(level string) Logger {
	l := new(slog.LevelVar)
	l.Set(slog.LevelDebug)

	s := &stlog{
		sl: slog.New(slog.NewJSONHandler(os.Stderr,
			&slog.HandlerOptions{Level: l})),
	}

	l.Set(parseLevel(level))

	return s
}

// WithField returns a derived Logger with a string field.
func (s *stlog) WithField(key, val string) Logger {
	return &stlog{sl: s.sl.With(key, val)}
}

// Debug logs a new message with debug level.
func (s *stlog) Debug(v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelDebug, fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func (s *stlog) Debugf(format string, v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelDebug, fmt.Sprintf(format,
		v...))
}

// Info logs a new message with info level.
func (s *stlog) Info(v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func (s *stlog) Infof(format string, v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelInfo, fmt.Sprintf(format,
		v...))
}

// Error logs a new message with error level.
func (s *stlog) Error(v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelError, fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func (s *stlog) Errorf(format string, v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelError, fmt.Sprintf(format,
		v...))
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
func (s *stlog) Fatal(v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelError+4, fmt.Sprint(v...))
	os.Exit(1)
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func (s *stlog) Fatalf(format string, v ...any) {
	s.sl.LogAttrs(context.Background(), slog.LevelError+4, fmt.Sprintf(format,
		v...))
	os.Exit(1)
}
