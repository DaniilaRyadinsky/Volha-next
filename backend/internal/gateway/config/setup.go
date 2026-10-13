package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	platformconfig "github.com/autumnterror/volha-backend/internal/platform/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/spf13/viper"
)

type Config struct {
	AddrProducts       string        `mapstructure:"addr_products"`
	AddrAuth           string        `mapstructure:"addr_auth"`
	Timeout            time.Duration `mapstructure:"timeout"`
	Port               int           `mapstructure:"port"`
	Mode               string        `mapstructure:"mode"`
	CORSAllowedOrigins []string      `mapstructure:"cors_allowed_origins"`
	CookieSecure       bool          `mapstructure:"cookie_secure"`
	RedisPw            string        `mapstructure:"redis_pw"`
	RedisAddr          string        `mapstructure:"redis_addr"`
	RedisDictTTL       time.Duration `mapstructure:"redis_dictionaries_ttl"`
	DevAdminToken      string        `mapstructure:"dev_admin_token"`
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
	cfg.AddrAuth = platformconfig.StringFromEnv("AUTH_GRPC_ADDR", cfg.AddrAuth)
	cfg.RedisAddr = platformconfig.StringFromEnv("REDIS_ADDR", cfg.RedisAddr)
	cfg.RedisPw = platformconfig.StringFromEnv("REDIS_PASSWORD", "")
	cfg.DevAdminToken = platformconfig.StringFromEnv("DEV_ADMIN_TOKEN", cfg.DevAdminToken)
	cfg.Mode = platformconfig.StringFromEnv("APP_ENV", cfg.Mode)
	if origins := platformconfig.StringFromEnv("CORS_ALLOWED_ORIGINS", ""); origins != "" {
		cfg.CORSAllowedOrigins = splitOrigins(origins)
	}
	if secureValue := os.Getenv("COOKIE_SECURE"); secureValue != "" {
		cfg.CookieSecure, err = strconv.ParseBool(secureValue)
		if err != nil {
			return nil, format.Error(op, errors.New("COOKIE_SECURE must be true or false"))
		}
	}
	if cfg.Port, err = platformconfig.IntFromEnv("GATEWAY_PORT", cfg.Port); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.Timeout, err = platformconfig.DurationFromEnv("PRODUCT_GRPC_TIMEOUT", cfg.Timeout); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.RedisDictTTL, err = platformconfig.DurationFromEnv("REDIS_DICTIONARIES_TTL", defaultDuration(cfg.RedisDictTTL, 10*time.Minute)); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.AddrProducts == "" || cfg.AddrAuth == "" || cfg.Port == 0 || cfg.RedisAddr == "" {
		return nil, format.Error(op, errors.New("missing required gateway config"))
	}
	if len(cfg.CORSAllowedOrigins) == 0 {
		cfg.CORSAllowedOrigins = []string{"http://localhost:3000"}
	}
	for _, origin := range cfg.CORSAllowedOrigins {
		if origin == "*" {
			return nil, format.Error(op, errors.New("CORS_ALLOWED_ORIGINS must list exact origins when cookies are enabled"))
		}
	}
	if isDevelopment(cfg.Mode) && cfg.DevAdminToken == "" {
		return nil, format.Error(op, errors.New("DEV_ADMIN_TOKEN is required in development mode"))
	}
	return &cfg, nil
}

func splitOrigins(origins string) []string {
	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))
	for _, origin := range parts {
		if origin = strings.TrimSpace(origin); origin != "" {
			result = append(result, origin)
		}
	}
	return result
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
