package view

type ProductCatalogItem struct {
	ID         string `json:"id"`
	Slug       string `json:"slug"`
	Category   string `json:"category"`
	Title      string `json:"title"`
	Price      int32  `json:"price"`
	Image      string `json:"image"`
	PriceSale  int32  `json:"pricesale"`
	IsFeatured bool   `json:"isFeatured"`
}

type ProductCatalogList struct {
	Items  []ProductCatalogItem `json:"items"`
	Total  int32                `json:"total"`
	Offset int32                `json:"offset"`
	Limit  int32                `json:"limit"`
}
