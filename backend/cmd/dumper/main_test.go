package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBackupPublishesOnlySuccessfulDump(t *testing.T) {
	dir := t.TempDir()
	cfg := dumperConfig{host: "db", user: "user", password: "secret", db: "volha", port: 5432}
	err := dumpWithRunner(context.Background(), cfg, dir, func(cmd *exec.Cmd) error {
		require.Contains(t, cmd.Args, "pg_dump")
		require.Contains(t, cmd.Env, "PGPASSWORD=secret")
		_, err := io.WriteString(cmd.Stdout, "SQL data")
		return err
	})
	require.NoError(t, err)
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	require.NoError(t, err)
	require.Len(t, files, 1)
	data, err := os.ReadFile(files[0])
	require.NoError(t, err)
	require.Equal(t, "SQL data", string(data))
	failedDir := t.TempDir()
	require.Error(t, dumpWithRunner(context.Background(), cfg, failedDir, func(cmd *exec.Cmd) error {
		_, _ = io.WriteString(cmd.Stdout, "partial")
		return errors.New("dump failed")
	}))
	entries, err := os.ReadDir(failedDir)
	require.NoError(t, err)
	require.Empty(t, entries)
}
func TestBackupRetentionAndConfiguration(t *testing.T) {
	dir := t.TempDir()
	for i := 0; i < 22; i++ {
		require.NoError(t, os.WriteFile(filepath.Join(dir, fmt.Sprintf("dump.20261022_%06d.sql", i)), []byte("sql"), 0600))
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "manual.sql"), []byte("manual"), 0600))
	cleanupDumps(context.Background(), dir)
	files, err := filepath.Glob(filepath.Join(dir, "dump.*.sql"))
	require.NoError(t, err)
	require.Len(t, files, 20)
	_, err = os.Stat(filepath.Join(dir, "manual.sql"))
	require.NoError(t, err)
	for key, value := range map[string]string{"POSTGRES_HOST": "db", "POSTGRES_USER": "u", "POSTGRES_PASSWORD": "secret", "POSTGRES_DB": "volha", "POSTGRES_PORT": "5432"} {
		t.Setenv(key, value)
	}
	cfg, err := loadConfig()
	require.NoError(t, err)
	require.Equal(t, 5432, cfg.port)
	t.Setenv("POSTGRES_PORT", "bad")
	_, err = loadConfig()
	require.Error(t, err)
}
