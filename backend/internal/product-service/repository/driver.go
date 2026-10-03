package repository

import (
	"context"
	"database/sql"
)

type SqlRepo interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type Driver struct {
	Driver SqlRepo
}
