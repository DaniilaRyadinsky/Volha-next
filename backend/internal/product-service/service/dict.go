package service

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

func (s *ProductsService) GetDictionaries(ctx context.Context, idCat string) (*domain.Dictionaries, error) {
	const op = "service.GetDictionaries"
	repo, err := s.dictRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if idCat == domain.NotByCategory || idCat == "" {
		return repo.GetDictionaries(ctx)
	}

	if err := validateID(idCat); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	catalog, err := s.catalogRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if _, err := catalog.GetCategory(ctx, idCat); err != nil {
		return nil, err
	}
	return repo.GetDictionariesByCategoryID(ctx, idCat)
}
