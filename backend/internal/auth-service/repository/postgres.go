package repository

import (
	"context"
	"crypto/subtle"
	"database/sql"
	"errors"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/lib/pq"
)

type Store struct {
	db *sql.DB
}

func New(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateUserAndSession(ctx context.Context, user *domain.User, passwordHash, sessionID, refreshHash string, expiresAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM auth_sessions WHERE expires_at <= now()`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx,
		`INSERT INTO users (id, email, login, password_hash) VALUES ($1, $2, $3, $4)`,
		user.ID, user.Email, user.Login, passwordHash,
	); err != nil {
		return mapPostgresError(err)
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO auth_sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3, $4)`,
		sessionID, user.ID, refreshHash, expiresAt,
	); err != nil {
		return mapPostgresError(err)
	}
	return tx.Commit()
}

func (s *Store) FindUserForLogin(ctx context.Context, identifier string) (*domain.User, string, error) {
	var user domain.User
	var passwordHash string
	err := s.db.QueryRowContext(ctx,
		`SELECT id, email, login, password_hash FROM users WHERE lower(email) = lower($1) OR lower(login) = lower($1)`,
		identifier,
	).Scan(&user.ID, &user.Email, &user.Login, &passwordHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, "", domain.ErrUnauthenticated
	}
	if err != nil {
		return nil, "", err
	}
	return &user, passwordHash, nil
}

func (s *Store) CreateSession(ctx context.Context, userID, sessionID, refreshHash string, expiresAt time.Time) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `DELETE FROM auth_sessions WHERE expires_at <= now()`); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx,
		`INSERT INTO auth_sessions (id, user_id, refresh_token_hash, expires_at) VALUES ($1, $2, $3, $4)`,
		sessionID, userID, refreshHash, expiresAt,
	); err != nil {
		return mapPostgresError(err)
	}
	return tx.Commit()
}

func (s *Store) RotateSession(ctx context.Context, userID, sessionID, oldRefreshHash, newRefreshHash string, expiresAt time.Time) (*domain.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var user domain.User
	var storedHash string
	err = tx.QueryRowContext(ctx, `
SELECT u.id, u.email, u.login, s.refresh_token_hash
FROM auth_sessions s
JOIN users u ON u.id = s.user_id
WHERE s.id = $1 AND s.user_id = $2 AND s.expires_at > now()
FOR UPDATE OF s`, sessionID, userID).Scan(&user.ID, &user.Email, &user.Login, &storedHash)
	if errors.Is(err, sql.ErrNoRows) || (err == nil && subtle.ConstantTimeCompare([]byte(storedHash), []byte(oldRefreshHash)) != 1) {
		return nil, domain.ErrUnauthenticated
	}
	if err != nil {
		return nil, err
	}
	if _, err := tx.ExecContext(ctx,
		`UPDATE auth_sessions SET refresh_token_hash = $2, expires_at = $3 WHERE id = $1`,
		sessionID, newRefreshHash, expiresAt,
	); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) AuthenticateSession(ctx context.Context, userID, sessionID string) (*domain.User, error) {
	var user domain.User
	err := s.db.QueryRowContext(ctx, `
SELECT u.id, u.email, u.login
FROM auth_sessions s
JOIN users u ON u.id = s.user_id
WHERE s.id = $1 AND s.user_id = $2 AND s.expires_at > now()`, sessionID, userID).
		Scan(&user.ID, &user.Email, &user.Login)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUnauthenticated
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Store) RevokeSession(ctx context.Context, userID, sessionID, refreshHash string) error {
	_, err := s.db.ExecContext(ctx,
		`DELETE FROM auth_sessions WHERE user_id = $1 AND id = $2 AND refresh_token_hash = $3`,
		userID, sessionID, refreshHash,
	)
	return err
}

