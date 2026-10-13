package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/autumnterror/volha-backend/docs"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	authgrpc "github.com/autumnterror/volha-backend/internal/gateway/grpc/auth"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	"github.com/autumnterror/volha-backend/internal/gateway/net/echo"
	"github.com/autumnterror/volha-backend/internal/gateway/redis"
	"github.com/autumnterror/volha-backend/internal/platform/buildinfo"
	"github.com/autumnterror/volha-backend/internal/platform/lifecycle"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
)

// @title Volha gateway REST API
// @version 0.1

// @contact.name Alex "bustard" Provor
// @contact.url https://breezy.su
// @contact.email help@breezy.su

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	const op = "main"

	configPath := flag.String("config", "", "path to gateway config file")
	flag.Parse()
	if *configPath != "" {
		_ = os.Setenv("CONFIG_PATH", *configPath)
	}
	log.Configure(log.Config{
		Service:     "volha-gateway",
		Env:         os.Getenv("APP_ENV"),
		Version:     buildinfo.Version,
		Level:       os.Getenv("LOG_LEVEL"),
		Format:      os.Getenv("LOG_FORMAT"),
		Development: os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local",
	})

	cfg, err := config.Setup()
	if err != nil {
		log.Error(context.Background(), "gateway config setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	obsCfg := observability.FromEnv("volha-gateway", cfg.Mode, ":9091")
	log.Configure(log.Config{
		Service:     obsCfg.Service,
		Env:         obsCfg.Environment,
		Version:     obsCfg.Version,
		Level:       os.Getenv("LOG_LEVEL"),
		Format:      os.Getenv("LOG_FORMAT"),
		Development: cfg.Mode == "DEV" || cfg.Mode == "local" || cfg.Mode == "development",
	})
	log.Info(runCtx, "service starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
		"build_time", buildinfo.BuildTime,
		"environment", obsCfg.Environment,
	)
	shutdownTracing, err := observability.SetupTracing(runCtx, obsCfg)
	if err != nil {
		log.Error(runCtx, "tracing setup failed", "operation", op, "error", err)
		shutdownTracing = func(context.Context) error { return nil }
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracing(ctx); err != nil {
			log.Error(ctx, "tracing shutdown failed", "operation", op, "error", err)
		}
	}()

	metrics := observability.NewMetrics(obsCfg.Service, obsCfg.Version, obsCfg.Commit)
	metricsServer, err := observability.StartMetricsServer(runCtx, obsCfg, metrics)
	if err != nil {
		log.Error(runCtx, "metrics setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	var metricsErr <-chan error
	if metricsServer != nil {
		metricsErr = metricsServer.Errors()
	}

	p, err := products.New(cfg, obsCfg, metrics)
	if err != nil {
		log.Error(runCtx, "product client setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := p.Close(); err != nil {
			log.Error(context.Background(), "grpc product client close failed", "operation", op, "error", err)
		}
	}()
	authClient, err := authgrpc.New(cfg, obsCfg, metrics)
	if err != nil {
		log.Error(runCtx, "auth client setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := authClient.Close(); err != nil {
			log.Error(context.Background(), "grpc auth client close failed", "operation", op, "error", err)
		}
	}()

	rds, err := redis.New(cfg, metrics)
	if err != nil {
		log.Error(runCtx, "redis setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := rds.Close(); err != nil {
			log.Error(context.Background(), "redis close failed", "operation", op, "error", err)
		}
	}()

	e := echo.New(rds, p, authClient, cfg, obsCfg, metrics)
	serverErr := make(chan error, 1)
	go func() {
		serverErr <- e.Run()
	}()

	termination := lifecycle.Wait(runCtx, serverErr, metricsErr)
	switch {
	case termination.Err != nil && termination.Source == "server":
		log.Error(runCtx, "http server failed", "operation", op, "error", termination.Err)
	case termination.Err != nil && termination.Source == "metrics":
		log.Error(runCtx, "metrics server failed", "operation", op, "error", termination.Err)
	case termination.Source == "signal":
		log.Info(runCtx, "shutdown requested", "operation", op, "reason", termination.Err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Stop(shutdownCtx); err != nil {
		log.Error(shutdownCtx, "http shutdown failed", "operation", op, "error", err)
	} else {
		log.Info(shutdownCtx, "http server stopped", "operation", op)
	}
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error(shutdownCtx, "metrics shutdown failed", "operation", op, "error", err)
	} else {
		log.Info(shutdownCtx, "metrics server stopped", "operation", op)
	}

	log.Info(context.Background(), "service stopped", "operation", op, "reason", termination.Source)
}
