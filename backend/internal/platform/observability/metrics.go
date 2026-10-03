package observability

import (
	"context"
	"database/sql"
	"errors"
	"net"
	"net/http"
	"net/http/pprof"
	"strconv"
	"strings"
	"time"

	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type Metrics struct {
	Registry *prometheus.Registry

	HTTPRequests *prometheus.CounterVec
	HTTPDuration *prometheus.HistogramVec
	HTTPInflight prometheus.Gauge

	GRPCRequests *prometheus.CounterVec
	GRPCDuration *prometheus.HistogramVec

	GRPCClientRequests *prometheus.CounterVec
	GRPCClientDuration *prometheus.HistogramVec

	RedisCacheOperations *prometheus.CounterVec
}

func NewMetrics(service, version, commit string) *Metrics {
	registry := prometheus.NewRegistry()
	m := &Metrics{
		Registry: registry,
		HTTPRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "volha_http_requests_total",
			Help: "Total HTTP requests by service, method, normalized route, and status.",
		}, []string{"service", "method", "route", "status"}),
		HTTPDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "volha_http_request_duration_seconds",
			Help:    "HTTP request duration by service, method, normalized route, and status.",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "method", "route", "status"}),
		HTTPInflight: prometheus.NewGauge(prometheus.GaugeOpts{
			Name:        "volha_http_inflight_requests",
			Help:        "Currently in-flight HTTP requests.",
			ConstLabels: prometheus.Labels{"service": service},
		}),
		GRPCRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "volha_grpc_server_requests_total",
			Help: "Total gRPC server requests by service, method, and code.",
		}, []string{"service", "grpc_service", "grpc_method", "grpc_code"}),
		GRPCDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "volha_grpc_server_request_duration_seconds",
			Help:    "gRPC server request duration by service, method, and code.",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "grpc_service", "grpc_method", "grpc_code"}),
		GRPCClientRequests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "volha_grpc_client_requests_total",
			Help: "Total downstream gRPC client requests by service, method, and code.",
		}, []string{"service", "grpc_service", "grpc_method", "grpc_code"}),
		GRPCClientDuration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "volha_grpc_client_request_duration_seconds",
			Help:    "Downstream gRPC client request duration by service, method, and code.",
			Buckets: prometheus.DefBuckets,
		}, []string{"service", "grpc_service", "grpc_method", "grpc_code"}),
		RedisCacheOperations: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "volha_redis_cache_operations_total",
			Help: "Redis cache operations by service, cache, operation, and result.",
		}, []string{"service", "cache", "operation", "result"}),
	}
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		m.HTTPRequests,
		m.HTTPDuration,
		m.HTTPInflight,
		m.GRPCRequests,
		m.GRPCDuration,
		m.GRPCClientRequests,
		m.GRPCClientDuration,
		m.RedisCacheOperations,
		prometheus.NewGaugeFunc(prometheus.GaugeOpts{
			Name: "volha_build_info",
			Help: "Build metadata for the running Volha service.",
			ConstLabels: prometheus.Labels{
				"service": service,
				"version": version,
				"commit":  commit,
			},
		}, func() float64 { return 1 }),
	)
	return m
}

type MetricsServer struct {
	server *http.Server
	errc   chan error
}

func StartMetricsServer(_ context.Context, cfg Config, metrics *Metrics) (*MetricsServer, error) {
	if !cfg.MetricsEnabled {
		return nil, nil
	}
	if metrics == nil {
		return nil, errors.New("metrics registry is nil")
	}
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry, promhttp.HandlerOpts{}))
	if cfg.PprofEnabled {
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	}
	server := &http.Server{Addr: cfg.MetricsAddr, Handler: mux}
	listener, err := net.Listen("tcp", cfg.MetricsAddr)
	if err != nil {
		return nil, err
	}
	ms := &MetricsServer{server: server, errc: make(chan error, 1)}
	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			ms.errc <- err
		}
		close(ms.errc)
	}()
	return ms, nil
}

func (s *MetricsServer) Errors() <-chan error {
	if s == nil {
		return nil
	}
	return s.errc
}

func (s *MetricsServer) Shutdown(ctx context.Context) error {
	if s == nil {
		return nil
	}
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}

