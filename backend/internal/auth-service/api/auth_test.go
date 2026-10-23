package api

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/auth-service/config"
	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/autumnterror/volha-backend/internal/auth-service/mocks"
	"github.com/autumnterror/volha-backend/internal/auth-service/service"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuthTransport(t *testing.T) {
	ctx := context.Background()
	user := &domain.User{ID: xid.New().String(), Email: "u@example.com", Login: "login"}
	hash, err := bcrypt.GenerateFromPassword([]byte("Password123"), bcrypt.MinCost)
	require.NoError(t, err)
	cart := &domain.Cart{Items: []*domain.CartItem{{ProductID: "p", Quantity: 2, Price: 10, LineTotal: 20}}, TotalQuantity: 2, TotalPrice: 20}
	store := mocks.Store{CreateUserAndSessionFn: func(context.Context, *domain.User, string, string, string, time.Time) error { return nil }, FindUserForLoginFn: func(context.Context, string) (*domain.User, string, error) { return user, string(hash), nil }, CreateSessionFn: func(context.Context, string, string, string, time.Time) error { return nil }, RotateSessionFn: func(context.Context, string, string, string, string, time.Time) (*domain.User, error) {
		return user, nil
	}, RevokeSessionFn: func(context.Context, string, string, string) error { return nil }, AuthenticateSessionFn: func(context.Context, string, string) (*domain.User, error) { return user, nil }, PasswordHashFn: func(context.Context, string) (string, error) { return string(hash), nil }, ChangePasswordFn: func(context.Context, string, string, string) error { return nil }, ChangeLoginFn: func(context.Context, string, string, string) (*domain.User, error) { return user, nil }, GetCartFn: func(context.Context, string) (*domain.Cart, error) { return cart, nil }, AddCartItemFn: func(context.Context, string, string, int32) error { return nil }, UpdateCartItemFn: func(context.Context, string, string, int32) error { return nil }, RemoveCartItemFn: func(context.Context, string, string) error { return nil }, ClearCartFn: func(context.Context, string) error { return nil }}
	s := &serverAPI{auth: service.New(store, strings.Repeat("s", 32), time.Minute, time.Hour)}
	registered, err := s.Register(ctx, &productsRPC.RegisterRequest{Email: user.Email, Login: user.Login, Password: "Password123"})
	require.NoError(t, err)
	require.NotEmpty(t, registered.GetAccessToken())
	pair, err := s.Login(ctx, &productsRPC.LoginRequest{Identifier: "login", Password: "Password123"})
	require.NoError(t, err)
	principal, err := s.Authenticate(ctx, &productsRPC.AuthenticateRequest{AccessToken: pair.GetAccessToken()})
	require.NoError(t, err)
	require.Equal(t, user.ID, principal.GetUser().GetId())
	fresh, err := s.Refresh(ctx, &productsRPC.RefreshRequest{RefreshToken: pair.GetRefreshToken()})
	require.NoError(t, err)
	require.NotEqual(t, pair.GetRefreshToken(), fresh.GetRefreshToken())
	_, err = s.Logout(ctx, &productsRPC.LogoutRequest{RefreshToken: fresh.GetRefreshToken()})
	require.NoError(t, err)
	_, err = s.ChangePassword(ctx, &productsRPC.ChangePasswordRequest{UserId: user.ID, SessionId: "sid", OldPassword: "Password123", NewPassword: "OtherPassword123"})
	require.NoError(t, err)
	_, err = s.ChangeLogin(ctx, &productsRPC.ChangeLoginRequest{UserId: user.ID, SessionId: "sid", Password: "Password123", NewLogin: "other_login"})
	require.NoError(t, err)
	got, err := s.GetCart(ctx, &productsRPC.CartRequest{UserId: user.ID})
	require.NoError(t, err)
	require.Equal(t, int64(20), got.GetTotalPrice())
	id := xid.New().String()
	_, err = s.AddCartItem(ctx, &productsRPC.CartItemRequest{UserId: user.ID, ProductId: id, Quantity: 1})
	require.NoError(t, err)
	_, err = s.UpdateCartItem(ctx, &productsRPC.CartItemRequest{UserId: user.ID, ProductId: id, Quantity: 2})
	require.NoError(t, err)
	_, err = s.RemoveCartItem(ctx, &productsRPC.CartItemKey{UserId: user.ID, ProductId: id})
	require.NoError(t, err)
	_, err = s.ClearCart(ctx, &productsRPC.CartRequest{UserId: user.ID})
	require.NoError(t, err)
	_, err = s.Register(ctx, nil)
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = s.Authenticate(ctx, &productsRPC.AuthenticateRequest{AccessToken: "wrong"})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	require.Nil(t, userToRPC(nil))
	require.Empty(t, cartToRPC(&domain.Cart{}).GetItems())
}
func TestAuthErrorsAndHealth(t *testing.T) {
	for input, want := range map[error]codes.Code{domain.ErrInvalidInput: codes.InvalidArgument, domain.ErrAlreadyExists: codes.AlreadyExists, domain.ErrUnauthenticated: codes.Unauthenticated, domain.ErrNotFound: codes.NotFound, domain.ErrCartLimit: codes.ResourceExhausted, errors.New("sensitive detail"): codes.Internal} {
		err := mapError(context.Background(), "test", input)
		require.Equal(t, want, status.Code(err))
		if want == codes.Internal {
			require.NotContains(t, err.Error(), "sensitive detail")
		}
	}
	_, err := (healthServer{}).Check(context.Background(), nil)
	require.NoError(t, err)
	_, err = (healthServer{ready: func(context.Context) error { return errors.New("down") }}).Check(context.Background(), nil)
	require.Equal(t, codes.Unavailable, status.Code(err))
	app := New(&config.Config{Port: -1}, service.New(mocks.Store{}, "secret", time.Minute, time.Hour), nil, observability.Config{}, observability.NewMetrics("test", "dev", "test"))
	require.Error(t, app.Run())
	app.Stop(context.Background())
}
