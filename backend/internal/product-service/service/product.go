package service

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

func (s *ProductsService) GetAllProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int, error) {
	const op = "service.GetAllProducts"

	if err := validateProductRange(offset, limit); err != nil {
		return nil, 0, wrapServiceCheck(op, err)
	}

	repo, err := s.productRepo(ctx)
	if err != nil {
		return nil, 0, wrapServiceCheck(op, err)
	}
	return repo.GetAllProducts(ctx, offset, limit)
}

func (s *ProductsService) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	const op = "service.GetProduct"

	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	repo, err := s.productRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return repo.GetProduct(ctx, id)
}

func (s *ProductsService) CreateProduct(ctx context.Context, p *domain.ProductId) error {
	const op = "service.CreateProduct"

	if err := validateProductPayload(p); err != nil {
		return wrapServiceCheck(op, err)
	}
	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.productRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateProduct(ctx, p)
	})
}

func (s *ProductsService) UpdateProduct(ctx context.Context, p *domain.ProductId) error {
	const op = "service.UpdateProduct"

	if err := validateProductPayload(p); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.productRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateProduct(ctx, p)
	})
}

func (s *ProductsService) DeleteProduct(ctx context.Context, id string) error {
	const op = "service.DeleteProduct"

	if err := validateID(id); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.productRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteProduct(ctx, id)
	})
}

func (s *ProductsService) FilterProducts(ctx context.Context, filter *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error) {
	const op = "service.FilterProducts"

	if err := validateFilter(filter); err != nil {
		return nil, 0, wrapServiceCheck(op, err)
	}

	repo, err := s.productRepo(ctx)
	if err != nil {
		return nil, 0, wrapServiceCheck(op, err)
	}
	return repo.FilterProducts(ctx, filter)
}
