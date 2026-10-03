package repository

import "context"

type Provider interface {
	Product(ctx context.Context) ProductRepo
	Catalog(ctx context.Context) CatalogRepo
	Dict(ctx context.Context) DictRepo
	Pcp(ctx context.Context) ProductColorPhotosRepo
	FileRefs(ctx context.Context) FileReferenceRepo
}