func EchoMiddleware(cfg Config, metrics *Metrics) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			route := normalizedRoute(c)
			ctx := c.Request().Context()

			start := time.Now()
			if metrics != nil {
				metrics.HTTPInflight.Inc()
				defer metrics.HTTPInflight.Dec()
			}
			err := next(c)
			if err != nil {
				c.Error(err)
			}
			statusCode := c.Response().Status
			statusLabel := strconv.Itoa(statusCode)
			duration := time.Since(start)

			if metrics != nil {
				metrics.HTTPRequests.WithLabelValues(cfg.Service, c.Request().Method, route, statusLabel).Inc()
				metrics.HTTPDuration.WithLabelValues(cfg.Service, c.Request().Method, route, statusLabel).Observe(duration.Seconds())
			}
			if !isNoiseRoute(route) && statusCode >= 500 {
				log.Error(ctx, "http request failed",
					"http_method", c.Request().Method,
					"http_route", route,
					"http_status", statusCode,
					"duration_ms", float64(duration.Microseconds())/1000,
				)
			}
			if cfg.SlowRequestThreshold > 0 && duration > cfg.SlowRequestThreshold {
				log.Warn(ctx, "slow http request",
					"http_method", c.Request().Method,
					"http_route", route,
					"http_status", statusCode,
					"duration_ms", float64(duration.Microseconds())/1000,
				)
			}
			return nil
		}
	}
}

func UnaryServerInterceptor(cfg Config, metrics *Metrics) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()
		resp, err := handler(ctx, req)
		code := status.Code(err).String()
		grpcService, grpcMethod := splitFullMethod(info.FullMethod)
		duration := time.Since(start)
		if metrics != nil {
			metrics.GRPCRequests.WithLabelValues(cfg.Service, grpcService, grpcMethod, code).Inc()
			metrics.GRPCDuration.WithLabelValues(cfg.Service, grpcService, grpcMethod, code).Observe(duration.Seconds())
		}
		if !strings.Contains(info.FullMethod, "grpc.health.v1.Health") && err != nil {
			log.Error(ctx, "grpc request failed",
				"grpc_service", grpcService,
				"grpc_method", grpcMethod,
				"grpc_code", code,
				"duration_ms", float64(duration.Microseconds())/1000,
			)
		}
		if !strings.Contains(info.FullMethod, "grpc.health.v1.Health") && cfg.SlowRequestThreshold > 0 && duration > cfg.SlowRequestThreshold {
			log.Warn(ctx, "slow grpc request",
				"grpc_service", grpcService,
				"grpc_method", grpcMethod,
				"grpc_code", code,
				"duration_ms", float64(duration.Microseconds())/1000,
			)
		}
		return resp, err
	}
}

func UnaryClientInterceptor(cfg Config, metrics *Metrics) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req any, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		code := status.Code(err).String()
		grpcService, grpcMethod := splitFullMethod(method)
		if metrics != nil {
			metrics.GRPCClientRequests.WithLabelValues(cfg.Service, grpcService, grpcMethod, code).Inc()
			metrics.GRPCClientDuration.WithLabelValues(cfg.Service, grpcService, grpcMethod, code).Observe(time.Since(start).Seconds())
		}
		return err
	}
}

func RegisterDBStats(metrics *Metrics, service string, db *sql.DB) {
	if metrics == nil || db == nil {
		return
	}
	metrics.Registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "volha_db_pool_open_connections",
		Help:        "Database pool open connections.",
		ConstLabels: prometheus.Labels{"service": service},
	}, func() float64 { return float64(db.Stats().OpenConnections) }))
	metrics.Registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "volha_db_pool_in_use_connections",
		Help:        "Database pool in-use connections.",
		ConstLabels: prometheus.Labels{"service": service},
	}, func() float64 { return float64(db.Stats().InUse) }))
	metrics.Registry.MustRegister(prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Name:        "volha_db_pool_idle_connections",
		Help:        "Database pool idle connections.",
		ConstLabels: prometheus.Labels{"service": service},
	}, func() float64 { return float64(db.Stats().Idle) }))
	metrics.Registry.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
		Name:        "volha_db_pool_wait_count_total",
		Help:        "Database pool wait count.",
		ConstLabels: prometheus.Labels{"service": service},
	}, func() float64 { return float64(db.Stats().WaitCount) }))
	metrics.Registry.MustRegister(prometheus.NewCounterFunc(prometheus.CounterOpts{
		Name:        "volha_db_pool_wait_duration_seconds_total",
		Help:        "Database pool wait duration.",
		ConstLabels: prometheus.Labels{"service": service},
	}, func() float64 { return db.Stats().WaitDuration.Seconds() }))
}

func ObserveRedisCache(metrics *Metrics, service, cache, operation, result string) {
	if metrics != nil {
		metrics.RedisCacheOperations.WithLabelValues(service, cache, operation, result).Inc()
	}
}

func normalizedRoute(c echo.Context) string {
	route := c.Path()
	if route == "" {
		route = "unknown"
	}
	return route
}

func isNoiseRoute(route string) bool {
	return route == "/health/live" || route == "/health/ready" || route == "/metrics"
}

func splitFullMethod(fullMethod string) (string, string) {
	trimmed := strings.TrimPrefix(fullMethod, "/")
	parts := strings.Split(trimmed, "/")
	if len(parts) != 2 {
		return "unknown", "unknown"
	}
	return parts[0], parts[1]
}
