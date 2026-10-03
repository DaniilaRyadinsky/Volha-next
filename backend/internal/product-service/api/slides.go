package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateSlide(ctx context.Context, req *productsRPC.Slide) (*emptypb.Empty, error) {
	const op = "grpc.CreateSlide"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateSlide(ctx, domain.SlideFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateSlide(ctx context.Context, req *productsRPC.Slide) (*emptypb.Empty, error) {
	const op = "grpc.UpdateSlide"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateSlide(ctx, domain.SlideFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteSlide(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteSlide"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteSlide(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllSlides(ctx context.Context, _ *emptypb.Empty) (*productsRPC.SlideList, error) {
	const op = "grpc.GetAllSlides"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllSlides(ctx)
	})
	if err != nil {
		return nil, err
	}

	slides := res.([]*domain.Slide)
	return &productsRPC.SlideList{
		Items: domain.SlidesToRpc(slides),
	}, nil
}

func (s *ServerAPI) GetSlide(ctx context.Context, req *productsRPC.Id) (*productsRPC.Slide, error) {
	const op = "grpc.GetSlide"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetSlide(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.SlideToRpc(res.(*domain.Slide)), nil
}
