package mw

import (
	"net/http"
	"strings"

	"github.com/autumnterror/volha-backend/internal/gateway/config"

	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

func DevAdminAuth(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if cfg.Mode != "local" && cfg.Mode != "development" && cfg.Mode != "DEV" {
				return c.JSON(http.StatusForbidden, map[string]string{
					"error": "admin endpoints are disabled",
				})
			}
			header := c.Request().Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" || token != cfg.DevAdminToken {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "unauthorized",
				})
			}
			return next(c)
		}
	}
}

func CheckId() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			id := c.QueryParam("id")
			if id == "" {
				return next(c)
			}
			if _, err := xid.FromString(id); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "bad id",
				})
			}
			return next(c)
		}
	}
}
