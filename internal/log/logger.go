package logger

import (
	"context"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Init configures a global Zap logger. If jsonOutput is true, logs are emitted
// as JSON; otherwise, a console (human-friendly) encoder is used.
//
// Levels supported (case-insensitive): "debug", "info", "warn", "error".
func Init(level string, jsonOutput bool) *zap.Logger {
	encCfg := zap.NewProductionEncoderConfig()
	encCfg.TimeKey = "ts"
	encCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	if !jsonOutput {
		// Console encoder with readable level names
		encCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(parseLevel(level)),
		Development:      false,
		Encoding:         map[bool]string{true: "json", false: "console"}[jsonOutput],
		EncoderConfig:    encCfg,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	l, err := cfg.Build()
	if err != nil {
		// Fallback to a no-op logger to avoid panics
		l = zap.NewNop()
	}
	zap.ReplaceGlobals(l)
	return l
}

// L returns the process-wide global logger set by Init.
func L() *zap.Logger { return zap.L() }

// With returns a child logger with constant fields attached.
// Example: logger.With("request_id", id)
func With(args ...any) *zap.Logger { return zap.L().With(argsToFields(args...)...) }

// Context helpers to carry a request-scoped logger.
type ctxKey struct{}

// IntoContext stores the provided logger in the context and returns the new context.
func IntoContext(ctx context.Context, l *zap.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// FromContext retrieves the logger from the context, falling back to the global
// logger if one is not present.
func FromContext(ctx context.Context) *zap.Logger {
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return zap.L()
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(strings.TrimSpace(level)) {
	case "debug":
		return zapcore.DebugLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default: // info + unknowns
		return zapcore.InfoLevel
	}
}

func argsToFields(args ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i+1 < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, zap.Any(key, args[i+1]))
	}
	return fields
}
