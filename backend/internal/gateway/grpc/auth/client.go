package auth

import (
	"context"
	"fmt"
	"time"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/autumnterror/volha-backend/internal/platform/requestid"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Client struct {
	API    productsRPC.AuthServiceClient
	Health healthpb.HealthClient
	cc     *grpc.ClientConn
	cfg    *config.Config
}

func New(cfg *config.Config, obsCfg observability.Config, metrics *observability.Metrics) (*Client, error) {
	const op = "grpc.auth.New"
	cc, err := grpc.NewClient(
		cfg.AddrAuth,
		grpc.WithChainUnaryInterceptor(
			requestid.UnaryClientInterceptor(),
			observability.UnaryClientInterceptor(obsCfg, metrics),
		),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, format.Error(op, err)
	}
	log.Info(context.Background(), "connecting to auth service", "target", cfg.AddrAuth)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cc.Connect()
	for {
		state := cc.GetState()
		if state == connectivity.Ready {
			log.Info(context.Background(), "connected to auth service", "target", cfg.AddrAuth)
			break
		}
		if !cc.WaitForStateChange(ctx, state) {
			_ = cc.Close()
			return nil, fmt.Errorf("auth service connection timeout")
		}
	}
	return &Client{API: productsRPC.NewAuthServiceClient(cc), Health: healthpb.NewHealthClient(cc), cc: cc, cfg: cfg}, nil
}

func (c *Client) Ready(ctx context.Context) bool {
	if c.cc.GetState() != connectivity.Ready {
		return false
	}
	timeout := c.cfg.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	response, err := c.Health.Check(checkCtx, &healthpb.HealthCheckRequest{})
	return err == nil && response.GetStatus() == healthpb.HealthCheckResponse_SERVING
}

func (c *Client) ValidateAccessToken(ctx context.Context, token string) (*productsRPC.AuthenticatedUser, error) {
	timeout := c.cfg.Timeout
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	checkCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return c.API.Authenticate(checkCtx, &productsRPC.AuthenticateRequest{AccessToken: token})
}

func (c *Client) Close() error {
	return c.cc.Close()
}
