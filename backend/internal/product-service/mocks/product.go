package mocks

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
)

type TxRunner struct{}

func (TxRunner) RunInTx(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

type Provider struct {
	ArticleRepo repository.ArticleRepo
	ProductRepo repository.ProductRepo
	CatalogRepo repository.CatalogRepo
	DictRepo    repository.DictRepo
	PcpRepo     repository.ProductColorPhotosRepo
	FileRepo    repository.FileReferenceRepo
}

func (p Provider) Product(context.Context) repository.ProductRepo {
	return p.ProductRepo
}

func (p Provider) Catalog(context.Context) repository.CatalogRepo {
	return p.CatalogRepo
}

func (p Provider) Dict(context.Context) repository.DictRepo {
	return p.DictRepo
}

func (p Provider) Pcp(context.Context) repository.ProductColorPhotosRepo {
	return p.PcpRepo
}

func (p Provider) FileRefs(context.Context) repository.FileReferenceRepo {
	return p.FileRepo
}

type ProductRepo struct {
	GetAllProductsFn func(context.Context, int, int) ([]*domain.Product, int, error)
	GetProductFn     func(context.Context, string) (*domain.Product, error)
	CreateProductFn  func(context.Context, *domain.ProductId) error
	UpdateProductFn  func(context.Context, *domain.ProductId) error
	DeleteProductFn  func(context.Context, string) error
	FilterProductsFn func(context.Context, *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error)
}

func (r ProductRepo) GetAllProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int, error) {
	return r.GetAllProductsFn(ctx, offset, limit)
}

func (r ProductRepo) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	return r.GetProductFn(ctx, id)
}

func (r ProductRepo) CreateProduct(ctx context.Context, p *domain.ProductId) error {
	return r.CreateProductFn(ctx, p)
}

func (r ProductRepo) UpdateProduct(ctx context.Context, p *domain.ProductId) error {
	return r.UpdateProductFn(ctx, p)
}

func (r ProductRepo) DeleteProduct(ctx context.Context, id string) error {
	return r.DeleteProductFn(ctx, id)
}

func (r ProductRepo) FilterProducts(ctx context.Context, filter *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error) {
	return r.FilterProductsFn(ctx, filter)
}

func (p Provider) Article(context.Context) repository.ArticleRepo { return p.ArticleRepo }
