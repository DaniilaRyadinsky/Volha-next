package config

import (
	"os"
	"time"
)

func Test() *Config {
	return &Config{
		AddrProducts:  envOr("TEST_PRODUCT_GRPC_ADDR", envOr("PRODUCT_GRPC_ADDR", "localhost:8008")),
		Timeout:       10 * time.Second,
		Port:          8080,
		Mode:          "DEV",
		DevAdminToken: envOr("TEST_DEV_ADMIN_TOKEN", os.Getenv("DEV_ADMIN_TOKEN")),

		RedisAddr:    envOr("TEST_REDIS_ADDR", envOr("REDIS_ADDR", "localhost:6379")),
		RedisPw:      envOr("TEST_REDIS_PASSWORD", os.Getenv("REDIS_PASSWORD")),
		RedisDictTTL: time.Minute,
	}
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
