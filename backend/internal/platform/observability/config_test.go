package observability

import (
	"context"
	"testing"
)

func TestValidateSamplingRatio(t *testing.T) {
	t.Parallel()

	if err := ValidateSamplingRatio(0.5); err != nil {
		t.Fatalf("expected valid ratio, got %v", err)
	}
	if err := ValidateSamplingRatio(1.5); err == nil {
		t.Fatal("expected invalid ratio")
	}
}

func TestSetupTracingDisabled(t *testing.T) {
	t.Parallel()

	shutdown, err := SetupTracing(context.Background(), Config{
		Service:            "test",
		Environment:        "test",
		Version:            "test",
		MetricsEnabled:     false,
		OTelEnabled:        false,
		TraceSamplingRatio: 0.1,
	})
	if err != nil {
		t.Fatalf("SetupTracing() error = %v", err)
	}
	if shutdown == nil {
		t.Fatal("expected shutdown function")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown error = %v", err)
	}
}
