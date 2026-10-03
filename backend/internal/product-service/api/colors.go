package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateColor(ctx context.Context, req *productsRPC.Color) (*emptypb.Empty, error) {
	const op = "grpc.CreateColor"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateColor(ctx, domain.ColorFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateColor(ctx context.Context, req *productsRPC.Color) (*emptypb.Empty, error) {
	const op = "grpc.UpdateColor"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateColor(ctx, domain.ColorFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteColor(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteColor"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteColor(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllColors(ctx context.Context, _ *emptypb.Empty) (*productsRPC.ColorList, error) {
	const op = "grpc.GetAllColors"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllColors(ctx)
	})
	if err != nil {
		return nil, err
	}

	colors := res.([]*domain.Color)
	return &productsRPC.ColorList{
		Items: domain.ColorsToRpc(colors),
	}, nil
}

func (s *ServerAPI) GetColor(ctx context.Context, req *productsRPC.Id) (*productsRPC.Color, error) {
	const op = "grpc.GetColor"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetColor(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.ColorToRpc(res.(*domain.Color)), nil
}
