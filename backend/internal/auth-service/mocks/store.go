package mocks

import (
	"context"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
)

type Store struct {
	CreateUserAndSessionFn func(ctx context.Context, user *domain.User, passwordHash, sessionID, refreshHash string, expiresAt time.Time) error
	FindUserForLoginFn     func(ctx context.Context, identifier string) (*domain.User, string, error)
	CreateSessionFn        func(ctx context.Context, userID, sessionID, refreshHash string, expiresAt time.Time) error
	RotateSessionFn        func(ctx context.Context, userID, sessionID, oldRefreshHash, newRefreshHash string, expiresAt time.Time) (*domain.User, error)
	AuthenticateSessionFn  func(ctx context.Context, userID, sessionID string) (*domain.User, error)
	RevokeSessionFn        func(ctx context.Context, userID, sessionID, refreshHash string) error
	PasswordHashFn         func(ctx context.Context, userID string) (string, error)
	ChangePasswordFn       func(ctx context.Context, userID, sessionID, passwordHash string) error
	ChangeLoginFn          func(ctx context.Context, userID, sessionID, login string) (*domain.User, error)
	GetCartFn              func(ctx context.Context, userID string) (*domain.Cart, error)
	AddCartItemFn          func(ctx context.Context, userID, productID string, quantity int32) error
	UpdateCartItemFn       func(ctx context.Context, userID, productID string, quantity int32) error
	RemoveCartItemFn       func(ctx context.Context, userID, productID string) error
	ClearCartFn            func(ctx context.Context, userID string) error
}

func (s Store) CreateUserAndSession(ctx context.Context, user *domain.User, passwordHash, sessionID, refreshHash string, expiresAt time.Time) error {
	return s.CreateUserAndSessionFn(ctx, user, passwordHash, sessionID, refreshHash, expiresAt)
}
func (s Store) FindUserForLogin(ctx context.Context, identifier string) (*domain.User, string, error) {
	return s.FindUserForLoginFn(ctx, identifier)
}
func (s Store) CreateSession(ctx context.Context, userID, sessionID, refreshHash string, expiresAt time.Time) error {
	return s.CreateSessionFn(ctx, userID, sessionID, refreshHash, expiresAt)
}
func (s Store) RotateSession(ctx context.Context, userID, sessionID, oldRefreshHash, newRefreshHash string, expiresAt time.Time) (*domain.User, error) {
	return s.RotateSessionFn(ctx, userID, sessionID, oldRefreshHash, newRefreshHash, expiresAt)
}
func (s Store) AuthenticateSession(ctx context.Context, userID, sessionID string) (*domain.User, error) {
	return s.AuthenticateSessionFn(ctx, userID, sessionID)
}
func (s Store) RevokeSession(ctx context.Context, userID, sessionID, refreshHash string) error {
	return s.RevokeSessionFn(ctx, userID, sessionID, refreshHash)
}
func (s Store) PasswordHash(ctx context.Context, userID string) (string, error) {
	return s.PasswordHashFn(ctx, userID)
}
func (s Store) ChangePassword(ctx context.Context, userID, sessionID, passwordHash string) error {
	return s.ChangePasswordFn(ctx, userID, sessionID, passwordHash)
}
func (s Store) ChangeLogin(ctx context.Context, userID, sessionID, login string) (*domain.User, error) {
	return s.ChangeLoginFn(ctx, userID, sessionID, login)
}
func (s Store) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return s.GetCartFn(ctx, userID)
}
func (s Store) AddCartItem(ctx context.Context, userID, productID string, quantity int32) error {
	return s.AddCartItemFn(ctx, userID, productID, quantity)
}
func (s Store) UpdateCartItem(ctx context.Context, userID, productID string, quantity int32) error {
	return s.UpdateCartItemFn(ctx, userID, productID, quantity)
}
func (s Store) RemoveCartItem(ctx context.Context, userID, productID string) error {
	return s.RemoveCartItemFn(ctx, userID, productID)
}
func (s Store) ClearCart(ctx context.Context, userID string) error { return s.ClearCartFn(ctx, userID) }
