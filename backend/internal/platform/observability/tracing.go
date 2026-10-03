package observability

import (
	"context"
	"errors"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials/insecure"
)

func SetupTracing(ctx context.Context, cfg Config) (func(context.Context) error, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))
	if !cfg.OTelEnabled {
		otel.SetTracerProvider(sdktrace.NewTracerProvider())
		return func(context.Context) error { return nil }, nil
	}
	setupCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	exporter, err := otlptracegrpc.New(
		setupCtx,
		otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
		otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	res, err := resource.New(ctx,
		resource.WithAttributes(
			attribute.String("service.name", cfg.Service),
			attribute.String("deployment.environment", cfg.Environment),
			attribute.String("service.version", cfg.Version),
		),
	)
	if err != nil {
		return nil, err
	}
	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.TraceSamplingRatio))),
	)
	otel.SetTracerProvider(provider)
	return provider.Shutdown, nil
}

func ValidateSamplingRatio(ratio float64) error {
	if ratio < 0 || ratio > 1 {
		return errors.New("trace sampling ratio must be between 0 and 1")
	}
	return nil
}
