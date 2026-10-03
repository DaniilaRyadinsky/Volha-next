package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
)

func (s *ServerAPI) IsFileReferenced(ctx context.Context, req *productsRPC.FileReferenceRequest) (*productsRPC.FileReferenceResponse, error) {
	const op = "grpc.IsFileReferenced"

	res, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.IsFileReferenced(ctx, req.GetFilename())
	})
	if err != nil {
		return nil, err
	}

	return &productsRPC.FileReferenceResponse{
		Referenced: res.(bool),
	}, nil
}
