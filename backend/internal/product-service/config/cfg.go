package config

import (
	"errors"
	"os"
	"time"

	platformconfig "github.com/autumnterror/volha-backend/internal/platform/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	platformpostgres "github.com/autumnterror/volha-backend/internal/platform/postgres"
	"github.com/spf13/viper"
)

type Config struct {
	ConnStr             string        `mapstructure:"conn_str"`
	Port                int           `mapstructure:"port"`
	Mode                string        `mapstructure:"mode"`
	MaxOpenConns        int           `mapstructure:"max_open_conns"`
	MaxIdleConns        int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime     time.Duration `mapstructure:"conn_max_idle_time"`
	PostgresPingTimeout time.Duration `mapstructure:"postgres_ping_timeout"`
}

// Setup creates config structure.
func Setup() (*Config, error) {
	const op = "config.setup"

	type _сonfig struct {
		DataSource          string        `mapstructure:"data_source"`
		PortPostgres        int           `mapstructure:"port_postgres"`
		Port                int           `mapstructure:"port"`
		Mode                string        `mapstructure:"mode"`
		MaxOpenConns        int           `mapstructure:"max_open_conns"`
		MaxIdleConns        int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
		ConnMaxIdleTime     time.Duration `mapstructure:"conn_max_idle_time"`
		PostgresPingTimeout time.Duration `mapstructure:"postgres_ping_timeout"`
	}
	platformconfig.LoadDotEnv()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/product-service.yaml"
	}

	v := viper.New()
	v.SetConfigFile(configPath)
	var cfg _сonfig
	if err := v.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, format.Error(op, err)
	}
	var err error
	cfg.DataSource = platformconfig.StringFromEnv("POSTGRES_HOST", cfg.DataSource)
	if cfg.PortPostgres, err = platformconfig.IntFromEnv("POSTGRES_PORT", cfg.PortPostgres); err != nil {
		return nil, format.Error(op, err)
	}
	if cfg.Port, err = platformconfig.IntFromEnv("PRODUCT_GRPC_PORT", cfg.Port); err != nil {
		return nil, format.Error(op, err)
	}
	cfg.Mode = platformconfig.StringFromEnv("APP_ENV", cfg.Mode)

	user := os.Getenv("POSTGRES_USER")
	pw := os.Getenv("POSTGRES_PASSWORD")
	db := os.Getenv("POSTGRES_DB")

	if user == "" || pw == "" || db == "" {
		return nil, format.Error(op, errors.New("missing environment variables"))
	}
	if cfg.DataSource == "" || cfg.PortPostgres == 0 || cfg.Port == 0 {
		return nil, format.Error(op, errors.New("missing required product-service config"))
	}

	return &Config{
		ConnStr:             platformpostgres.URL(cfg.DataSource, cfg.PortPostgres, db, user, pw, "disable"),
		Port:                cfg.Port,
		Mode:                cfg.Mode,
		MaxOpenConns:        envInt("POSTGRES_MAX_OPEN_CONNS", defaultInt(cfg.MaxOpenConns, 10)),
		MaxIdleConns:        envInt("POSTGRES_MAX_IDLE_CONNS", defaultInt(cfg.MaxIdleConns, 5)),
		ConnMaxLifetime:     envDuration("POSTGRES_CONN_MAX_LIFETIME", defaultDuration(cfg.ConnMaxLifetime, 30*time.Minute)),
		ConnMaxIdleTime:     envDuration("POSTGRES_CONN_MAX_IDLE_TIME", defaultDuration(cfg.ConnMaxIdleTime, 5*time.Minute)),
		PostgresPingTimeout: envDuration("POSTGRES_PING_TIMEOUT", defaultDuration(cfg.PostgresPingTimeout, 5*time.Second)),
	}, nil
}

func defaultInt(value int, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultDuration(value time.Duration, fallback time.Duration) time.Duration {
	if value == 0 {
		return fallback
	}
	return value
}

func envInt(name string, fallback int) int {
	value, err := platformconfig.IntFromEnv(name, fallback)
	if err != nil {
		return fallback
	}
	return value
}

func envDuration(name string, fallback time.Duration) time.Duration {
	value, err := platformconfig.DurationFromEnv(name, fallback)
	if err != nil {
		return fallback
	}
	return value
}
