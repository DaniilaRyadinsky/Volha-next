package psql

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/autumnterror/volha-backend/internal/auth-service/config"
	"github.com/stretchr/testify/require"
)

func TestConnectionLifecycle(t *testing.T) {
	db, m, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	m.ExpectPing()
	cfg := &config.Config{ConnStr: "dsn", PostgresPingTimeout: time.Second, MaxOpenConns: 5, MaxIdleConns: 2}
	conn, err := connect(cfg, func(driver, dsn string) (*sql.DB, error) {
		require.Equal(t, "postgres", driver)
		require.Equal(t, "dsn", dsn)
		return db, nil
	})
	require.NoError(t, err)
	require.Equal(t, 5, db.Stats().MaxOpenConnections)
	m.ExpectPing()
	require.NoError(t, conn.Ping(context.Background()))
	m.ExpectClose()
	require.NoError(t, conn.Close())
	require.NoError(t, m.ExpectationsWereMet())
}
func TestConnectionFailures(t *testing.T) {
	boom := errors.New("open failure")
	_, err := connect(&config.Config{}, func(string, string) (*sql.DB, error) { return nil, boom })
	require.ErrorIs(t, err, boom)
	db, m, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	require.NoError(t, err)
	m.ExpectPing().WillReturnError(boom)
	m.ExpectClose()
	_, err = connect(&config.Config{PostgresPingTimeout: time.Second, MaxIdleConns: 2}, func(string, string) (*sql.DB, error) { return db, nil })
	require.ErrorIs(t, err, boom)
	require.NoError(t, m.ExpectationsWereMet())
	_, err = Connect(&config.Config{ConnStr: "not a postgres connection string"})
	require.Error(t, err)
}
