package service

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

func (s *ProductsService) GetPhotosByProductColor(ctx context.Context, productID, colorID string) ([]string, error) {
	const op = "service.GetPhotosByProductColor"
	repo, err := s.pcpRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	if err := validatePCPIDs(productID, colorID); err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return repo.GetPhotosByProductColor(ctx, productID, colorID)
}

func (s *ProductsService) GetAllProductColorPhotos(ctx context.Context) ([]domain.ProductColorPhotos, error) {
	const op = "service.GetAllProductColorPhotos"
	repo, err := s.pcpRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return repo.GetAllProductColorPhotos(ctx)
}

func (s *ProductsService) GetProductColorPhotosByProduct(ctx context.Context, productID string) ([]domain.ProductColorPhotos, error) {
	const op = "service.GetProductColorPhotosByProduct"
	if err := validateID(productID); err != nil {
		return nil, wrapServiceCheck(op, err)
	}
	repo, err := s.pcpRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck(op, err)
	}

	return repo.GetProductColorPhotosByProduct(ctx, productID)
}

func (s *ProductsService) CreateProductColorPhotos(ctx context.Context, pcp *domain.ProductColorPhotos) error {
	const op = "service.CreateProductColorPhotos"
	if err := validatePCP(pcp); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.pcpRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.CreateProductColorPhotos(ctx, pcp)
	})
}

func (s *ProductsService) UpdateProductColorPhotos(ctx context.Context, pcp *domain.ProductColorPhotos) error {
	const op = "service.UpdateProductColorPhotos"
	if err := validatePCP(pcp); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.pcpRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.UpdateProductColorPhotos(ctx, pcp)
	})
}

func (s *ProductsService) DeleteProductColorPhotos(ctx context.Context, productID, colorID string) error {
	const op = "service.DeleteProductColorPhotos"
	if err := validatePCPIDs(productID, colorID); err != nil {
		return wrapServiceCheck(op, err)
	}

	return s.runInTx(ctx, op, func(ctx context.Context) error {
		repo, err := s.pcpRepo(ctx)
		if err != nil {
			return wrapServiceCheck(op, err)
		}
		return repo.DeleteProductColorPhotos(ctx, productID, colorID)
	})
}
