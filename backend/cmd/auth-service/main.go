package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/api"
	"github.com/autumnterror/volha-backend/internal/auth-service/config"
	"github.com/autumnterror/volha-backend/internal/auth-service/infra/psql"
	"github.com/autumnterror/volha-backend/internal/auth-service/repository"
	"github.com/autumnterror/volha-backend/internal/auth-service/service"
	"github.com/autumnterror/volha-backend/internal/platform/buildinfo"
	"github.com/autumnterror/volha-backend/internal/platform/lifecycle"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
)

func main() {
	const op = "auth.main"
	configPath := flag.String("config", "", "path to auth-service config file")
	flag.Parse()
	if *configPath != "" {
		_ = os.Setenv("CONFIG_PATH", *configPath)
	}
	log.Configure(log.Config{
		Service:     "volha-auth-service",
		Env:         os.Getenv("APP_ENV"),
		Version:     buildinfo.Version,
		Level:       os.Getenv("LOG_LEVEL"),
		Format:      os.Getenv("LOG_FORMAT"),
		Development: os.Getenv("APP_ENV") == "" || os.Getenv("APP_ENV") == "local",
	})

	cfg, err := config.Setup()
	if err != nil {
		log.Error(context.Background(), "auth service config setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	runCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	obsCfg := observability.FromEnv("volha-auth-service", cfg.Mode, ":9093")
	log.Configure(log.Config{
		Service:     obsCfg.Service,
		Env:         obsCfg.Environment,
		Version:     obsCfg.Version,
		Level:       os.Getenv("LOG_LEVEL"),
		Format:      os.Getenv("LOG_FORMAT"),
		Development: cfg.Mode == "DEV" || cfg.Mode == "local" || cfg.Mode == "development",
	})
	log.Info(runCtx, "service starting", "version", buildinfo.Version, "commit", buildinfo.Commit, "build_time", buildinfo.BuildTime, "environment", obsCfg.Environment)
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

	db, err := psql.Connect(cfg)
	if err != nil {
		log.Error(runCtx, "postgres setup failed", "operation", op, "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Error(context.Background(), "postgres disconnect failed", "operation", op, "error", err)
		}
	}()
	observability.RegisterDBStats(metrics, obsCfg.Service, db.Driver)

	authService := service.New(repository.New(db.Driver), cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)
	a := api.New(cfg, authService, db.Ping, obsCfg, metrics)
	serverErr := make(chan error, 1)
	go func() { serverErr <- a.Run() }()
	log.Info(runCtx, "grpc server starting", "operation", op)

	termination := lifecycle.Wait(runCtx, serverErr, metricsErr)
	switch {
	case termination.Err != nil && termination.Source == "server":
		log.Error(runCtx, "grpc server failed", "operation", op, "error", termination.Err)
	case termination.Err != nil && termination.Source == "metrics":
		log.Error(runCtx, "metrics server failed", "operation", op, "error", termination.Err)
	case termination.Source == "signal":
		log.Info(runCtx, "shutdown requested", "operation", op, "reason", termination.Err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	a.Stop(shutdownCtx)
	if err := metricsServer.Shutdown(shutdownCtx); err != nil {
		log.Error(shutdownCtx, "metrics server shutdown failed", "operation", op, "error", err)
	}
	log.Info(context.Background(), "service stopped", "operation", op, "reason", termination.Source)
}
