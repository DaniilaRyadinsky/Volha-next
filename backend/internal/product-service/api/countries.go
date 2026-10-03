package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *ServerAPI) CreateCountry(ctx context.Context, req *productsRPC.Country) (*emptypb.Empty, error) {
	const op = "grpc.CreateCountry"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.CreateCountry(ctx, domain.CountryFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) UpdateCountry(ctx context.Context, req *productsRPC.Country) (*emptypb.Empty, error) {
	const op = "grpc.UpdateCountry"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.UpdateCountry(ctx, domain.CountryFromRpc(req))
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) DeleteCountry(ctx context.Context, req *productsRPC.Id) (*emptypb.Empty, error) {
	const op = "grpc.DeleteCountry"

	_, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return nil, s.s.DeleteCountry(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *ServerAPI) GetAllCountries(ctx context.Context, _ *emptypb.Empty) (*productsRPC.CountryList, error) {
	const op = "grpc.GetAllCountries"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetAllCountries(ctx)
	})
	if err != nil {
		return nil, err
	}

	countries := res.([]*domain.Country)
	return &productsRPC.CountryList{
		Items: domain.CountriesToRpc(countries),
	}, nil
}

func (s *ServerAPI) GetCountry(ctx context.Context, req *productsRPC.Id) (*productsRPC.Country, error) {
	const op = "grpc.GetCountry"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetCountry(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}

	return domain.CountryToRpc(res.(*domain.Country)), nil
}
