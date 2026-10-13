package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/autumnterror/volha-backend/internal/auth-service/repository"
	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

const tokenIssuer = "volha-auth-service"

var loginPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{2,31}$`)

type Service struct {
	store      *repository.Store
	secret     []byte
	accessTTL  time.Duration
	refreshTTL time.Duration
}

type TokenPair struct {
	User             *domain.User
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int64
	RefreshExpiresIn int64
}

type Principal struct {
	User      *domain.User
	SessionID string
}

type claims struct {
	Subject   string `json:"sub"`
	SessionID string `json:"sid"`
	JWTID     string `json:"jti"`
	TokenUse  string `json:"token_use"`
	Issuer    string `json:"iss"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
}

func New(store *repository.Store, secret string, accessTTL, refreshTTL time.Duration) *Service {
	return &Service{
		store:      store,
		secret:     []byte(secret),
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (s *Service) Register(ctx context.Context, email, login, password string) (*TokenPair, error) {
	normalizedEmail, err := normalizeEmail(email)
	if err != nil {
		return nil, err
	}
	normalizedLogin, err := normalizeLogin(login)
	if err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{ID: xid.New().String(), Email: normalizedEmail, Login: normalizedLogin}
	sessionID := xid.New().String()
	pair, expiresAt, err := s.issuePair(user.ID, sessionID)
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateUserAndSession(ctx, user, string(passwordHash), sessionID, hashToken(pair.RefreshToken), expiresAt); err != nil {
		return nil, err
	}
	pair.User = user
	return pair, nil
}

func (s *Service) Login(ctx context.Context, identifier, password string) (*TokenPair, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" || password == "" {
		return nil, domain.ErrUnauthenticated
	}
	user, passwordHash, err := s.store.FindUserForLogin(ctx, identifier)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, domain.ErrUnauthenticated
	}
	sessionID := xid.New().String()
	pair, expiresAt, err := s.issuePair(user.ID, sessionID)
	if err != nil {
		return nil, err
	}
	if err := s.store.CreateSession(ctx, user.ID, sessionID, hashToken(pair.RefreshToken), expiresAt); err != nil {
		return nil, err
	}
	pair.User = user
	return pair, nil
}

func (s *Service) Refresh(ctx context.Context, rawToken string) (*TokenPair, error) {
	tokenClaims, err := s.verifyToken(rawToken, "refresh", false)
	if err != nil {
		return nil, domain.ErrUnauthenticated
	}
	pair, expiresAt, err := s.issuePair(tokenClaims.Subject, tokenClaims.SessionID)
	if err != nil {
		return nil, err
	}
	user, err := s.store.RotateSession(ctx, tokenClaims.Subject, tokenClaims.SessionID, hashToken(rawToken), hashToken(pair.RefreshToken), expiresAt)
	if err != nil {
		return nil, err
	}
	pair.User = user
	return pair, nil
}

func (s *Service) Logout(ctx context.Context, rawToken string) error {
	tokenClaims, err := s.verifyToken(rawToken, "refresh", true)
	if err != nil {
		return nil
	}
	return s.store.RevokeSession(ctx, tokenClaims.Subject, tokenClaims.SessionID, hashToken(rawToken))
}

func (s *Service) Authenticate(ctx context.Context, rawToken string) (*Principal, error) {
	tokenClaims, err := s.verifyToken(rawToken, "access", false)
	if err != nil {
		return nil, domain.ErrUnauthenticated
	}
	user, err := s.store.AuthenticateSession(ctx, tokenClaims.Subject, tokenClaims.SessionID)
	if err != nil {
		return nil, err
	}
	return &Principal{User: user, SessionID: tokenClaims.SessionID}, nil
}

