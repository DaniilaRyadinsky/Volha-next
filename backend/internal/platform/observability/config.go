package observability

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/autumnterror/volha-backend/internal/platform/buildinfo"
)

type Config struct {
	Service              string
	Environment          string
	Version              string
	Commit               string
	BuildTime            string
	MetricsEnabled       bool
	MetricsAddr          string
	PprofEnabled         bool
	OTelEnabled          bool
	OTLPEndpoint         string
	TraceSamplingRatio   float64
	SlowRequestThreshold time.Duration
}

func FromEnv(service, environment, defaultMetricsAddr string) Config {
	return Config{
		Service:              service,
		Environment:          environment,
		Version:              buildinfo.Version,
		Commit:               buildinfo.Commit,
		BuildTime:            buildinfo.BuildTime,
		MetricsEnabled:       envBool("METRICS_ENABLED", true),
		MetricsAddr:          env("METRICS_ADDR", defaultMetricsAddr),
		PprofEnabled:         envBool("PPROF_ENABLED", false),
		OTelEnabled:          envBool("OTEL_ENABLED", false),
		OTLPEndpoint:         env("OTEL_EXPORTER_OTLP_ENDPOINT", "alloy:4317"),
		TraceSamplingRatio:   envFloat("OTEL_TRACES_SAMPLER_RATIO", 0.1),
		SlowRequestThreshold: envDuration("SLOW_REQUEST_THRESHOLD", time.Second),
	}
}

func (c Config) Validate() error {
	if c.TraceSamplingRatio < 0 || c.TraceSamplingRatio > 1 {
		return errors.New("trace sampling ratio must be between 0 and 1")
	}
	if c.MetricsEnabled && c.MetricsAddr == "" {
		return errors.New("metrics address is required when metrics are enabled")
	}
	if c.OTelEnabled && c.OTLPEndpoint == "" {
		return errors.New("OTLP endpoint is required when tracing is enabled")
	}
	return nil
}

func env(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envBool(name string, fallback bool) bool {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}

func envFloat(name string, fallback float64) float64 {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fallback
	}
	return parsed
}

func envDuration(name string, fallback time.Duration) time.Duration {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		return fallback
	}
	return parsed
}
