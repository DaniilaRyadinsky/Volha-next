package observability

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGRPCAndDependencyMetrics(t *testing.T) {
	m := NewMetrics("test", "v1", "commit")
	cfg := Config{Service: "test", SlowRequestThreshold: time.Hour}
	want := status.Error(codes.NotFound, "missing")
	_, err := UnaryServerInterceptor(cfg, m)(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/products.Products/GetProduct"}, func(context.Context, any) (any, error) { return nil, want })
	require.ErrorIs(t, err, want)
	err = UnaryClientInterceptor(cfg, m)(context.Background(), "/products.Products/GetProduct", nil, nil, nil, func(context.Context, string, any, any, *grpc.ClientConn, ...grpc.CallOption) error { return want })
	require.ErrorIs(t, err, want)
	ObserveRedisCache(m, "test", "dictionaries", "get", "hit")
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	RegisterDBStats(m, "test", db)
	families, err := m.Registry.Gather()
	require.NoError(t, err)
	names := map[string]bool{}
	for _, family := range families {
		names[family.GetName()] = true
		if family.GetName() == "volha_grpc_server_requests_total" {
			require.Equal(t, float64(1), family.Metric[0].GetCounter().GetValue())
		}
	}
	require.True(t, names["volha_grpc_server_requests_total"])
	require.True(t, names["volha_grpc_client_requests_total"])
	require.True(t, names["volha_redis_cache_operations_total"])
	ObserveRedisCache(nil, "test", "cache", "get", "error")
	RegisterDBStats(nil, "test", db)
}
func TestMetricsServerLifecycleAndConfig(t *testing.T) {
	m := NewMetrics("test", "dev", "c")
	server, err := StartMetricsServer(context.Background(), Config{MetricsEnabled: true, MetricsAddr: "127.0.0.1:0", PprofEnabled: true}, m)
	require.NoError(t, err)
	require.NotNil(t, server)
	require.NoError(t, server.Shutdown(context.Background()))
	select {
	case err := <-server.Errors():
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("server did not stop")
	}
	server, err = StartMetricsServer(context.Background(), Config{}, nil)
	require.NoError(t, err)
	require.Nil(t, server)
	require.NoError(t, server.Shutdown(context.Background()))
	require.Nil(t, server.Errors())
	_, err = StartMetricsServer(context.Background(), Config{MetricsEnabled: true, MetricsAddr: "127.0.0.1:-1"}, m)
	require.Error(t, err)
	_, err = StartMetricsServer(context.Background(), Config{MetricsEnabled: true}, nil)
	require.Error(t, err)
	for _, cfg := range []Config{{TraceSamplingRatio: -1}, {TraceSamplingRatio: 2}, {MetricsEnabled: true}, {OTelEnabled: true}} {
		require.Error(t, cfg.Validate())
	}
	t.Setenv("METRICS_ENABLED", "bad")
	t.Setenv("OTEL_TRACES_SAMPLER_RATIO", "bad")
	t.Setenv("SLOW_REQUEST_THRESHOLD", "bad")
	cfg := FromEnv("service", "test", ":1234")
	require.True(t, cfg.MetricsEnabled)
	require.Equal(t, 0.1, cfg.TraceSamplingRatio)
	require.Equal(t, time.Second, cfg.SlowRequestThreshold)
	_, err = SetupTracing(context.Background(), Config{TraceSamplingRatio: -1})
	require.Error(t, err)
}