func (s *Service) ChangePassword(ctx context.Context, userID, sessionID, oldPassword, newPassword string) (*domain.User, error) {
	if err := validatePassword(newPassword); err != nil {
		return nil, err
	}
	if oldPassword == "" || oldPassword == newPassword {
		return nil, domain.ErrInvalidInput
	}
	passwordHash, err := s.store.PasswordHash(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(oldPassword)); err != nil {
		return nil, domain.ErrUnauthenticated
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(newPassword)); err == nil {
		return nil, domain.ErrInvalidInput
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	if err := s.store.ChangePassword(ctx, userID, sessionID, string(newHash)); err != nil {
		return nil, err
	}
	return s.store.AuthenticateSession(ctx, userID, sessionID)
}

func (s *Service) ChangeLogin(ctx context.Context, userID, sessionID, password, newLogin string) (*domain.User, error) {
	login, err := normalizeLogin(newLogin)
	if err != nil {
		return nil, err
	}
	passwordHash, err := s.store.PasswordHash(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		return nil, domain.ErrUnauthenticated
	}
	return s.store.ChangeLogin(ctx, userID, sessionID, login)
}

func (s *Service) GetCart(ctx context.Context, userID string) (*domain.Cart, error) {
	return s.store.GetCart(ctx, userID)
}

func (s *Service) AddCartItem(ctx context.Context, userID, productID string, quantity int32) (*domain.Cart, error) {
	if err := validateProductID(productID); err != nil {
		return nil, err
	}
	if quantity < 1 || quantity > 99 {
		return nil, domain.ErrInvalidInput
	}
	if err := s.store.AddCartItem(ctx, userID, productID, quantity); err != nil {
		return nil, err
	}
	return s.store.GetCart(ctx, userID)
}

func (s *Service) UpdateCartItem(ctx context.Context, userID, productID string, quantity int32) (*domain.Cart, error) {
	if err := validateProductID(productID); err != nil {
		return nil, err
	}
	if quantity < 1 || quantity > 99 {
		return nil, domain.ErrInvalidInput
	}
	if err := s.store.UpdateCartItem(ctx, userID, productID, quantity); err != nil {
		return nil, err
	}
	return s.store.GetCart(ctx, userID)
}

func (s *Service) RemoveCartItem(ctx context.Context, userID, productID string) (*domain.Cart, error) {
	if err := validateProductID(productID); err != nil {
		return nil, err
	}
	if err := s.store.RemoveCartItem(ctx, userID, productID); err != nil {
		return nil, err
	}
	return s.store.GetCart(ctx, userID)
}

func (s *Service) ClearCart(ctx context.Context, userID string) (*domain.Cart, error) {
	if err := s.store.ClearCart(ctx, userID); err != nil {
		return nil, err
	}
	return s.store.GetCart(ctx, userID)
}

func (s *Service) issuePair(userID, sessionID string) (*TokenPair, time.Time, error) {
	now := time.Now().UTC()
	accessExpiresAt := now.Add(s.accessTTL)
	refreshExpiresAt := now.Add(s.refreshTTL)
	accessToken, err := s.signToken(userID, sessionID, "access", now, accessExpiresAt)
	if err != nil {
		return nil, time.Time{}, err
	}
	refreshToken, err := s.signToken(userID, sessionID, "refresh", now, refreshExpiresAt)
	if err != nil {
		return nil, time.Time{}, err
	}
	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		AccessExpiresIn:  int64(s.accessTTL / time.Second),
		RefreshExpiresIn: int64(s.refreshTTL / time.Second),
	}, refreshExpiresAt, nil
}

func (s *Service) signToken(userID, sessionID, tokenUse string, issuedAt, expiresAt time.Time) (string, error) {
	headerJSON, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(claims{
		Subject:   userID,
		SessionID: sessionID,
		JWTID:     xid.New().String(),
		TokenUse:  tokenUse,
		Issuer:    tokenIssuer,
		IssuedAt:  issuedAt.Unix(),
		ExpiresAt: expiresAt.Unix(),
	})
	if err != nil {
		return "", err
	}
	encodedHeader := base64.RawURLEncoding.EncodeToString(headerJSON)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	message := encodedHeader + "." + encodedPayload
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(message))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return message + "." + signature, nil
}

func (s *Service) verifyToken(rawToken, expectedUse string, allowExpired bool) (*claims, error) {
	parts := strings.Split(rawToken, ".")
	if len(parts) != 3 || parts[0] == "" || parts[1] == "" || parts[2] == "" {
		return nil, errors.New("malformed token")
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("malformed token header")
	}
	var header struct {
		Algorithm string `json:"alg"`
		Type      string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Algorithm != "HS256" || header.Type != "JWT" {
		return nil, errors.New("unsupported token header")
	}
	providedSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("malformed token signature")
	}
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(parts[0] + "." + parts[1]))
	if !hmac.Equal(providedSignature, mac.Sum(nil)) {
		return nil, errors.New("invalid token signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("malformed token payload")
	}
	var tokenClaims claims
	if err := json.Unmarshal(payload, &tokenClaims); err != nil {
		return nil, errors.New("malformed token claims")
	}
	now := time.Now().Unix()
	if tokenClaims.Issuer != tokenIssuer || tokenClaims.TokenUse != expectedUse || tokenClaims.Subject == "" || tokenClaims.SessionID == "" || tokenClaims.JWTID == "" || tokenClaims.ExpiresAt <= tokenClaims.IssuedAt || tokenClaims.IssuedAt > now+60 {
		return nil, errors.New("invalid token claims")
	}
	if !allowExpired && tokenClaims.ExpiresAt <= now {
		return nil, errors.New("expired token")
	}
	return &tokenClaims, nil
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func normalizeEmail(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 254 {
		return "", fmt.Errorf("%w: email is too long", domain.ErrInvalidInput)
	}
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address != value || !strings.Contains(value, "@") {
		return "", fmt.Errorf("%w: email is invalid", domain.ErrInvalidInput)
	}
	return strings.ToLower(value), nil
}

func normalizeLogin(value string) (string, error) {
	value = strings.ToLower(strings.TrimSpace(value))
	if !loginPattern.MatchString(value) {
		return "", fmt.Errorf("%w: login must contain 3-32 latin letters, digits, dots, underscores or hyphens", domain.ErrInvalidInput)
	}
	return value, nil
}

func validatePassword(password string) error {
	if !utf8.ValidString(password) || len(password) < 10 || len(password) > 72 {
		return fmt.Errorf("%w: password must be 10-72 bytes", domain.ErrInvalidInput)
	}
	hasLetter, hasDigit := false, false
	for _, r := range password {
		hasLetter = hasLetter || unicode.IsLetter(r)
		hasDigit = hasDigit || unicode.IsDigit(r)
	}
	if !hasLetter || !hasDigit {
		return fmt.Errorf("%w: password must contain a letter and a number", domain.ErrInvalidInput)
	}
	return nil
}

func validateProductID(productID string) error {
	if _, err := xid.FromString(productID); err != nil {
		return fmt.Errorf("%w: product id is invalid", domain.ErrInvalidInput)
	}
	return nil
}
