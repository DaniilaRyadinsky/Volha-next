package mw

import (
	"context"
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/authcookie"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Authenticator interface {
	ValidateAccessToken(context.Context, string) (*productsRPC.AuthenticatedUser, error)
}

func RequireAuth(auth Authenticator) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie(authcookie.AccessTokenName)
			if err != nil || cookie.Value == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
			}
			principal, err := auth.ValidateAccessToken(c.Request().Context(), cookie.Value)
			if err != nil {
				if status.Code(err) == codes.Unauthenticated {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "unauthenticated"})
				}
				return c.JSON(http.StatusServiceUnavailable, map[string]string{"error": "auth service unavailable"})
			}
			c.Set(authcookie.PrincipalContextKey, principal)
			return next(c)
		}
	}
}

func CurrentUser(c echo.Context) *productsRPC.AuthenticatedUser {
	principal, _ := c.Get(authcookie.PrincipalContextKey).(*productsRPC.AuthenticatedUser)
	return principal
}
