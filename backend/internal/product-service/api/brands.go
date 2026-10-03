package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateBrand(ctx context.Context, req *productsRPC.Brand) (*emptypb.Empty, error) {
	const op = "grpc.CreateBrand"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateBrand(ctx, domain.BrandFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateBrand(ctx context.Context, req *productsRPC.Brand) (*emptypb.Empty, error) {
	const op = "grpc.UpdateBrand"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateBrand(ctx, domain.BrandFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteBrand(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteBrand"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteBrand(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllBrands(ctx context.Context, _ *emptypb.Empty) (*productsRPC.BrandList, error) {
	const op = "grpc.GetAllBrands"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllBrands(ctx)
	})
	if err != nil {
		return nil, err
	}

	brands := res.([]*domain.Brand)
	return &productsRPC.BrandList{
		Items: domain.BrandsToRpc(brands),
	}, nil
}

func (s *ServerAPI) GetBrand(ctx context.Context, req *productsRPC.Id) (*productsRPC.Brand, error) {
	const op = "grpc.GetBrand"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetBrand(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.BrandToRpc(res.(*domain.Brand)), nil
}
