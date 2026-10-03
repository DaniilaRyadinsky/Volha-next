package api

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

func (s *ServerAPI) GetDictionaries(ctx context.Context, req *productsRPC.Id) (*productsRPC.Dictionaries, error) {
	const op = "grpc.GetDictionaries"

	d, err := handleCRUDResponse(ctx, op, func() (any, error) {
		return s.s.GetDictionaries(ctx, req.GetId())
	})
	if err != nil {
		return nil, err
	}
	return domain.DictionariesToRpc(d.(*domain.Dictionaries)), nil
}
