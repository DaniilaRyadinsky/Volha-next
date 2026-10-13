package handlers

import (
	"context"
	"net/http"
	"time"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/authcookie"
	"github.com/autumnterror/volha-backend/internal/gateway/net/mw"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type registerRequest struct {
	Email    string `json:"email" example:"user@example.com"`
	Login    string `json:"login" example:"user_name"`
	Password string `json:"password" example:"StrongPassword123"`
}

type loginRequest struct {
	Identifier string `json:"identifier" example:"user@example.com"`
	Password   string `json:"password" example:"StrongPassword123"`
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" example:"OldPassword123"`
	NewPassword string `json:"new_password" example:"NewPassword123"`
}

type changeLoginRequest struct {
	Password string `json:"password" example:"StrongPassword123"`
	NewLogin string `json:"new_login" example:"new_login"`
}

type cartItemRequest struct {
	ProductID string `json:"product_id" example:"c61j5e0k4pc7j3n1k2g0"`
	Quantity  int32  `json:"quantity" example:"2"`
}

// Register godoc
// @Summary Регистрация пользователя
// @Description Access и refresh JWT устанавливаются в HttpOnly cookies и не возвращаются в JSON.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body registerRequest true "Данные пользователя"
// @Success 201 {object} productsRPC.AuthUser
// @Failure 400 {object} view.ErrorResponse
// @Failure 409 {object} view.ErrorResponse
// @Router /api/auth/register [post]
func (a *Apis) Register(c echo.Context) error {
	var req registerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	result, err := a.apiAuth.API.Register(ctx, &productsRPC.RegisterRequest{
		Email: req.Email, Login: req.Login, Password: req.Password,
	})
	if err != nil {
		return a.authError(c, ctx, "Register", err)
	}
	setAuthCookies(c, result, a.cfg.CookieSecure)
	return c.JSON(http.StatusCreated, result.GetUser())
}

// Login godoc
// @Summary Вход по логину или email
// @Description Access и refresh JWT устанавливаются в HttpOnly cookies и не возвращаются в JSON.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body loginRequest true "Учетные данные"
// @Success 200 {object} productsRPC.AuthUser
// @Failure 400 {object} view.ErrorResponse
// @Failure 401 {object} view.ErrorResponse
// @Router /api/auth/login [post]
func (a *Apis) Login(c echo.Context) error {
	var req loginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	result, err := a.apiAuth.API.Login(ctx, &productsRPC.LoginRequest{
		Identifier: req.Identifier, Password: req.Password,
	})
	if err != nil {
		return a.authError(c, ctx, "Login", err)
	}
	setAuthCookies(c, result, a.cfg.CookieSecure)
	return c.JSON(http.StatusOK, result.GetUser())
}

// Refresh godoc
// @Summary Обновить пару JWT токенов
// @Description Refresh token ротируется; новая пара JWT устанавливается в HttpOnly cookies.
// @Tags auth
// @Produce json
// @Success 200 {object} productsRPC.AuthUser
// @Failure 401 {object} view.ErrorResponse
// @Router /api/auth/refresh [post]
func (a *Apis) Refresh(c echo.Context) error {
	cookie, err := c.Cookie(authcookie.RefreshTokenName)
	if err != nil || cookie.Value == "" {
		clearAuthCookies(c, a.cfg.CookieSecure)
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	result, err := a.apiAuth.API.Refresh(ctx, &productsRPC.RefreshRequest{RefreshToken: cookie.Value})
	if err != nil {
		clearAuthCookies(c, a.cfg.CookieSecure)
		return a.authError(c, ctx, "Refresh", err)
	}
	setAuthCookies(c, result, a.cfg.CookieSecure)
	return c.JSON(http.StatusOK, result.GetUser())
}

// Logout godoc
// @Summary Завершить сессию
// @Description Refresh-сессия отзывается, access и refresh cookies очищаются.
// @Tags auth
// @Produce json
// @Success 200 {object} view.MessageResponse
// @Router /api/auth/logout [post]
func (a *Apis) Logout(c echo.Context) error {
	cookie, _ := c.Cookie(authcookie.RefreshTokenName)
	if cookie != nil && cookie.Value != "" {
		ctx, cancel := requestContext(c)
		defer cancel()
		if _, err := a.apiAuth.API.Logout(ctx, &productsRPC.LogoutRequest{RefreshToken: cookie.Value}); err != nil {
			clearAuthCookies(c, a.cfg.CookieSecure)
			return a.authError(c, ctx, "Logout", err)
		}
	}
	clearAuthCookies(c, a.cfg.CookieSecure)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "logged out"})
}

