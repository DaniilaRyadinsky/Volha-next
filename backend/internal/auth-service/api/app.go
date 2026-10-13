package api

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/config"
	"github.com/autumnterror/volha-backend/internal/auth-service/service"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/autumnterror/volha-backend/internal/platform/requestid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
)

type App struct {
	server *grpc.Server
	cfg    *config.Config
}

type HealthCheckFunc func(context.Context) error

func New(cfg *config.Config, auth *service.Service, ready HealthCheckFunc, obsCfg observability.Config, metrics *observability.Metrics) *App {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.UnaryServerInterceptor(),
			observability.UnaryServerInterceptor(obsCfg, metrics),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{MaxConnectionIdle: 0}),
	)
	Register(server, auth)
	healthpb.RegisterHealthServer(server, healthServer{ready: ready})
	return &App{server: server, cfg: cfg}
}

func (a *App) Run() error {
	const op = "auth.grpc.App.Run"
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return format.Error(op, err)
	}
	log.Info(context.Background(), "grpc server running", "operation", op, "port", a.cfg.Port)
	if err := a.server.Serve(listener); err != nil {
		return format.Error(op, err)
	}
	return nil
}

func (a *App) Stop(ctx context.Context) {
	done := make(chan struct{})
	go func() {
		a.server.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		a.server.Stop()
	}
	log.Info(context.Background(), "grpc server stopped", "port", a.cfg.Port)
}

type healthServer struct {
	healthpb.UnimplementedHealthServer
	ready HealthCheckFunc
}

func (s healthServer) Check(ctx context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	if s.ready == nil {
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
	}
	checkCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := s.ready(checkCtx); err != nil {
		return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_NOT_SERVING}, status.Error(codes.Unavailable, err.Error())
	}
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}
