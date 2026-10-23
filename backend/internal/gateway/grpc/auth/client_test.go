package auth

import (
	"context"
	"net"
	"testing"
	"time"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
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
	productsRPC.RegisterAuthServiceServer(server, authServer{})
	go func() { _ = server.Serve(listener) }()
	t.Cleanup(server.Stop)
	cfg := &config.Config{AddrAuth: listener.Addr().String(), Timeout: time.Second}
	client, err := New(cfg, observability.Config{}, observability.NewMetrics("gateway-test", "dev", "test"))
	require.NoError(t, err)
	require.True(t, client.Ready(context.Background()))
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	require.False(t, client.Ready(context.Background()))
	principal, err := client.ValidateAccessToken(context.Background(), "access")
	require.NoError(t, err)
	require.Equal(t, "user", principal.GetUser().GetId())
	require.NoError(t, client.Close())
	require.False(t, client.Ready(context.Background()))
}

type authServer struct {
	productsRPC.UnimplementedAuthServiceServer
}

func (authServer) Authenticate(ctx context.Context, req *productsRPC.AuthenticateRequest) (*productsRPC.AuthenticatedUser, error) {
	return &productsRPC.AuthenticatedUser{User: &productsRPC.AuthUser{Id: "user"}}, nil
}
