package psqltx

import (
	"context"
	"database/sql"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
)

type RepoProvider struct {
	db *sql.DB
}

func NewRepoProvider(db *sql.DB) *RepoProvider {
	return &RepoProvider{db: db}
}

func (p *RepoProvider) Product(ctx context.Context) repository.ProductRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}

func (p *RepoProvider) Catalog(ctx context.Context) repository.CatalogRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}

func (p *RepoProvider) Dict(ctx context.Context) repository.DictRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}

func (p *RepoProvider) Pcp(ctx context.Context) repository.ProductColorPhotosRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}

func (p *RepoProvider) FileRefs(ctx context.Context) repository.FileReferenceRepo {
	if tx, ok := TxFromContext(ctx); ok {
		return repository.Driver{Driver: tx}
	}
	return repository.Driver{Driver: p.db}
}
