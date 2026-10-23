package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTypedEnvAndDotenv(t *testing.T) {
	t.Setenv("VALUE", "")
	require.Equal(t, "fallback", StringFromEnv("VALUE", "fallback"))
	n, err := IntFromEnv("VALUE", 42)
	require.NoError(t, err)
	require.Equal(t, 42, n)
	d, err := DurationFromEnv("VALUE", time.Second)
	require.NoError(t, err)
	require.Equal(t, time.Second, d)
	t.Setenv("VALUE", "123")
	n, err = IntFromEnv("VALUE", 0)
	require.NoError(t, err)
	require.Equal(t, 123, n)
	require.Equal(t, "123", StringFromEnv("VALUE", ""))
	t.Setenv("VALUE", "2s")
	d, err = DurationFromEnv("VALUE", 0)
	require.NoError(t, err)
	require.Equal(t, 2*time.Second, d)
	_, err = IntFromEnv("VALUE", 0)
	require.Error(t, err)
	t.Setenv("VALUE", "bad")
	_, err = DurationFromEnv("VALUE", 0)
	require.Error(t, err)
	root := t.TempDir()
	require.NoError(t, os.Mkdir(filepath.Join(root, "backend"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(root, ".env"), []byte("VOLHA_TEST_DOTENV=from-parent"), 0600))
	t.Chdir(filepath.Join(root, "backend"))
	t.Setenv("VOLHA_TEST_DOTENV", "")
	require.NoError(t, os.Unsetenv("VOLHA_TEST_DOTENV"))
	LoadDotEnv()
	require.Equal(t, "from-parent", os.Getenv("VOLHA_TEST_DOTENV"))
}
