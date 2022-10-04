package alog

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

// zlog contains methods to write logs using zerolog and implements the Logger
// interface.
type zlog struct {
	zl zerolog.Logger
}

// Verify zlog implements Logger.
var _ Logger = &zlog{}

// newZlogConsole returns a new Logger with console formatting at the debug
// level.
func newZlogConsole() Logger {
	cw := zerolog.ConsoleWriter{
		Out: os.Stdout,
		// VSCode does not support colors in the output channel:
		// https://github.com/Microsoft/vscode/issues/571
		NoColor:    true,
		TimeFormat: time.RFC3339,
	}

	return &zlog{
		zl: zerolog.New(cw).With().Timestamp().Logger().
			Level(zerolog.DebugLevel),
	}
}

// newZlogJSON returns a new Logger with JSON formatting at the debug level.
func newZlogJSON() Logger {
	return &zlog{
		zl: zerolog.New(os.Stderr).With().Timestamp().Logger().
			Level(zerolog.DebugLevel),
	}
}

// WithLevel returns a derived Logger with the level set to level.
func (z *zlog) WithLevel(level string) Logger {
	zlevel, err := zerolog.ParseLevel(strings.ToLower(level))
	if err != nil {
		z.zl.Error().Msgf("SetLevel unknown level, using INFO: %v", level)

		return &zlog{zl: z.zl.Level(zerolog.InfoLevel)}
	}

	return &zlog{zl: z.zl.Level(zlevel)}
}

// WithField returns a derived Logger with a string field.
func (z *zlog) WithField(key, val string) Logger {
	return &zlog{zl: z.zl.With().Str(key, val).Logger()}
}

// Debug logs a new message with debug level.
func (z *zlog) Debug(v ...interface{}) {
	z.zl.Debug().Msg(fmt.Sprint(v...))
}

// Debugf logs a new formatted message with debug level.
func (z *zlog) Debugf(format string, v ...interface{}) {
	z.zl.Debug().Msgf(format, v...)
}

// Info logs a new message with info level.
func (z *zlog) Info(v ...interface{}) {
	z.zl.Info().Msg(fmt.Sprint(v...))
}

// Infof logs a new formatted message with info level.
func (z *zlog) Infof(format string, v ...interface{}) {
	z.zl.Info().Msgf(format, v...)
}

// Error logs a new message with error level.
func (z *zlog) Error(v ...interface{}) {
	z.zl.Error().Msg(fmt.Sprint(v...))
}

// Errorf logs a new formatted message with error level.
func (z *zlog) Errorf(format string, v ...interface{}) {
	z.zl.Error().Msgf(format, v...)
}

// Fatal logs a new message with fatal level followed by a call to os.Exit(1).
func (z *zlog) Fatal(v ...interface{}) {
	z.zl.Fatal().Msg(fmt.Sprint(v...))
}

// Fatalf logs a new formatted message with fatal level followed by a call to
// os.Exit(1).
func (z *zlog) Fatalf(format string, v ...interface{}) {
	z.zl.Fatal().Msgf(format, v...)
}
