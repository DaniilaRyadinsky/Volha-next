package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/autumnterror/volha-backend/internal/auth-service/mocks"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func authFixture(store Store) *Service {
	return New(store, strings.Repeat("s", 32), 15*time.Minute, 24*time.Hour)
}
func TestRegisterHashesAndNormalizes(t *testing.T) {
	var stored *domain.User
	var refreshHash string
	store := mocks.Store{CreateUserAndSessionFn: func(_ context.Context, user *domain.User, hash, sid, rhash string, expires time.Time) error {
		stored = user
		refreshHash = rhash
		require.NoError(t, bcrypt.CompareHashAndPassword([]byte(hash), []byte("Password123")))
		require.Len(t, sid, 20)
		require.True(t, expires.After(time.Now()))
		return nil
	}}
	svc := authFixture(store)
	pair, err := svc.Register(context.Background(), " USER@EXAMPLE.COM ", " My_Login ", "Password123")
	require.NoError(t, err)
	require.Equal(t, "user@example.com", stored.Email)
	require.Equal(t, "my_login", stored.Login)
	require.Equal(t, hashToken(pair.RefreshToken), refreshHash)
	require.Equal(t, int64(900), pair.AccessExpiresIn)
	require.Same(t, stored, pair.User)
	claims, err := svc.verifyToken(pair.AccessToken, "access", false)
	require.NoError(t, err)
	require.Equal(t, stored.ID, claims.Subject)
	for _, input := range [][3]string{{"bad", "login", "Password123"}, {"a@b.com", "!", "Password123"}, {"a@b.com", "login", "short"}, {"a@b.com", "login", "abcdefghijk"}} {
		_, err := svc.Register(context.Background(), input[0], input[1], input[2])
		require.ErrorIs(t, err, domain.ErrInvalidInput)
	}
}
func TestLoginRefreshAuthenticationAndLogout(t *testing.T) {
	ctx := context.Background()
	hash, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.MinCost)
	require.NoError(t, err)
	user := &domain.User{ID: xid.New().String(), Email: "u@e.com", Login: "user"}
	var sid string
	revoked := false
	store := mocks.Store{
		FindUserForLoginFn: func(_ context.Context, id string) (*domain.User, string, error) {
			require.Equal(t, "user", id)
			return user, string(hash), nil
		},
		CreateSessionFn: func(_ context.Context, uid, session, hash string, expires time.Time) error {
			require.Equal(t, user.ID, uid)
			sid = session
			return nil
		},
		RotateSessionFn: func(_ context.Context, uid, session, oldHash, newHash string, expires time.Time) (*domain.User, error) {
			require.Equal(t, sid, session)
			require.NotEqual(t, oldHash, newHash)
			return user, nil
		},
		AuthenticateSessionFn: func(_ context.Context, uid, session string) (*domain.User, error) {
			require.Equal(t, user.ID, uid)
			require.Equal(t, sid, session)
			if revoked {
				return nil, domain.ErrUnauthenticated
			}
			return user, nil
		},
		RevokeSessionFn: func(_ context.Context, uid, session, hash string) error {
			require.Equal(t, sid, session)
			revoked = true
			return nil
		}}
	svc := authFixture(store)
	pair, err := svc.Login(ctx, " user ", "Password123")
	require.NoError(t, err)
	_, err = svc.Login(ctx, "user", "wrong")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	_, err = svc.Login(ctx, "", "Password123")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	principal, err := svc.Authenticate(ctx, pair.AccessToken)
	require.NoError(t, err)
	require.Equal(t, sid, principal.SessionID)
	fresh, err := svc.Refresh(ctx, pair.RefreshToken)
	require.NoError(t, err)
	require.NotEqual(t, pair.RefreshToken, fresh.RefreshToken)
	_, err = svc.Refresh(ctx, pair.AccessToken)
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	require.NoError(t, svc.Logout(ctx, "invalid"))
	require.NoError(t, svc.Logout(ctx, fresh.RefreshToken))
	_, err = svc.Authenticate(ctx, pair.AccessToken)
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
}
func TestJWTRejectsTamperingExpiryAndWrongUse(t *testing.T) {
	svc := authFixture(mocks.Store{})
	now := time.Now()
	valid, err := svc.signToken("user", "session", "access", now, now.Add(time.Minute))
	require.NoError(t, err)
	expired, err := svc.signToken("user", "session", "access", now.Add(-2*time.Minute), now.Add(-time.Minute))
	require.NoError(t, err)
	future, err := svc.signToken("user", "session", "access", now.Add(2*time.Minute), now.Add(3*time.Minute))
	require.NoError(t, err)
	for _, token := range []string{"", "bad.parts", valid + "x", expired, future, "!!!.!!!.!!!"} {
		_, err := svc.verifyToken(token, "access", false)
		require.Error(t, err)
	}
	_, err = svc.verifyToken(valid, "refresh", false)
	require.Error(t, err)
	_, err = svc.verifyToken(expired, "access", true)
	require.NoError(t, err)
	other := New(mocks.Store{}, "different secret", time.Minute, time.Hour)
	_, err = other.verifyToken(valid, "access", false)
	require.Error(t, err)
}
func TestAccountChangesAndCartRules(t *testing.T) {
	ctx := context.Background()
	hash, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.MinCost)
	require.NoError(t, err)
	user := &domain.User{ID: "user", Login: "new_login"}
	cart := &domain.Cart{Items: []*domain.CartItem{}, TotalQuantity: 2}
	pid := xid.New().String()
	calls := 0
	store := mocks.Store{PasswordHashFn: func(context.Context, string) (string, error) { return string(hash), nil }, ChangePasswordFn: func(_ context.Context, uid, sid, hashed string) error {
		require.NoError(t, bcrypt.CompareHashAndPassword([]byte(hashed), []byte("NewPassword123")))
		calls++
		return nil
	}, AuthenticateSessionFn: func(context.Context, string, string) (*domain.User, error) { return user, nil }, ChangeLoginFn: func(_ context.Context, uid, sid, login string) (*domain.User, error) {
		require.Equal(t, "new_login", login)
		calls++
		return user, nil
	}, GetCartFn: func(context.Context, string) (*domain.Cart, error) { return cart, nil }, AddCartItemFn: func(_ context.Context, uid, id string, q int32) error {
		require.Equal(t, pid, id)
		require.Equal(t, int32(2), q)
		calls++
		return nil
	}, UpdateCartItemFn: func(context.Context, string, string, int32) error { calls++; return nil }, RemoveCartItemFn: func(context.Context, string, string) error { calls++; return nil }, ClearCartFn: func(context.Context, string) error { calls++; return nil }}
	svc := authFixture(store)
	got, err := svc.ChangePassword(ctx, "user", "session", "Password123", "NewPassword123")
	require.NoError(t, err)
	require.Same(t, user, got)
	_, err = svc.ChangePassword(ctx, "user", "session", "wrong", "NewPassword123")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	_, err = svc.ChangePassword(ctx, "user", "session", "Password123", "Password123")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
	got, err = svc.ChangeLogin(ctx, "user", "session", "Password123", " NEW_LOGIN ")
	require.NoError(t, err)
	require.Same(t, user, got)
	_, err = svc.ChangeLogin(ctx, "user", "session", "wrong", "new_login")
	require.ErrorIs(t, err, domain.ErrUnauthenticated)
	current, err := svc.GetCart(ctx, "user")
	require.NoError(t, err)
	require.Same(t, cart, current)
	current, err = svc.AddCartItem(ctx, "user", pid, 2)
	require.NoError(t, err)
	require.Same(t, cart, current)
	_, err = svc.UpdateCartItem(ctx, "user", pid, 1)
	require.NoError(t, err)
	_, err = svc.RemoveCartItem(ctx, "user", pid)
	require.NoError(t, err)
	_, err = svc.ClearCart(ctx, "user")
	require.NoError(t, err)
	require.Equal(t, 6, calls)
	for _, quantity := range []int32{0, 100} {
		_, err = svc.AddCartItem(ctx, "user", pid, quantity)
		require.ErrorIs(t, err, domain.ErrInvalidInput)
		_, err = svc.UpdateCartItem(ctx, "user", pid, quantity)
		require.ErrorIs(t, err, domain.ErrInvalidInput)
	}
	_, err = svc.RemoveCartItem(ctx, "user", "../invalid")
	require.ErrorIs(t, err, domain.ErrInvalidInput)
}
func TestServicePropagatesStoreFailure(t *testing.T) {
	boom := errors.New("storage down")
	svc := authFixture(mocks.Store{FindUserForLoginFn: func(context.Context, string) (*domain.User, string, error) { return nil, "", boom }, GetCartFn: func(context.Context, string) (*domain.Cart, error) { return nil, boom }})
	_, err := svc.Login(context.Background(), "user", "password")
	require.ErrorIs(t, err, boom)
	_, err = svc.GetCart(context.Background(), "user")
	require.ErrorIs(t, err, boom)
}
