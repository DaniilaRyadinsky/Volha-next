package api

import (
	"context"
	"fmt"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/autumnterror/volha-backend/internal/platform/requestid"
	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/status"
	"net"
	"time"
)

type App struct {
	gRPCServer *grpc.Server
	cfg        *config.Config
}

type HealthCheckFunc func(ctx context.Context) error

// New construct new App structure
func New(
	cfg *config.Config,
	ps *service.ProductsService,
	ready HealthCheckFunc,
	obsCfg observability.Config,
	metrics *observability.Metrics,
) *App {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			requestid.UnaryServerInterceptor(),
			observability.UnaryServerInterceptor(obsCfg, metrics),
		),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: 0,
		}),
	)
	Register(s, ps)
	healthpb.RegisterHealthServer(s, dependencyHealthServer{ready: ready})
	return &App{
		gRPCServer: s,
		cfg:        cfg,
	}
}

func (a *App) Run() error {
	const op = "grpc.App"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.cfg.Port))
	if err != nil {
		return format.Error(op, err)
	}
	log.Info(context.Background(), "grpc server running", "operation", op, "port", a.cfg.Port)

	if err := a.gRPCServer.Serve(l); err != nil {
		return format.Error(op, err)
	}

	return nil
}

func (a *App) Stop(ctx context.Context) {
	const op = "grpc.Stop"
	done := make(chan struct{})
	go func() {
		a.gRPCServer.GracefulStop()
		close(done)
	}()
	select {
	case <-done:
	case <-ctx.Done():
		a.gRPCServer.Stop()
	}
	log.Info(context.Background(), "grpc server stopped", "operation", op, "port", a.cfg.Port)
}

type dependencyHealthServer struct {
	healthpb.UnimplementedHealthServer
	ready HealthCheckFunc
}

func (s dependencyHealthServer) Check(ctx context.Context, _ *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
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
