package service

import (
	"context"
	"errors"

	"github.com/autumnterror/volha-backend/internal/product-service/repository"
)

func (s *ProductsService) runInTx(ctx context.Context, op string, fn func(ctx context.Context) error) error {
	if s.tx == nil {
		return wrapServiceCheck(op, errors.New("tx runner is nil"))
	}
	return s.tx.RunInTx(ctx, fn)
}

func (s *ProductsService) productRepo(ctx context.Context) (repository.ProductRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.Product(ctx)
	if repo == nil {
		return nil, errors.New("product repository is nil")
	}
	return repo, nil
}

func (s *ProductsService) catalogRepo(ctx context.Context) (repository.CatalogRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.Catalog(ctx)
	if repo == nil {
		return nil, errors.New("catalog repository is nil")
	}
	return repo, nil
}

func (s *ProductsService) dictRepo(ctx context.Context) (repository.DictRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.Dict(ctx)
	if repo == nil {
		return nil, errors.New("dictionary repository is nil")
	}
	return repo, nil
}

func (s *ProductsService) pcpRepo(ctx context.Context) (repository.ProductColorPhotosRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.Pcp(ctx)
	if repo == nil {
		return nil, errors.New("product color photos repository is nil")
	}
	return repo, nil
}

func (s *ProductsService) fileReferenceRepo(ctx context.Context) (repository.FileReferenceRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.FileRefs(ctx)
	if repo == nil {
		return nil, errors.New("file reference repository is nil")
	}
	return repo, nil
}
