package psql

import (
	"context"
	"database/sql"

	"github.com/autumnterror/volha-backend/internal/auth-service/config"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	_ "github.com/lib/pq"
)

type Database struct {
	Driver *sql.DB
}

func Connect(cfg *config.Config) (*Database, error) { return connect(cfg, sql.Open) }

func connect(cfg *config.Config, openDB func(string, string) (*sql.DB, error)) (*Database, error) {
	const op = "auth.psql.Connect"
	db, err := openDB("postgres", cfg.ConnStr)
	if err != nil {
		return nil, format.Error(op, err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)
	ctx, cancel := context.WithTimeout(context.Background(), cfg.PostgresPingTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, format.Error(op, err)
	}
	log.Info(context.Background(), "postgres connection established")
	return &Database{Driver: db}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	if err := d.Driver.PingContext(ctx); err != nil {
		return format.Error("auth.psql.Ping", err)
	}
	return nil
}

func (d *Database) Close() error {
	if err := d.Driver.Close(); err != nil {
		return format.Error("auth.psql.Close", err)
	}
	log.Info(context.Background(), "postgres connection terminated")
	return nil
}
