package mw

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/gateway/net/authcookie"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authFunc func(context.Context, string) (*productsRPC.AuthenticatedUser, error)

func (f authFunc) ValidateAccessToken(ctx context.Context, token string) (*productsRPC.AuthenticatedUser, error) {
	return f(ctx, token)
}
func TestAccessCookieMiddleware(t *testing.T) {
	for _, tc := range []struct {
		name, token string
		code        codes.Code
		want        int
	}{{"missing", "", codes.OK, 401}, {"valid", "token", codes.OK, 204}, {"revoked", "token", codes.Unauthenticated, 401}, {"backend down", "token", codes.Unavailable, 503}} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			if tc.token != "" {
				req.AddCookie(&http.Cookie{Name: authcookie.AccessTokenName, Value: tc.token})
			}
			rec := httptest.NewRecorder()
			c := echo.New().NewContext(req, rec)
			auth := authFunc(func(_ context.Context, token string) (*productsRPC.AuthenticatedUser, error) {
				require.Equal(t, tc.token, token)
				if tc.code != codes.OK {
					return nil, status.Error(tc.code, "auth failure")
				}
				return &productsRPC.AuthenticatedUser{User: &productsRPC.AuthUser{Id: "u"}}, nil
			})
			require.NoError(t, RequireAuth(auth)(func(c echo.Context) error {
				require.Equal(t, "u", CurrentUser(c).GetUser().GetId())
				return c.NoContent(204)
			})(c))
			require.Equal(t, tc.want, rec.Code)
		})
	}
}
func TestAdminAndIDMiddleware(t *testing.T) {
	for _, tc := range []struct {
		mode, token string
		want        int
	}{{"PROD", "secret", 403}, {"local", "wrong", 401}, {"local", "secret", 204}, {"DEV", "secret", 204}} {
		req := httptest.NewRequest("POST", "/", nil)
		req.Header.Set("Authorization", "Bearer "+tc.token)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		require.NoError(t, DevAdminAuth(&config.Config{Mode: tc.mode, DevAdminToken: "secret"})(func(c echo.Context) error { return c.NoContent(204) })(c))
		require.Equal(t, tc.want, rec.Code)
	}
	for _, id := range []string{"", xid.New().String(), "invalid"} {
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(httptest.NewRequest("GET", "/?id="+id, nil), rec)
		require.NoError(t, CheckId()(func(c echo.Context) error { return c.NoContent(204) })(c))
		if id == "invalid" {
			require.Equal(t, 400, rec.Code)
		} else {
			require.Equal(t, 204, rec.Code)
		}
	}
}
