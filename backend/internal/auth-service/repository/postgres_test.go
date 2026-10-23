package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func newMockStore(t *testing.T) (*Store, sqlmock.Sqlmock) {
	t.Helper()
	db, m, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, m.ExpectationsWereMet()); db.Close() })
	return New(db), m
}
func userRows() *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "email", "login"}).AddRow("user", "user@example.com", "login")
}
func TestRegistrationAndSessionTransactions(t *testing.T) {
	ctx := context.Background()
	expires := time.Now().Add(time.Hour)
	user := &domain.User{ID: "user", Email: "user@example.com", Login: "login"}
	t.Run("registration commits both records", func(t *testing.T) {
		s, m := newMockStore(t)
		m.ExpectBegin()
		m.ExpectExec("DELETE FROM auth_sessions").WillReturnResult(sqlmock.NewResult(0, 0))
		m.ExpectExec("INSERT INTO users").WithArgs(user.ID, user.Email, user.Login, "hash").WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectExec("INSERT INTO auth_sessions").WithArgs("sid", "user", "refresh", expires).WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectCommit()
		require.NoError(t, s.CreateUserAndSession(ctx, user, "hash", "sid", "refresh", expires))
	})
	t.Run("duplicate rolls back", func(t *testing.T) {
		s, m := newMockStore(t)
		m.ExpectBegin()
		m.ExpectExec("DELETE FROM auth_sessions").WillReturnResult(sqlmock.NewResult(0, 0))
		m.ExpectExec("INSERT INTO users").WillReturnError(&pq.Error{Code: "23505"})
		m.ExpectRollback()
		require.ErrorIs(t, s.CreateUserAndSession(ctx, user, "hash", "sid", "refresh", expires), domain.ErrAlreadyExists)
	})
	t.Run("new session and failure", func(t *testing.T) {
		s, m := newMockStore(t)
		m.ExpectBegin()
		m.ExpectExec("DELETE FROM auth_sessions").WillReturnResult(sqlmock.NewResult(0, 0))
		m.ExpectExec("INSERT INTO auth_sessions").WithArgs("sid", "user", "refresh", expires).WillReturnResult(sqlmock.NewResult(0, 1))
		m.ExpectCommit()
		require.NoError(t, s.CreateSession(ctx, "user", "sid", "refresh", expires))
		boom := errors.New("begin failed")
		m.ExpectBegin().WillReturnError(boom)
		require.ErrorIs(t, s.CreateSession(ctx, "user", "sid", "refresh", expires), boom)
	})
}
func TestSessionRotationLocksAndChecksOldToken(t *testing.T) {
	ctx := context.Background()
	expires := time.Now().Add(time.Hour)
	for _, old := range []string{"old", "wrong"} {
		t.Run(old, func(t *testing.T) {
			s, m := newMockStore(t)
			m.ExpectBegin()
			m.ExpectQuery("FOR UPDATE OF s").WithArgs("sid", "user").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "login", "hash"}).AddRow("user", "user@example.com", "login", "old"))
			if old == "wrong" {
				m.ExpectRollback()
				_, err := s.RotateSession(ctx, "user", "sid", old, "new", expires)
				require.ErrorIs(t, err, domain.ErrUnauthenticated)
				return
			}
			m.ExpectExec("UPDATE auth_sessions").WithArgs("sid", "new", expires).WillReturnResult(sqlmock.NewResult(0, 1))
			m.ExpectCommit()
			user, err := s.RotateSession(ctx, "user", "sid", old, "new", expires)
			require.NoError(t, err)
			require.Equal(t, "user", user.ID)
		})
	}
	s, m := newMockStore(t)
	m.ExpectQuery("JOIN users").WithArgs("sid", "user").WillReturnRows(userRows())
	u, err := s.AuthenticateSession(ctx, "user", "sid")
	require.NoError(t, err)
	require.Equal(t, "login", u.Login)
	m.ExpectQuery("JOIN users").WillReturnError(sql.ErrNoRows)
	_, err = s.AuthenticateSession(ctx, "user", "gone")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	m.ExpectExec("DELETE FROM auth_sessions").WithArgs("user", "sid", "hash").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, s.RevokeSession(ctx, "user", "sid", "hash"))
}
func TestUserLookupAndChanges(t *testing.T) {
	ctx := context.Background()
	s, m := newMockStore(t)
	m.ExpectQuery("lower\\(email\\)").WithArgs("login").WillReturnRows(sqlmock.NewRows([]string{"id", "email", "login", "hash"}).AddRow("user", "u@e.com", "login", "hash"))
	u, hash, err := s.FindUserForLogin(ctx, "login")
	require.NoError(t, err)
	require.Equal(t, "user", u.ID)
	require.Equal(t, "hash", hash)
	m.ExpectQuery("lower\\(email\\)").WillReturnError(sql.ErrNoRows)
	_, _, err = s.FindUserForLogin(ctx, "missing")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	m.ExpectQuery("SELECT password_hash").WithArgs("user").WillReturnRows(sqlmock.NewRows([]string{"hash"}).AddRow("hash"))
	hash, err = s.PasswordHash(ctx, "user")
	require.NoError(t, err)
	require.Equal(t, "hash", hash)
	m.ExpectBegin()
	m.ExpectQuery("SELECT id FROM auth_sessions").WithArgs("user", "sid").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("sid"))
	m.ExpectExec("UPDATE users SET password_hash").WithArgs("user", "new").WillReturnResult(sqlmock.NewResult(0, 1))
	m.ExpectExec("DELETE FROM auth_sessions").WithArgs("user", "sid").WillReturnResult(sqlmock.NewResult(0, 2))
	m.ExpectCommit()
	require.NoError(t, s.ChangePassword(ctx, "user", "sid", "new"))
	m.ExpectBegin()
	m.ExpectQuery("SELECT id FROM auth_sessions").WithArgs("user", "sid").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("sid"))
	m.ExpectQuery("UPDATE users SET login").WithArgs("user", "new").WillReturnRows(userRows())
	m.ExpectCommit()
	u, err = s.ChangeLogin(ctx, "user", "sid", "new")
	require.NoError(t, err)
	require.Equal(t, "user", u.ID)
	m.ExpectBegin()
	m.ExpectQuery("SELECT id FROM auth_sessions").WillReturnError(sql.ErrNoRows)
	m.ExpectRollback()
	_, err = s.ChangeLogin(ctx, "user", "gone", "new")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
}
func TestCartProjectionAndMutations(t *testing.T) {
	ctx := context.Background()
	s, m := newMockStore(t)
	m.ExpectQuery("FROM cart_items ci").WithArgs("user").WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "quantity", "price", "sale", "image", "total"}).AddRow("p1", "desk", "Desk", 2, 10, 0, "a.jpg", 20).AddRow("p2", "chair", "Chair", 1, 7, 5, "b.jpg", 5))
	cart, err := s.GetCart(ctx, "user")
	require.NoError(t, err)
	require.Len(t, cart.Items, 2)
	require.Equal(t, int32(3), cart.TotalQuantity)
	require.Equal(t, int64(25), cart.TotalPrice)
	m.ExpectBegin()
	m.ExpectQuery("SELECT EXISTS").WithArgs("p1").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	m.ExpectQuery("INSERT INTO cart_items").WithArgs("user", "p1", int32(2)).WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(2))
	m.ExpectCommit()
	require.NoError(t, s.AddCartItem(ctx, "user", "p1", 2))
	m.ExpectBegin()
	m.ExpectQuery("SELECT EXISTS").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	m.ExpectQuery("INSERT INTO cart_items").WillReturnError(sql.ErrNoRows)
	m.ExpectRollback()
	require.ErrorIs(t, s.AddCartItem(ctx, "user", "p1", 99), domain.ErrCartLimit)
	m.ExpectBegin()
	m.ExpectQuery("SELECT EXISTS").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	m.ExpectRollback()
	require.ErrorIs(t, s.AddCartItem(ctx, "user", "missing", 1), domain.ErrNotFound)
	m.ExpectExec("UPDATE cart_items").WithArgs("user", "p1", int32(3)).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, s.UpdateCartItem(ctx, "user", "p1", 3))
	m.ExpectExec("UPDATE cart_items").WillReturnResult(sqlmock.NewResult(0, 0))
	require.ErrorIs(t, s.UpdateCartItem(ctx, "user", "missing", 3), domain.ErrNotFound)
	m.ExpectExec("DELETE FROM cart_items").WithArgs("user", "p1").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, s.RemoveCartItem(ctx, "user", "p1"))
	m.ExpectExec("DELETE FROM cart_items").WithArgs("user", "missing").WillReturnResult(sqlmock.NewResult(0, 0))
	require.ErrorIs(t, s.RemoveCartItem(ctx, "user", "missing"), domain.ErrNotFound)
	m.ExpectExec("DELETE FROM cart_items").WithArgs("user").WillReturnResult(sqlmock.NewResult(0, 3))
	require.NoError(t, s.ClearCart(ctx, "user"))
}
func TestConstraintErrors(t *testing.T) {
	for code, want := range map[string]error{"23505": domain.ErrAlreadyExists, "23503": domain.ErrNotFound, "23514": domain.ErrInvalidInput, "22001": domain.ErrInvalidInput} {
		require.ErrorIs(t, mapPostgresError(&pq.Error{Code: pq.ErrorCode(code)}), want)
	}
	require.NoError(t, mapPostgresError(nil))
	boom := errors.New("io")
	require.ErrorIs(t, mapPostgresError(boom), boom)
}
