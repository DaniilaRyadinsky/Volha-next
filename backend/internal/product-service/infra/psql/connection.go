package psql

import (
	"context"
	"database/sql"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/product-service/config"
	_ "github.com/lib/pq"
)

type PostgresDb struct {
	Driver *sql.DB
}

// NewConnect is constructor of PostgresDb. Construct with connection
func NewConnect(cfg *config.Config) (*PostgresDb, error) {
	const op = "psql.NewConnect"

	db, err := sql.Open("postgres", cfg.ConnStr)
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
	return &PostgresDb{Driver: db}, nil
}

// Disconnect is terminated connection
func (d *PostgresDb) Disconnect() error {
	const op = "psql.PostgresDb.Disconnect"
	err := d.Driver.Close()
	if err != nil {
		return format.Error(op, err)
	}

	log.Info(context.Background(), "postgres connection terminated")
	return nil
}

func (d *PostgresDb) Ping(ctx context.Context) error {
	const op = "psql.PostgresDb.Ping"
	if err := d.Driver.PingContext(ctx); err != nil {
		return format.Error(op, err)
	}
	return nil
}
