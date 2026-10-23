package repository

import "context"

type Provider interface {
	Article(context.Context) ArticleRepo
	Product(ctx context.Context) ProductRepo
	Catalog(ctx context.Context) CatalogRepo
	Dict(ctx context.Context) DictRepo
	Pcp(ctx context.Context) ProductColorPhotosRepo
	FileRefs(ctx context.Context) FileReferenceRepo
}
