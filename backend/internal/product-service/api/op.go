package api

import (
	"context"
	"errors"

	"github.com/autumnterror/volha-backend/api/proto/gen"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerAPI struct {
	productsRPC.UnimplementedProductsServer
	s *service.ProductsService
}

func Register(
	server *grpc.Server,
	s *service.ProductsService,
) {
	productsRPC.RegisterProductsServer(server, &ServerAPI{
		s: s,
	})
}

func handleCRUDResponse(ctx context.Context, op string, action func() (any, error)) (any, error) {
	res, err := action()
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrUnknownType):
			return nil, status.Error(codes.Unimplemented, err.Error())
		case errors.Is(err, domain.ErrInvalidType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, domain.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, domain.ErrAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, domain.ErrForeignKey):
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		case errors.Is(err, service.ErrBadServiceCheck):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Error(ctx, "grpc action failed", "operation", op, "error", err)
			return nil, status.Error(codes.Internal, "check logs")
		}
	}
	if res != nil {
		return res, nil
	}
	return nil, nil
}
