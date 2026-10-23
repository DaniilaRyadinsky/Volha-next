package logger

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggingLevelsAndArguments(t *testing.T) {
	old := global.Load()
	t.Cleanup(func() { global.Store(old) })
	var buffer bytes.Buffer
	global.Store(slog.New(slog.NewJSONHandler(&buffer, &slog.HandlerOptions{Level: slog.LevelDebug})))
	Debug(context.Background(), "debug", "number", 1)
	Info(context.Background(), "info", slog.String("field", "value"))
	Warn(context.Background(), "warning", "odd")
	Error(context.Background(), "error", 99)
	require.Contains(t, buffer.String(), `"number":1`)
	require.Contains(t, buffer.String(), `"field":"value"`)
	require.Contains(t, buffer.String(), "invalid_log_arg")
	for input, want := range map[string]slog.Level{"debug": slog.LevelDebug, "warning": slog.LevelWarn, "warn": slog.LevelWarn, "error": slog.LevelError, "unknown": slog.LevelInfo} {
		require.Equal(t, want, parseLevel(input))
	}
	Configure(Config{Service: "test", Env: "unit", Format: "json"})
	require.NotNil(t, L())
	Configure(Config{Development: true})
	require.NotNil(t, L())
	require.Equal(t, "fallback", defaultString("", "fallback"))
}
