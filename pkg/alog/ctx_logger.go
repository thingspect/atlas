package alog

import "context"

// CtxLogger wraps a Logger for use in Context values. This allows passing a
// consistent struct while still supporting method chaining.
type CtxLogger struct {
	Logger
}

// loggerKey is the key for CtxLogger values in Contexts. It is unexported,
// clients should use NewContext and FromContext instead of using this key
// directly.
type loggerKey struct{}

// NewContext returns a new Context that carries a CtxLogger.
func NewContext(ctx context.Context, logger *CtxLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// FromContext returns the CtxLogger value stored in a Context, if any. If a
// CtxLogger is not present, one carrying the global logger will be returned.
func FromContext(ctx context.Context) *CtxLogger {
	logger, ok := ctx.Value(loggerKey{}).(*CtxLogger)
	if !ok {
		return &CtxLogger{Logger: Global()}
	}

	return logger
}
