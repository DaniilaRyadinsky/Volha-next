package api

import (
	"context"
	"errors"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/auth-service/domain"
	"github.com/autumnterror/volha-backend/internal/auth-service/service"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func Register(server *grpc.Server, auth *service.Service) {
	productsRPC.RegisterAuthServiceServer(server, &serverAPI{auth: auth})
}

type serverAPI struct {
	productsRPC.UnimplementedAuthServiceServer
	auth *service.Service
}

func (s *serverAPI) Register(ctx context.Context, req *productsRPC.RegisterRequest) (*productsRPC.AuthResponse, error) {
	result, err := s.auth.Register(ctx, req.GetEmail(), req.GetLogin(), req.GetPassword())
	if err != nil {
		return nil, mapError(ctx, "Register", err)
	}
	return authResponseToRPC(result), nil
}

func (s *serverAPI) Login(ctx context.Context, req *productsRPC.LoginRequest) (*productsRPC.AuthResponse, error) {
	result, err := s.auth.Login(ctx, req.GetIdentifier(), req.GetPassword())
	if err != nil {
		return nil, mapError(ctx, "Login", err)
	}
	return authResponseToRPC(result), nil
}

func (s *serverAPI) Refresh(ctx context.Context, req *productsRPC.RefreshRequest) (*productsRPC.AuthResponse, error) {
	result, err := s.auth.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, mapError(ctx, "Refresh", err)
	}
	return authResponseToRPC(result), nil
}

func (s *serverAPI) Logout(ctx context.Context, req *productsRPC.LogoutRequest) (*emptypb.Empty, error) {
	if err := s.auth.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, mapError(ctx, "Logout", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) Authenticate(ctx context.Context, req *productsRPC.AuthenticateRequest) (*productsRPC.AuthenticatedUser, error) {
	principal, err := s.auth.Authenticate(ctx, req.GetAccessToken())
	if err != nil {
		return nil, mapError(ctx, "Authenticate", err)
	}
	return &productsRPC.AuthenticatedUser{User: userToRPC(principal.User), SessionId: principal.SessionID}, nil
}

func (s *serverAPI) ChangePassword(ctx context.Context, req *productsRPC.ChangePasswordRequest) (*productsRPC.AuthUser, error) {
	user, err := s.auth.ChangePassword(ctx, req.GetUserId(), req.GetSessionId(), req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		return nil, mapError(ctx, "ChangePassword", err)
	}
	return userToRPC(user), nil
}

func (s *serverAPI) ChangeLogin(ctx context.Context, req *productsRPC.ChangeLoginRequest) (*productsRPC.AuthUser, error) {
	user, err := s.auth.ChangeLogin(ctx, req.GetUserId(), req.GetSessionId(), req.GetPassword(), req.GetNewLogin())
	if err != nil {
		return nil, mapError(ctx, "ChangeLogin", err)
	}
	return userToRPC(user), nil
}

func (s *serverAPI) GetCart(ctx context.Context, req *productsRPC.CartRequest) (*productsRPC.Cart, error) {
	cart, err := s.auth.GetCart(ctx, req.GetUserId())
	if err != nil {
		return nil, mapError(ctx, "GetCart", err)
	}
	return cartToRPC(cart), nil
}

func (s *serverAPI) AddCartItem(ctx context.Context, req *productsRPC.CartItemRequest) (*productsRPC.Cart, error) {
	cart, err := s.auth.AddCartItem(ctx, req.GetUserId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		return nil, mapError(ctx, "AddCartItem", err)
	}
	return cartToRPC(cart), nil
}

func (s *serverAPI) UpdateCartItem(ctx context.Context, req *productsRPC.CartItemRequest) (*productsRPC.Cart, error) {
	cart, err := s.auth.UpdateCartItem(ctx, req.GetUserId(), req.GetProductId(), req.GetQuantity())
	if err != nil {
		return nil, mapError(ctx, "UpdateCartItem", err)
	}
	return cartToRPC(cart), nil
}

func (s *serverAPI) RemoveCartItem(ctx context.Context, req *productsRPC.CartItemKey) (*productsRPC.Cart, error) {
	cart, err := s.auth.RemoveCartItem(ctx, req.GetUserId(), req.GetProductId())
	if err != nil {
		return nil, mapError(ctx, "RemoveCartItem", err)
	}
	return cartToRPC(cart), nil
}

func (s *serverAPI) ClearCart(ctx context.Context, req *productsRPC.CartRequest) (*productsRPC.Cart, error) {
	cart, err := s.auth.ClearCart(ctx, req.GetUserId())
	if err != nil {
		return nil, mapError(ctx, "ClearCart", err)
	}
	return cartToRPC(cart), nil
}

func userToRPC(user *domain.User) *productsRPC.AuthUser {
	if user == nil {
		return nil
	}
	return &productsRPC.AuthUser{Id: user.ID, Email: user.Email, Login: user.Login}
}

func authResponseToRPC(result *service.TokenPair) *productsRPC.AuthResponse {
	return &productsRPC.AuthResponse{
		User:             userToRPC(result.User),
		AccessToken:      result.AccessToken,
		RefreshToken:     result.RefreshToken,
		AccessExpiresIn:  result.AccessExpiresIn,
		RefreshExpiresIn: result.RefreshExpiresIn,
	}
}

func cartToRPC(cart *domain.Cart) *productsRPC.Cart {
	response := &productsRPC.Cart{TotalQuantity: cart.TotalQuantity, TotalPrice: cart.TotalPrice}
	for _, item := range cart.Items {
		response.Items = append(response.Items, &productsRPC.CartItem{
			ProductId: item.ProductID,
			Slug:      item.Slug,
			Title:     item.Title,
			Quantity:  item.Quantity,
			Price:     item.Price,
			Pricesale: item.PriceSale,
			Image:     item.Image,
			LineTotal: item.LineTotal,
		})
	}
	return response
}

func mapError(ctx context.Context, operation string, err error) error {
	switch {
	case errors.Is(err, domain.ErrInvalidInput):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, domain.ErrAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrUnauthenticated):
		return status.Error(codes.Unauthenticated, err.Error())
	case errors.Is(err, domain.ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrCartLimit):
		return status.Error(codes.ResourceExhausted, err.Error())
	default:
		log.Error(ctx, "auth grpc request failed", "operation", operation, "error", err)
		return status.Error(codes.Internal, "check logs")
	}
}
