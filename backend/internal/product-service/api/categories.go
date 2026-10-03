package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateCategory(ctx context.Context, req *productsRPC.Category) (*emptypb.Empty, error) {
	const op = "grpc.CreateCategory"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateCategory(ctx, domain.CategoryFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateCategory(ctx context.Context, req *productsRPC.Category) (*emptypb.Empty, error) {
	const op = "grpc.UpdateCategory"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateCategory(ctx, domain.CategoryFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteCategory(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteCategory"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteCategory(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllCategories(ctx context.Context, _ *emptypb.Empty) (*productsRPC.CategoryList, error) {
	const op = "grpc.GetAllCategories"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllCategories(ctx)
	})
	if err != nil {
		return nil, err
	}

	categories := res.([]*domain.Category)
	return &productsRPC.CategoryList{
		Items: domain.CategoriesToRpc(categories),
	}, nil
}

func (s *ServerAPI) GetCategory(ctx context.Context, req *productsRPC.Id) (*productsRPC.Category, error) {
	const op = "grpc.GetCategory"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetCategory(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.CategoryToRpc(res.(*domain.Category)), nil
}
