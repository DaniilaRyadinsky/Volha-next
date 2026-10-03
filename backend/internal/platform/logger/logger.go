package logger

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync/atomic"

	"github.com/autumnterror/volha-backend/internal/platform/requestid"
	"go.opentelemetry.io/otel/trace"
)

type Config struct {
	Service     string
	Env         string
	Version     string
	Level       string
	Format      string
	Development bool
}

var global atomic.Pointer[slog.Logger]

func Configure(cfg Config) {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.Level)}
	format := strings.ToLower(cfg.Format)
	if format == "" {
		if cfg.Development {
			format = "text"
		} else {
			format = "json"
		}
	}

	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	service := defaultString(cfg.Service, "volha")
	env := defaultString(cfg.Env, "local")
	version := defaultString(cfg.Version, "dev")
	logger := slog.New(handler).With(
		slog.String("service", service),
		slog.String("environment", env),
		slog.String("version", version),
	)
	global.Store(logger)
	slog.SetDefault(logger)
}

func Debug(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelDebug, msg, args...)
}

func Info(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelInfo, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelWarn, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
	log(ctx, slog.LevelError, msg, args...)
}

func L() *slog.Logger {
	if l := global.Load(); l != nil {
		return l
	}
	Configure(Config{Development: true})
	return global.Load()
}

func log(ctx context.Context, level slog.Level, msg string, args ...any) {
	L().LogAttrs(ctx, level, msg, append(contextAttrs(ctx), slogAttrs(args...)...)...)
}

func contextAttrs(ctx context.Context) []slog.Attr {
	if ctx == nil {
		return nil
	}
	attrs := make([]slog.Attr, 0, 3)
	if id := requestid.FromContext(ctx); id != "" {
		attrs = append(attrs, slog.String("request_id", id))
	}
	span := trace.SpanContextFromContext(ctx)
	if span.IsValid() {
		attrs = append(attrs, slog.String("trace_id", span.TraceID().String()))
		attrs = append(attrs, slog.String("span_id", span.SpanID().String()))
	}
	return attrs
}

func slogAttrs(args ...any) []slog.Attr {
	attrs := make([]slog.Attr, 0, len(args)/2)
	for i := 0; i < len(args); {
		if attr, ok := args[i].(slog.Attr); ok {
			attrs = append(attrs, attr)
			i++
			continue
		}
		key, ok := args[i].(string)
		if !ok || i+1 >= len(args) {
			attrs = append(attrs, slog.Any("invalid_log_arg", args[i]))
			i++
			continue
		}
		attrs = append(attrs, slog.Any(key, args[i+1]))
		i += 2
	}
	return attrs
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func defaultString(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
