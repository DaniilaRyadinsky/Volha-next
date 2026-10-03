package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"syscall"
	"time"

	"github.com/autumnterror/volha-backend/internal/platform/buildinfo"
	platformconfig "github.com/autumnterror/volha-backend/internal/platform/config"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
)

const (
	dumpsDir = "/dumps"
	maxDumps = 20
	dumpTime = time.Hour
)

var dumpRegexp = regexp.MustCompile(`^dump\.\d{8}_\d{6}\.sql$`)

type dumperConfig struct {
	user     string
	password string
	host     string
	db       string
	port     int
}

func main() {
	platformconfig.LoadDotEnv()
	log.Configure(log.Config{
		Service: "volha-dumper",
		Env:     envOr("APP_ENV", "local"),
		Version: buildinfo.Version,
		Level:   os.Getenv("LOG_LEVEL"),
		Format:  os.Getenv("LOG_FORMAT"),
	})

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log.Info(ctx, "service starting",
		"version", buildinfo.Version,
		"commit", buildinfo.Commit,
		"build_time", buildinfo.BuildTime,
	)

	cfg, err := loadConfig()
	if err != nil {
		log.Error(ctx, "dumper config failed", "operation", "dumper.config", "error", err)
		os.Exit(1)
	}
	if err := os.MkdirAll(dumpsDir, 0755); err != nil {
		log.Error(ctx, "create dumps directory failed", "operation", "dumper.setup", "error", err)
		os.Exit(1)
	}

	if err := dump(ctx, cfg, dumpsDir); err != nil {
		log.Error(ctx, "initial dump failed", "operation", "dumper.dump.initial", "error", err)
	}
	log.Info(ctx, "next dump scheduled", "next_at", time.Now().UTC().Add(dumpTime).Format(time.RFC3339))

	ticker := time.NewTicker(dumpTime)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			log.Info(context.Background(), "service stopped", "operation", "dumper.shutdown")
			return
		case <-ticker.C:
			if err := dump(ctx, cfg, dumpsDir); err != nil {
				log.Error(ctx, "scheduled dump failed", "operation", "dumper.dump.scheduled", "error", err)
			}
			log.Info(ctx, "next dump scheduled", "next_at", time.Now().UTC().Add(dumpTime).Format(time.RFC3339))
		}
	}
}

func loadConfig() (dumperConfig, error) {
	port, err := strconv.Atoi(envOr("POSTGRES_PORT", "5432"))
	if err != nil {
		return dumperConfig{}, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
	}
	cfg := dumperConfig{
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		host:     os.Getenv("POSTGRES_HOST"),
		db:       os.Getenv("POSTGRES_DB"),
		port:     port,
	}
	if cfg.user == "" || cfg.password == "" || cfg.host == "" || cfg.db == "" {
		return cfg, errors.New("missing required postgres environment")
	}
	return cfg, nil
}

func dump(ctx context.Context, cfg dumperConfig, dir string) error {
	start := time.Now()
	ts := start.UTC().Format("20060102_150405")
	filename := filepath.Join(dir, fmt.Sprintf("dump.%s.sql", ts))
	tmpFilename := filename + ".tmp"

	f, err := os.Create(tmpFilename)
	if err != nil {
		return fmt.Errorf("create temp dump file: %w", err)
	}
	defer func() {
		_ = f.Close()
	}()

	cmd := exec.CommandContext(
		ctx,
		"pg_dump",
		"-h", cfg.host,
		"-p", strconv.Itoa(cfg.port),
		"-U", cfg.user,
		cfg.db,
	)
	var stderr bytes.Buffer
	cmd.Stdout = f
	cmd.Stderr = &stderr
	cmd.Env = append(os.Environ(), "PGPASSWORD="+cfg.password)

	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpFilename)
		return fmt.Errorf("pg_dump failed: %w: %s", err, stderr.String())
	}
	if err := f.Close(); err != nil {
		_ = os.Remove(tmpFilename)
		return fmt.Errorf("close temp dump file: %w", err)
	}
	if err := os.Rename(tmpFilename, filename); err != nil {
		_ = os.Remove(tmpFilename)
		return fmt.Errorf("finalize dump file: %w", err)
	}

	log.Info(ctx, "dump completed",
		"file", filepath.Base(filename),
		"duration_ms", time.Since(start).Milliseconds(),
	)
	cleanupDumps(ctx, dir)
	return nil
}

func cleanupDumps(ctx context.Context, dir string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Warn(ctx, "read dumps directory failed", "operation", "dumper.cleanup", "error", err)
		return
	}
	var dumpFiles []string
	for _, f := range files {
		if !f.IsDir() && dumpRegexp.MatchString(f.Name()) {
			dumpFiles = append(dumpFiles, f.Name())
		}
	}
	if len(dumpFiles) <= maxDumps {
		return
	}

	sort.Strings(dumpFiles)
	toDelete := len(dumpFiles) - maxDumps
	removed := 0
	for i := 0; i < toDelete; i++ {
		path := filepath.Join(dir, dumpFiles[i])
		if err := os.Remove(path); err != nil {
			log.Warn(ctx, "remove old dump failed", "operation", "dumper.cleanup", "file", dumpFiles[i], "error", err)
			continue
		}
		removed++
	}
	log.Info(ctx, "old dumps removed", "operation", "dumper.cleanup", "removed_count", removed)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
