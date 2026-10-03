package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	platformpostgres "github.com/autumnterror/volha-backend/internal/platform/postgres"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	typeMigration := flag.String("type", "up", "type of migration action")
	flag.Parse()

	switch *typeMigration {
	case "up":
		if err := upMigrate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "down":
		if err := downMigrate(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		_, _ = fmt.Fprintln(os.Stderr, "type not recognized")
		os.Exit(1)
	}
}

func upMigrate() error {
	const op = "migrator-main.upMigrate"

	cfg, err := migrationConfigFromEnv()
	if err != nil {
		return format.Error(op, err)
	}

	db, err := sql.Open("postgres", cfg.connStr())
	if err != nil {
		return format.Error(op, err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return format.Error(op, err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsSourceURL(),
		"postgres", driver)
	if err != nil {
		return format.Error(op, err)
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return format.Error(op, err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		_, _ = fmt.Fprintln(os.Stdout, "No migrations to apply")
		return nil
	}
	_, _ = fmt.Fprintln(os.Stdout, "Migrations applied successfully!")
	return nil
}

func downMigrate() error {
	const op = "migrator-main.downMigrate"

	cfg, err := migrationConfigFromEnv()
	if err != nil {
		return format.Error(op, err)
	}

	db, err := sql.Open("postgres", cfg.connStr())
	if err != nil {
		return format.Error(op, err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return format.Error(op, err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		migrationsSourceURL(),
		"postgres", driver)
	if err != nil {
		return format.Error(op, err)
	}
	defer m.Close()

	err = m.Steps(-1)
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return format.Error(op, err)
	}
	if errors.Is(err, migrate.ErrNoChange) {
		_, _ = fmt.Fprintln(os.Stdout, "No migrations to rollback")
		return nil
	}

	_, _ = fmt.Fprintln(os.Stdout, "Migration rolled back successfully!")
	return nil
}

type migrationConfig struct {
	host     string
	port     int
	user     string
	password string
	db       string
	sslMode  string
}

func migrationConfigFromEnv() (migrationConfig, error) {
	cfg := migrationConfig{
		host:     os.Getenv("POSTGRES_HOST"),
		user:     os.Getenv("POSTGRES_USER"),
		password: os.Getenv("POSTGRES_PASSWORD"),
		db:       os.Getenv("POSTGRES_DB"),
		sslMode:  os.Getenv("POSTGRES_SSLMODE"),
		port:     5432,
	}
	if cfg.sslMode == "" {
		cfg.sslMode = "disable"
	}
	if port := os.Getenv("POSTGRES_PORT"); port != "" {
		parsed, err := strconv.Atoi(port)
		if err != nil {
			return cfg, fmt.Errorf("invalid POSTGRES_PORT: %w", err)
		}
		cfg.port = parsed
	}
	if cfg.host == "" || cfg.user == "" || cfg.password == "" || cfg.db == "" {
		return cfg, errors.New("missing required postgres migration environment")
	}
	return cfg, nil
}

func (c migrationConfig) connStr() string {
	return platformpostgres.URL(c.host, c.port, c.db, c.user, c.password, c.sslMode)
}

func migrationsSourceURL() string {
	path := os.Getenv("MIGRATIONS_PATH")
	if path == "" {
		path = filepath.Join("build", "volha", "migrations")
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "file://" + filepath.ToSlash(path)
	}
	return "file://" + filepath.ToSlash(abs)
}
