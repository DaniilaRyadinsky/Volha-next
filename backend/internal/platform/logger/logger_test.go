package logger

import (
	"context"
	"testing"

	"github.com/autumnterror/volha-backend/internal/platform/requestid"
	"go.opentelemetry.io/otel/trace"
)

func TestContextAttrsIncludeRequestAndTraceIDs(t *testing.T) {
	t.Parallel()

	traceID := trace.TraceID{1, 2, 3}
	spanID := trace.SpanID{4, 5, 6}
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
		Remote:  true,
	})
	ctx := requestid.WithContext(context.Background(), "req-123")
	ctx = trace.ContextWithSpanContext(ctx, spanCtx)

	attrs := contextAttrs(ctx)
	got := map[string]string{}
	for _, attr := range attrs {
		got[attr.Key] = attr.Value.String()
	}

	if got["request_id"] != "req-123" {
		t.Fatalf("request_id not preserved: %q", got["request_id"])
	}
	if got["trace_id"] != traceID.String() {
		t.Fatalf("trace_id not extracted: %q", got["trace_id"])
	}
	if got["span_id"] != spanID.String() {
		t.Fatalf("span_id not extracted: %q", got["span_id"])
	}
}
