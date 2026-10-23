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
	require.NoError(t, os.WriteFile(path, []byte("port: 8080\nmode: local\naddr_products: localhost:8008\naddr_auth: localhost:8009\nredis_addr: localhost:6379\ntimeout: 1s\n"), 0600))
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
func TestGatewayInvalidEnvAndTestDefaults(t *testing.T) {
	t.Setenv("TEST_PRODUCT_GRPC_ADDR", "products:9")
	require.Equal(t, "products:9", Test().AddrProducts)
	require.Equal(t, []string{"a", "b"}, splitOrigins(" a, b ,"))
}
