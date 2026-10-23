package psqltx

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
	"github.com/stretchr/testify/require"
)

func TestTransactionCommitRollbackAndBeginFailure(t *testing.T) {
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "commit", true: "rollback"}[failure], func(t *testing.T) {
			db, m, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()
			m.ExpectBegin()
			boom := errors.New("operation failed")
			if failure {
				m.ExpectRollback()
			} else {
				m.ExpectCommit()
			}
			r := NewTxRunner(db)
			err = r.RunInTx(context.Background(), func(ctx context.Context) error {
				tx, ok := TxFromContext(ctx)
				require.True(t, ok)
				require.NotNil(t, tx)
				if failure {
					return boom
				}
				return nil
			})
			if failure {
				require.ErrorIs(t, err, boom)
			} else {
				require.NoError(t, err)
			}
			require.NoError(t, m.ExpectationsWereMet())
		})
	}
	db, m, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	boom := errors.New("begin failed")
	m.ExpectBegin().WillReturnError(boom)
	require.ErrorIs(t, NewTxRunner(db).RunInTx(context.Background(), func(context.Context) error { t.Fatal("must not run"); return nil }), boom)
	require.NoError(t, m.ExpectationsWereMet())
}
func TestProviderUsesContextTransactionForEveryRepository(t *testing.T) {
	db, m, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()
	provider := NewRepoProvider(db)
	check := func(ctx context.Context, want repository.SqlRepo) {
		for _, repo := range []any{provider.Product(ctx), provider.Catalog(ctx), provider.Dict(ctx), provider.Pcp(ctx), provider.FileRefs(ctx), provider.Article(ctx)} {
			require.Same(t, want, repo.(repository.Driver).Driver)
		}
	}
	check(context.Background(), db)
	m.ExpectBegin()
	tx, err := db.Begin()
	require.NoError(t, err)
	check(context.WithValue(context.Background(), txKey{}, tx), tx)
	m.ExpectRollback()
	require.NoError(t, tx.Rollback())
	require.NoError(t, m.ExpectationsWereMet())
}
