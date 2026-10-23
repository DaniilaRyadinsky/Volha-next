package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMigrationConfigurationAndSource(t *testing.T) {
	for key, value := range map[string]string{"POSTGRES_HOST": "db", "POSTGRES_USER": "user", "POSTGRES_PASSWORD": "secret", "POSTGRES_DB": "volha", "POSTGRES_PORT": "5432", "POSTGRES_SSLMODE": ""} {
		t.Setenv(key, value)
	}
	cfg, err := migrationConfigFromEnv()
	require.NoError(t, err)
	require.Equal(t, "disable", cfg.sslMode)
	require.Contains(t, cfg.connStr(), "db:5432")
	t.Setenv("MIGRATIONS_PATH", "")
	require.True(t, strings.HasSuffix(migrationsSourceURL(), "/migrations"))
	t.Setenv("MIGRATIONS_PATH", t.TempDir())
	require.Contains(t, migrationsSourceURL(), "file://")
	t.Setenv("POSTGRES_PORT", "bad")
	_, err = migrationConfigFromEnv()
	require.Error(t, err)
	t.Setenv("POSTGRES_PORT", "5432")
	t.Setenv("POSTGRES_PASSWORD", "")
	_, err = migrationConfigFromEnv()
	require.Error(t, err)
}
