package view

import (
	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

type ProductCatalogItem struct {
	ID       string           `json:"id"`
	Slug     string           `json:"slug"`
	Title    string           `json:"title"`
	Category *domain.Category `json:"category"`
	Photos   []string         `json:"photos"`
	Price    int32            `json:"price"`
}
type ProductCatalogList struct {
	Items []ProductCatalogItem `json:"items"`
	Total int32                `json:"total"`
	Page  int32                `json:"page"`
	Limit int32                `json:"limit"`
}

func CatalogItemFromRPC(p *productsRPC.ProductCatalogItem) ProductCatalogItem {
	return ProductCatalogItem{ID: p.GetId(), Slug: p.GetSlug(), Title: p.GetTitle(), Category: domain.CategoryFromRpc(p.GetCategoryDetails()), Photos: append([]string{}, p.GetPhotos()...), Price: p.GetPrice()}
}
func CatalogItemFromProduct(p *productsRPC.Product) ProductCatalogItem {
	return ProductCatalogItem{ID: p.GetId(), Slug: p.GetSlug(), Title: p.GetTitle(), Category: domain.CategoryFromRpc(p.GetCategory()), Photos: append([]string{}, p.GetPhotos()...), Price: p.GetPrice()}
}
