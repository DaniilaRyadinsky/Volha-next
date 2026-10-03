package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func TestEchoMiddlewareUsesNormalizedRouteMetricLabel(t *testing.T) {
	t.Parallel()

	e := echo.New()
	metrics := NewMetrics("volha-gateway", "dev", "unknown")
	cfg := Config{Service: "volha-gateway", SlowRequestThreshold: 0}
	e.Use(EchoMiddleware(cfg, metrics))
	e.GET("/api/product/:id", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/api/product/abc123?raw=value", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}

	families, err := metrics.Registry.Gather()
	if err != nil {
		t.Fatalf("gather metrics: %v", err)
	}
	var routeLabels []string
	for _, family := range families {
		if family.GetName() != "volha_http_requests_total" {
			continue
		}
		for _, metric := range family.GetMetric() {
			routeLabels = append(routeLabels, labelValue(metric, "route"))
		}
	}
	if len(routeLabels) != 1 {
		t.Fatalf("expected one HTTP request metric, got %d", len(routeLabels))
	}
	if routeLabels[0] != "/api/product/:id" {
		t.Fatalf("route label must be normalized, got %q", routeLabels[0])
	}
}

func TestOTelEchoMiddlewareExtractsTraceparent(t *testing.T) {
	oldTP := otel.GetTracerProvider()
	oldPropagator := otel.GetTextMapPropagator()
	tp := trace.NewTracerProvider(trace.WithSampler(trace.AlwaysSample()))
	defer tp.Shutdown(t.Context())
	defer otel.SetTracerProvider(oldTP)
	defer otel.SetTextMapPropagator(oldPropagator)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	const traceID = "0af7651916cd43dd8448eb211c80319c"
	e := echo.New()
	e.Use(otelecho.Middleware("volha-gateway"))
	e.GET("/trace", func(c echo.Context) error {
		sc := oteltrace.SpanContextFromContext(c.Request().Context())
		if !sc.IsValid() {
			t.Fatal("expected valid span context")
		}
		if got := sc.TraceID().String(); got != traceID {
			t.Fatalf("unexpected trace id: %s", got)
		}
		return c.NoContent(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/trace", nil)
	req.Header.Set("traceparent", "00-"+traceID+"-b7ad6b7169203331-01")
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", rec.Code)
	}
}

func labelValue(metric *io_prometheus_client.Metric, name string) string {
	for _, label := range metric.GetLabel() {
		if label.GetName() == name {
			return label.GetValue()
		}
	}
	return ""
}
