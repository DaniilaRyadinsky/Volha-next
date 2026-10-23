package service

import (
	"context"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
)

type Store interface {
	CreateUserAndSession(ctx context.Context, user *domain.User, passwordHash, sessionID, refreshHash string, expiresAt time.Time) error
	FindUserForLogin(ctx context.Context, identifier string) (*domain.User, string, error)
	CreateSession(ctx context.Context, userID, sessionID, refreshHash string, expiresAt time.Time) error
	RotateSession(ctx context.Context, userID, sessionID, oldRefreshHash, newRefreshHash string, expiresAt time.Time) (*domain.User, error)
	AuthenticateSession(ctx context.Context, userID, sessionID string) (*domain.User, error)
	RevokeSession(ctx context.Context, userID, sessionID, refreshHash string) error
	PasswordHash(ctx context.Context, userID string) (string, error)
	ChangePassword(ctx context.Context, userID, sessionID, passwordHash string) error
	ChangeLogin(ctx context.Context, userID, sessionID, login string) (*domain.User, error)
	GetCart(ctx context.Context, userID string) (*domain.Cart, error)
	AddCartItem(ctx context.Context, userID, productID string, quantity int32) error
	UpdateCartItem(ctx context.Context, userID, productID string, quantity int32) error
	RemoveCartItem(ctx context.Context, userID, productID string) error
	ClearCart(ctx context.Context, userID string) error
}