// Me godoc
// @Summary Текущий пользователь
// @Tags auth
// @Produce json
// @Success 200 {object} productsRPC.AuthUser
// @Failure 401 {object} view.ErrorResponse
// @Router /api/auth/me [get]
func (a *Apis) Me(c echo.Context) error {
	principal := mw.CurrentUser(c)
	if principal == nil || principal.GetUser() == nil {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	return c.JSON(http.StatusOK, principal.GetUser())
}

// ChangePassword godoc
// @Summary Сменить пароль
// @Tags auth
// @Accept json
// @Produce json
// @Param request body changePasswordRequest true "Текущий и новый пароль"
// @Success 200 {object} productsRPC.AuthUser
// @Failure 400 {object} view.ErrorResponse
// @Failure 401 {object} view.ErrorResponse
// @Router /api/auth/change-password [put]
func (a *Apis) ChangePassword(c echo.Context) error {
	var req changePasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	principal := mw.CurrentUser(c)
	if principal == nil || principal.GetUser() == nil {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	user, err := a.apiAuth.API.ChangePassword(ctx, &productsRPC.ChangePasswordRequest{
		UserId: principal.GetUser().GetId(), SessionId: principal.GetSessionId(),
		OldPassword: req.OldPassword, NewPassword: req.NewPassword,
	})
	if err != nil {
		return a.authError(c, ctx, "ChangePassword", err)
	}
	return c.JSON(http.StatusOK, user)
}

// ChangeLogin godoc
// @Summary Сменить логин
// @Tags auth
// @Accept json
// @Produce json
// @Param request body changeLoginRequest true "Пароль и новый логин"
// @Success 200 {object} productsRPC.AuthUser
// @Failure 400 {object} view.ErrorResponse
// @Failure 401 {object} view.ErrorResponse
// @Failure 409 {object} view.ErrorResponse
// @Router /api/auth/change-login [put]
func (a *Apis) ChangeLogin(c echo.Context) error {
	var req changeLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	principal := mw.CurrentUser(c)
	if principal == nil || principal.GetUser() == nil {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	user, err := a.apiAuth.API.ChangeLogin(ctx, &productsRPC.ChangeLoginRequest{
		UserId: principal.GetUser().GetId(), SessionId: principal.GetSessionId(),
		Password: req.Password, NewLogin: req.NewLogin,
	})
	if err != nil {
		return a.authError(c, ctx, "ChangeLogin", err)
	}
	return c.JSON(http.StatusOK, user)
}

// GetCart godoc
// @Summary Получить корзину пользователя
// @Tags cart
// @Produce json
// @Success 200 {object} productsRPC.Cart
// @Failure 401 {object} view.ErrorResponse
// @Router /api/cart [get]
func (a *Apis) GetCart(c echo.Context) error {
	principal, ok := authenticatedPrincipal(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cart, err := a.apiAuth.API.GetCart(ctx, &productsRPC.CartRequest{UserId: principal.GetUser().GetId()})
	if err != nil {
		return a.authError(c, ctx, "GetCart", err)
	}
	return c.JSON(http.StatusOK, cart)
}

// AddCartItem godoc
// @Summary Добавить товар в корзину
// @Tags cart
// @Accept json
// @Produce json
// @Param request body cartItemRequest true "Товар и количество"
// @Success 200 {object} productsRPC.Cart
// @Failure 400 {object} view.ErrorResponse
// @Failure 401 {object} view.ErrorResponse
// @Router /api/cart/items [post]
func (a *Apis) AddCartItem(c echo.Context) error {
	var req cartItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	principal, ok := authenticatedPrincipal(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cart, err := a.apiAuth.API.AddCartItem(ctx, &productsRPC.CartItemRequest{
		UserId: principal.GetUser().GetId(), ProductId: req.ProductID, Quantity: req.Quantity,
	})
	if err != nil {
		return a.authError(c, ctx, "AddCartItem", err)
	}
	return c.JSON(http.StatusOK, cart)
}

// UpdateCartItem godoc
// @Summary Изменить количество товара в корзине
// @Tags cart
// @Accept json
// @Produce json
// @Param product_id path string true "ID товара"
// @Param request body cartItemRequest true "Новое количество"
// @Success 200 {object} productsRPC.Cart
// @Failure 400 {object} view.ErrorResponse
// @Failure 401 {object} view.ErrorResponse
// @Router /api/cart/items/{product_id} [put]
func (a *Apis) UpdateCartItem(c echo.Context) error {
	var req cartItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	principal, ok := authenticatedPrincipal(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cart, err := a.apiAuth.API.UpdateCartItem(ctx, &productsRPC.CartItemRequest{
		UserId: principal.GetUser().GetId(), ProductId: c.Param("product_id"), Quantity: req.Quantity,
	})
	if err != nil {
		return a.authError(c, ctx, "UpdateCartItem", err)
	}
	return c.JSON(http.StatusOK, cart)
}

// RemoveCartItem godoc
// @Summary Удалить товар из корзины
// @Tags cart
// @Produce json
// @Param product_id path string true "ID товара"
// @Success 200 {object} productsRPC.Cart
// @Failure 401 {object} view.ErrorResponse
// @Router /api/cart/items/{product_id} [delete]
func (a *Apis) RemoveCartItem(c echo.Context) error {
	principal, ok := authenticatedPrincipal(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cart, err := a.apiAuth.API.RemoveCartItem(ctx, &productsRPC.CartItemKey{
		UserId: principal.GetUser().GetId(), ProductId: c.Param("product_id"),
	})
	if err != nil {
		return a.authError(c, ctx, "RemoveCartItem", err)
	}
	return c.JSON(http.StatusOK, cart)
}

// ClearCart godoc
// @Summary Очистить корзину
// @Tags cart
// @Produce json
// @Success 200 {object} productsRPC.Cart
// @Failure 401 {object} view.ErrorResponse
// @Router /api/cart [delete]
func (a *Apis) ClearCart(c echo.Context) error {
	principal, ok := authenticatedPrincipal(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, view.ErrorResponse{Error: "unauthenticated"})
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	cart, err := a.apiAuth.API.ClearCart(ctx, &productsRPC.CartRequest{UserId: principal.GetUser().GetId()})
	if err != nil {
		return a.authError(c, ctx, "ClearCart", err)
	}
	return c.JSON(http.StatusOK, cart)
}

func authenticatedPrincipal(c echo.Context) (*productsRPC.AuthenticatedUser, bool) {
	principal := mw.CurrentUser(c)
	return principal, principal != nil && principal.GetUser() != nil && principal.GetUser().GetId() != ""
}

func setAuthCookies(c echo.Context, result *productsRPC.AuthResponse, secure bool) {
	http.SetCookie(c.Response(), &http.Cookie{
		Name: authcookie.AccessTokenName, Value: result.GetAccessToken(), Path: "/",
		HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode,
		MaxAge: int(result.GetAccessExpiresIn()),
	})
	http.SetCookie(c.Response(), &http.Cookie{
		Name: authcookie.RefreshTokenName, Value: result.GetRefreshToken(), Path: "/api/auth",
		HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode,
		MaxAge: int(result.GetRefreshExpiresIn()),
	})
}

func clearAuthCookies(c echo.Context, secure bool) {
	for _, cookie := range []*http.Cookie{
		{Name: authcookie.AccessTokenName, Path: "/"},
		{Name: authcookie.RefreshTokenName, Path: "/api/auth"},
	} {
		cookie.Value = ""
		cookie.HttpOnly = true
		cookie.Secure = secure
		cookie.SameSite = http.SameSiteLaxMode
		cookie.MaxAge = -1
		cookie.Expires = unixEpoch
		http.SetCookie(c.Response(), cookie)
	}
}

var unixEpoch = time.Unix(0, 0)

func (a *Apis) authError(c echo.Context, ctx context.Context, operation string, err error) error {
	st, ok := status.FromError(err)
	if !ok {
		log.Error(ctx, "non-grpc auth service error", "operation", operation, "error", err)
		return c.JSON(http.StatusBadGateway, view.ErrorResponse{Error: "auth service error"})
	}
	code := http.StatusBadGateway
	switch st.Code() {
	case codes.InvalidArgument:
		code = http.StatusBadRequest
	case codes.AlreadyExists:
		code = http.StatusConflict
	case codes.Unauthenticated:
		code = http.StatusUnauthorized
	case codes.NotFound:
		code = http.StatusNotFound
	case codes.ResourceExhausted:
		code = http.StatusTooManyRequests
	case codes.Unavailable:
		code = http.StatusServiceUnavailable
	case codes.DeadlineExceeded:
		code = http.StatusGatewayTimeout
	case codes.PermissionDenied:
		code = http.StatusForbidden
	case codes.FailedPrecondition:
		code = http.StatusBadRequest
	default:
		log.Error(ctx, "auth service request failed", "operation", operation, "grpc_code", st.Code().String(), "error", err)
	}
	message := st.Message()
	if code >= http.StatusInternalServerError {
		message = "auth service unavailable"
	}
	return c.JSON(code, view.ErrorResponse{Error: message})
}
