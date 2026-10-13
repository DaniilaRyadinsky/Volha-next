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
	ConnStr             string
	Port                int
	Mode                string
	MaxOpenConns        int
	MaxIdleConns        int
	ConnMaxLifetime     time.Duration
	ConnMaxIdleTime     time.Duration
	PostgresPingTimeout time.Duration
	JWTSecret           string
	AccessTokenTTL      time.Duration
	RefreshTokenTTL     time.Duration
}

func Setup() (*Config, error) {
	const op = "auth.config.Setup"
	type fileConfig struct {
		DataSource          string        `mapstructure:"data_source"`
		PortPostgres        int           `mapstructure:"port_postgres"`
		Port                int           `mapstructure:"port"`
		Mode                string        `mapstructure:"mode"`
		MaxOpenConns        int           `mapstructure:"max_open_conns"`
		MaxIdleConns        int           `mapstructure:"max_idle_conns"`
		ConnMaxLifetime     time.Duration `mapstructure:"conn_max_lifetime"`
		ConnMaxIdleTime     time.Duration `mapstructure:"conn_max_idle_time"`
		PostgresPingTimeout time.Duration `mapstructure:"postgres_ping_timeout"`
		AccessTokenTTL      time.Duration `mapstructure:"access_token_ttl"`
		RefreshTokenTTL     time.Duration `mapstructure:"refresh_token_ttl"`
	}

	platformconfig.LoadDotEnv()
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/auth-service.yaml"
	}
	v := viper.New()
	v.SetConfigFile(configPath)
	var fileCfg fileConfig
	if err := v.ReadInConfig(); err != nil {
		return nil, format.Error(op, err)
	}
	if err := v.Unmarshal(&fileCfg); err != nil {
		return nil, format.Error(op, err)
	}

	var err error
	fileCfg.DataSource = platformconfig.StringFromEnv("POSTGRES_HOST", fileCfg.DataSource)
	if fileCfg.PortPostgres, err = platformconfig.IntFromEnv("POSTGRES_PORT", fileCfg.PortPostgres); err != nil {
		return nil, format.Error(op, err)
	}
	if fileCfg.Port, err = platformconfig.IntFromEnv("AUTH_GRPC_PORT", fileCfg.Port); err != nil {
		return nil, format.Error(op, err)
	}
	fileCfg.Mode = platformconfig.StringFromEnv("APP_ENV", fileCfg.Mode)
	accessTTL, err := platformconfig.DurationFromEnv("ACCESS_TOKEN_TTL", defaultDuration(fileCfg.AccessTokenTTL, 15*time.Minute))
	if err != nil {
		return nil, format.Error(op, err)
	}
	refreshTTL, err := platformconfig.DurationFromEnv("REFRESH_TOKEN_TTL", defaultDuration(fileCfg.RefreshTokenTTL, 30*24*time.Hour))
	if err != nil {
		return nil, format.Error(op, err)
	}
	secret := platformconfig.StringFromEnv("JWT_SECRET", "")
	if os.Getenv("POSTGRES_USER") == "" || os.Getenv("POSTGRES_PASSWORD") == "" || os.Getenv("POSTGRES_DB") == "" {
		return nil, format.Error(op, errors.New("missing postgres credentials"))
	}
	if fileCfg.DataSource == "" || fileCfg.PortPostgres <= 0 || fileCfg.Port <= 0 {
		return nil, format.Error(op, errors.New("missing required auth-service config"))
	}
	if len(secret) < 32 {
		return nil, format.Error(op, errors.New("JWT_SECRET must contain at least 32 bytes"))
	}
	if accessTTL < time.Minute || accessTTL > 30*time.Minute || refreshTTL < time.Hour || refreshTTL > 30*24*time.Hour || accessTTL >= refreshTTL {
		return nil, format.Error(op, errors.New("token lifetimes are outside allowed ranges"))
	}

	return &Config{
		ConnStr:             platformpostgres.URL(fileCfg.DataSource, fileCfg.PortPostgres, os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), "disable"),
		Port:                fileCfg.Port,
		Mode:                fileCfg.Mode,
		MaxOpenConns:        defaultInt(fileCfg.MaxOpenConns, 10),
		MaxIdleConns:        defaultInt(fileCfg.MaxIdleConns, 5),
		ConnMaxLifetime:     defaultDuration(fileCfg.ConnMaxLifetime, 30*time.Minute),
		ConnMaxIdleTime:     defaultDuration(fileCfg.ConnMaxIdleTime, 5*time.Minute),
		PostgresPingTimeout: defaultDuration(fileCfg.PostgresPingTimeout, 5*time.Second),
		JWTSecret:           secret,
		AccessTokenTTL:      accessTTL,
		RefreshTokenTTL:     refreshTTL,
	}, nil
}

func defaultInt(value, fallback int) int {
	if value == 0 {
		return fallback
	}
	return value
}

func defaultDuration(value, fallback time.Duration) time.Duration {
	if value == 0 {
		return fallback
	}
	return value
}
