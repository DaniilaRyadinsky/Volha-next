package products

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func TestClientConnectionHealthAndClose(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	server := grpc.NewServer()
	healthServer := health.NewServer()
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)
	healthpb.RegisterHealthServer(server, healthServer)
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	cfg := &config.Config{AddrProducts: listener.Addr().String(), Timeout: time.Second}
	client, err := New(cfg, observability.Config{}, observability.NewMetrics("gateway-test", "dev", "test"))
	require.NoError(t, err)
	require.True(t, client.Ready(context.Background()))
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	require.False(t, client.Ready(context.Background()))
	require.NoError(t, client.Close())
	require.False(t, client.Ready(context.Background()))
}
