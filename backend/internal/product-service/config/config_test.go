package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetupConfigAndEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "service.yaml")
	require.NoError(t, os.WriteFile(path, []byte("port: 8008\nmode: local\ndata_source: localhost\nport_postgres: 5432\n"), 0600))
	t.Setenv("CONFIG_PATH", path)
	t.Setenv("APP_ENV", "local")
	t.Setenv("POSTGRES_HOST", "database")
	t.Setenv("POSTGRES_USER", "user")
	t.Setenv("POSTGRES_PASSWORD", "password")
	t.Setenv("POSTGRES_DB", "volha")
	t.Setenv("JWT_SECRET", strings.Repeat("x", 32))
	t.Setenv("REDIS_PASSWORD", "redispass")
	t.Setenv("DEV_ADMIN_TOKEN", "admin")
	cfg, err := Setup()
	require.NoError(t, err)
	require.NotNil(t, cfg)
	t.Setenv("CONFIG_PATH", path+".missing")
	_, err = Setup()
	require.Error(t, err)
}
func TestConfigDefaults(t *testing.T) {
	require.Equal(t, 7, defaultInt(0, 7))
	require.Equal(t, 8, defaultInt(8, 7))
	require.Equal(t, defaultDuration(0, 1), defaultDuration(1, 2))
}
func TestIntegrationConfigUsesIsolatedVariables(t *testing.T) {
	t.Setenv("TEST_POSTGRES_HOST", "isolated")
	t.Setenv("TEST_POSTGRES_PORT", "5544")
	t.Setenv("TEST_POSTGRES_USER", "test")
	t.Setenv("TEST_POSTGRES_PASSWORD", "secret")
	require.Contains(t, Test().ConnStr, "isolated:5544")
}