func (s *Store) PasswordHash(ctx context.Context, userID string) (string, error) {
	var hash string
	err := s.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id = $1`, userID).Scan(&hash)
	if errors.Is(err, sql.ErrNoRows) {
		return "", domain.ErrNotFound
	}
	return hash, err
}

func (s *Store) ChangePassword(ctx context.Context, userID, sessionID, passwordHash string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := lockActiveSession(ctx, tx, userID, sessionID); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `UPDATE users SET password_hash = $2, updated_at = now() WHERE id = $1`, userID, passwordHash); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM auth_sessions WHERE user_id = $1 AND id <> $2`, userID, sessionID); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *Store) ChangeLogin(ctx context.Context, userID, sessionID, login string) (*domain.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if err := lockActiveSession(ctx, tx, userID, sessionID); err != nil {
		return nil, err
	}
	var user domain.User
	err = tx.QueryRowContext(ctx, `UPDATE users SET login = $2, updated_at = now() WHERE id = $1 RETURNING id, email, login`, userID, login).
		Scan(&user.ID, &user.Email, &user.Login)
	if err != nil {
		return nil, mapPostgresError(err)
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &user, nil
}

func lockActiveSession(ctx context.Context, tx *sql.Tx, userID, sessionID string) error {
	var id string
	err := tx.QueryRowContext(ctx,
		`SELECT id FROM auth_sessions WHERE user_id = $1 AND id = $2 AND expires_at > now() FOR UPDATE`,
		userID, sessionID,
	).Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrUnauthenticated
	}
	return err
}

func (s *Store) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT ci.product_id, p.slug, p.title, ci.quantity, p.price, p.pricesale,
       COALESCE(p.photos[1], ''),
       COALESCE(NULLIF(p.pricesale, 0), p.price)::bigint * ci.quantity
FROM cart_items ci
JOIN products p ON p.id = ci.product_id
WHERE ci.user_id = $1
ORDER BY ci.created_at, ci.product_id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cart := &domain.Cart{Items: make([]*domain.CartItem, 0)}
	for rows.Next() {
		item := &domain.CartItem{}
		if err := rows.Scan(&item.ProductID, &item.Slug, &item.Title, &item.Quantity, &item.Price, &item.PriceSale, &item.Image, &item.LineTotal); err != nil {
			return nil, err
		}
		cart.Items = append(cart.Items, item)
		cart.TotalQuantity += item.Quantity
		cart.TotalPrice += item.LineTotal
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *Store) AddCartItem(ctx context.Context, userID, productID string, quantity int32) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM products WHERE id = $1)`, productID).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return domain.ErrNotFound
	}
	var resultingQuantity int32
	err = tx.QueryRowContext(ctx, `
INSERT INTO cart_items (user_id, product_id, quantity)
VALUES ($1, $2, $3)
ON CONFLICT (user_id, product_id) DO UPDATE
SET quantity = cart_items.quantity + EXCLUDED.quantity, updated_at = now()
WHERE cart_items.quantity + EXCLUDED.quantity <= 99
RETURNING quantity`, userID, productID, quantity).Scan(&resultingQuantity)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrCartLimit
	}
	if err != nil {
		return mapPostgresError(err)
	}
	return tx.Commit()
}

func (s *Store) UpdateCartItem(ctx context.Context, userID, productID string, quantity int32) error {
	result, err := s.db.ExecContext(ctx,
		`UPDATE cart_items SET quantity = $3, updated_at = now() WHERE user_id = $1 AND product_id = $2`,
		userID, productID, quantity,
	)
	if err != nil {
		return mapPostgresError(err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) RemoveCartItem(ctx context.Context, userID, productID string) error {
	result, err := s.db.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id = $1 AND product_id = $2`, userID, productID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (s *Store) ClearCart(ctx context.Context, userID string) error {
	_, err := s.db.ExecContext(ctx, `DELETE FROM cart_items WHERE user_id = $1`, userID)
	return err
}

func mapPostgresError(err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return domain.ErrAlreadyExists
		case "23503":
			return domain.ErrNotFound
		case "23514", "22001":
			return domain.ErrInvalidInput
		}
	}
	return err
}
