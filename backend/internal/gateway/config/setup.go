package config

import (
	"errors"
	"os"
	"time"

	platformconfig "github.com/autumnterror/volha-backend/internal/platform/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/spf13/viper"
)

type Config struct {
	AddrProducts  string        `mapstructure:"addr_products"`
	Timeout       time.Duration `mapstructure:"timeout"`
	Port          int           `mapstructure:"port"`
	Mode          string        `mapstructure:"mode"`
	RedisPw       string        `mapstructure:"redis_pw"`
	RedisAddr     string        `mapstructure:"redis_addr"`
	RedisDictTTL  time.Duration `mapstructure:"redis_dictionaries_ttl"`
	DevAdminToken string        `mapstructure:"dev_admin_token"`
}

// Setup creates config structure.
func Setup() (*Config, error) {
	const op = "config.setup"

	platformconfig.LoadDotEnv()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/volha-gateway.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	var cfg Config
	if err := v.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}

	var err error
	cfg.AddrProducts = platformconfig.StringFromEnv("PRODUCT_GRPC_ADDR", cfg.AddrProducts)
	cfg.RedisAddr = platformconfig.StringFromEnv("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisPw = platformconfig.StringFromEnv("REDIS_PASSWORD", "")
	cfg.DevAdminToken = platformconfig.StringFromEnv("DEV_ADMIN_TOKEN", cfg.DevAdminToken)
	cfg.Mode = platformconfig.StringFromEnv("APP_ENV", cfg.Mode)
	if cfg.Port, err = platformconfig.IntFromEnv("GATEWAY_PORT", cfg.Port); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.Timeout, err = platformconfig.DurationFromEnv("PRODUCT_GRPC_TIMEOUT", cfg.Timeout); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.RedisDictTTL, err = platformconfig.DurationFromEnv("REDIS_DICTIONARIES_TTL", defaultDuration(cfg.RedisDictTTL, 10*time.Minute)); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.AddrProducts == "" || cfg.Port == 0 || cfg.RedisAddr == "" {
		return nil, format.Error(op, errors.New("missing required gateway config"))
	}
	if isDevelopment(cfg.Mode) && cfg.DevAdminToken == "" {
		return nil, format.Error(op, errors.New("DEV_ADMIN_TOKEN is required in development mode"))
	}
	return &cfg, nil
}

func isDevelopment(mode string) bool {
	return mode == "local" || mode == "development" || mode == "DEV"
}

func defaultDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value == 0 {
		return fallback
	}
	return value
}
