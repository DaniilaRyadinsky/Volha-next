package config

import (
	"os"
	"strconv"
	"time"

	platformpostgres "github.com/autumnterror/volha-backend/internal/platform/postgres"
)

func Test() *Config {
	host := envOr("TEST_POSTGRES_HOST", envOr("POSTGRES_HOST", "localhost"))
	port := envIntOr("TEST_POSTGRES_PORT", envIntOr("POSTGRES_PORT", 5432))
	user := envOr("TEST_POSTGRES_USER", os.Getenv("POSTGRES_USER"))
	password := envOr("TEST_POSTGRES_PASSWORD", os.Getenv("POSTGRES_PASSWORD"))
	db := envOr("TEST_POSTGRES_DB", os.Getenv("POSTGRES_DB"))

	return &Config{
		ConnStr:             platformpostgres.URL(host, port, db, user, password, "disable"),
		Port:                8080,
		Mode:                "TEST",
		MaxOpenConns:        5,
		MaxIdleConns:        2,
		ConnMaxLifetime:     time.Minute,
		ConnMaxIdleTime:     time.Minute,
		PostgresPingTimeout: time.Second,
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func envIntOr(name string, fallback int) int {
	value := os.Getenv(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
