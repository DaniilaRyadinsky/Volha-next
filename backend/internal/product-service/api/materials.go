package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateMaterial(ctx context.Context, req *productsRPC.Material) (*emptypb.Empty, error) {
	const op = "grpc.CreateMaterial"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateMaterial(ctx, domain.MaterialFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateMaterial(ctx context.Context, req *productsRPC.Material) (*emptypb.Empty, error) {
	const op = "grpc.UpdateMaterial"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateMaterial(ctx, domain.MaterialFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteMaterial(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteMaterial"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteMaterial(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllMaterials(ctx context.Context, _ *emptypb.Empty) (*productsRPC.MaterialList, error) {
	const op = "grpc.GetAllMaterials"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllMaterials(ctx)
	})
	if err != nil {
		return nil, err
	}

	materials := res.([]*domain.Material)
	return &productsRPC.MaterialList{
		Items: domain.MaterialsToRpc(materials),
	}, nil
}

func (s *ServerAPI) GetMaterial(ctx context.Context, req *productsRPC.Id) (*productsRPC.Material, error) {
	const op = "grpc.GetMaterial"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetMaterial(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.MaterialToRpc(res.(*domain.Material)), nil
}
