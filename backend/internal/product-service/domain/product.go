package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Product struct {
	Id          string      `json:"id"`
	Title       string      `json:"title"`
	Slug        string      `json:"slug"`
	Article     string      `json:"article"`
	Brand       *Brand      `json:"brand"`
	Category    *Category   `json:"category"`
	Country     *Country    `json:"country"`
	Width       int32       `json:"width"`
	Height      int32       `json:"height"`
	Depth       int32       `json:"depth"`
	Materials   []*Material `json:"materials"`
	Colors      []*Color    `json:"colors"`
	Photos      []string    `json:"photos"`
	Seems       []*Product  `json:"seems"`
	Price       int32       `json:"price"`
	Description string      `json:"description"`
	IsFeatured  bool        `json:"isFeatured"`
	PriceSale   int32       `json:"pricesale"`
}

type ProductId struct {
	Id          string   `json:"id"`
	Title       string   `json:"title"`
	Slug        string   `json:"slug"`
	Article     string   `json:"article"`
	Brand       string   `json:"brand"`
	Category    string   `json:"category"`
	Country     string   `json:"country"`
	Width       int32    `json:"width"`
	Height      int32    `json:"height"`
	Depth       int32    `json:"depth"`
	Materials   []string `json:"materials"`
	Colors      []string `json:"colors"`
	Photos      []string `json:"photos"`
	Seems       []string `json:"seems"`
	Price       int32    `json:"price"`
	Description string   `json:"description"`
	IsFeatured  bool     `json:"isFeatured"`
	PriceSale   int32    `json:"pricesale"`
}

type ProductCatalogItem struct {
	Id         string `json:"id"`
	Slug       string `json:"slug"`
	Category   string `json:"category"`
	Title      string `json:"title"`
	Price      int32  `json:"price"`
	Image      string `json:"image"`
	PriceSale  int32  `json:"pricesale"`
	IsFeatured bool   `json:"isFeatured"`
}

func ProductCatalogItemsToRpc(items []*ProductCatalogItem) []*productsRPC.ProductCatalogItem {
	res := make([]*productsRPC.ProductCatalogItem, 0, len(items))
	for _, item := range items {
		if item == nil {
			continue
		}
		res = append(res, &productsRPC.ProductCatalogItem{
			Id:         item.Id,
			Slug:       item.Slug,
			Category:   item.Category,
			Title:      item.Title,
			Price:      item.Price,
			Image:      item.Image,
			Pricesale:  item.PriceSale,
			IsFeatured: item.IsFeatured,
		})
	}
	return res
}

func NewEmptyProduct() Product {
	return Product{
		Id:      "",
		Title:   "",
		Slug:    "",
		Article: "",
		Brand: &Brand{
			Id:    "",
			Title: "",
		},
		Category: &Category{
			Id:    "",
			Title: "",
			Slug:  "",
			Img:   "",
		},
		Country: &Country{
			Id:       "",
			Title:    "",
			Friendly: "",
		},
		Width:       0,
		Height:      0,
		Depth:       0,
		Materials:   []*Material{},
		Colors:      []*Color{},
		Photos:      []string{},
		Seems:       []*Product{},
		Price:       0,
		Description: "",
		IsFeatured:  false,
		PriceSale:   0,
	}
}

type ProductFilter struct {
	Brand     []string `json:"brand"`
	Country   []string `json:"country"`
	Category  []string `json:"category"`
	MinWidth  int32    `json:"min_width"`
	MaxWidth  int32    `json:"max_width"`
	MinHeight int32    `json:"min_height"`
	MaxHeight int32    `json:"max_height"`
	MinDepth  int32    `json:"min_depth"`
	MaxDepth  int32    `json:"max_depth"`
	Materials []string `json:"materials"`
	Colors    []string `json:"colors"`
	MinPrice  int32    `json:"min_price"`
	MaxPrice  int32    `json:"max_price"`
	SortBy    string   `json:"sort_by"`
	SortOrder string   `json:"sort_order"`
	Offset    int32    `json:"offset"`
	Limit     int32    `json:"limit"`
	Title     string   `json:"title"`
}

func ProductFromRpc(p *productsRPC.Product) *Product {
	if p == nil {
		return nil
	}

	materials := make([]*Material, 0, len(p.GetMaterials()))
	for _, m := range p.GetMaterials() {
		if val := MaterialFromRpc(m); val != nil {
			materials = append(materials, val)
		}
	}

	colors := make([]*Color, 0, len(p.GetColors()))
	for _, c := range p.GetColors() {
		if val := ColorFromRpc(c); val != nil {
			colors = append(colors, val)
		}
	}

	seems := make([]*Product, 0, len(p.GetSeems()))
	for _, s := range p.GetSeems() {
		if val := ProductFromRpc(s); val != nil {
			seems = append(seems, val)
		}
	}

	return &Product{
		Id:          p.GetId(),
		Title:       p.GetTitle(),
		Slug:        p.GetSlug(),
		Article:     p.GetArticle(),
		Brand:       BrandFromRpc(p.GetBrand()),
		Category:    CategoryFromRpc(p.GetCategory()),
		Country:     CountryFromRpc(p.GetCountry()),
		Width:       p.GetWidth(),
		Height:      p.GetHeight(),
		Depth:       p.GetDepth(),
		Materials:   materials,
		Colors:      colors,
		Photos:      append([]string{}, p.GetPhotos()...),
		Seems:       seems,
		Price:       p.GetPrice(),
		Description: p.GetDescription(),
		IsFeatured:  p.GetIsFeatured(),
		PriceSale:   p.GetPricesale(),
	}
}

func ProductToRpc(p *Product) *productsRPC.Product {
	if p == nil {
		return nil
	}

	materials := make([]*productsRPC.Material, 0, len(p.Materials))
	for _, m := range p.Materials {
		if val := MaterialToRpc(m); val != nil {
			materials = append(materials, val)
		}
	}

	colors := make([]*productsRPC.Color, 0, len(p.Colors))
	for _, c := range p.Colors {
		if val := ColorToRpc(c); val != nil {
			colors = append(colors, val)
		}
	}

	seems := make([]*productsRPC.Product, 0, len(p.Seems))
	for _, s := range p.Seems {
		if val := ProductToRpc(s); val != nil {
			seems = append(seems, val)
		}
	}

	return &productsRPC.Product{
		Id:          p.Id,
		Title:       p.Title,
		Slug:        p.Slug,
		Article:     p.Article,
		Brand:       BrandToRpc(p.Brand),
		Category:    CategoryToRpc(p.Category),
		Country:     CountryToRpc(p.Country),
		Width:       p.Width,
		Height:      p.Height,
		Depth:       p.Depth,
		Materials:   materials,
		Colors:      colors,
		Photos:      append([]string{}, p.Photos...),
		Seems:       seems,
		Price:       p.Price,
		Description: p.Description,
		IsFeatured:  p.IsFeatured,
		Pricesale:   p.PriceSale,
	}
}

func ProductsFromRpc(items []*productsRPC.Product) []*Product {
	res := make([]*Product, 0, len(items))
	for _, p := range items {
		if val := ProductFromRpc(p); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func ProductsToRpc(items []*Product) []*productsRPC.Product {
	res := make([]*productsRPC.Product, 0, len(items))
	for _, p := range items {
		if val := ProductToRpc(p); val != nil {
			res = append(res, val)
		}
	}
	return res
}

func ProductIdFromRpc(p *productsRPC.ProductId) *ProductId {
	if p == nil {
		return nil
	}
	return &ProductId{
		Id:          p.GetId(),
		Title:       p.GetTitle(),
		Slug:        p.GetSlug(),
		Article:     p.GetArticle(),
		Brand:       p.GetBrand(),
		Category:    p.GetCategory(),
		Country:     p.GetCountry(),
		Width:       p.GetWidth(),
		Height:      p.GetHeight(),
		Depth:       p.GetDepth(),
		Materials:   append([]string{}, p.GetMaterials()...),
		Colors:      append([]string{}, p.GetColors()...),
		Photos:      append([]string{}, p.GetPhotos()...),
		Seems:       append([]string{}, p.GetSeems()...),
		Price:       p.GetPrice(),
		Description: p.GetDescription(),
		IsFeatured:  p.GetIsFeatured(),
		PriceSale:   p.GetPricesale(),
	}
}

func ProductIdToRpc(p *ProductId) *productsRPC.ProductId {
	if p == nil {
		return nil
	}
	return &productsRPC.ProductId{
		Id:          p.Id,
		Title:       p.Title,
		Slug:        p.Slug,
		Article:     p.Article,
		Brand:       p.Brand,
		Category:    p.Category,
		Country:     p.Country,
		Width:       p.Width,
		Height:      p.Height,
		Depth:       p.Depth,
		Materials:   append([]string{}, p.Materials...),
		Colors:      append([]string{}, p.Colors...),
		Photos:      append([]string{}, p.Photos...),
		Seems:       append([]string{}, p.Seems...),
		Price:       p.Price,
		Description: p.Description,
		IsFeatured:  p.IsFeatured,
		Pricesale:   p.PriceSale,
	}
}

func ProductFilterFromRpc(f *productsRPC.ProductFilter) *ProductFilter {
	if f == nil {
		return nil
	}
	return &ProductFilter{
		Brand:     append([]string{}, f.GetBrand()...),
		Country:   append([]string{}, f.GetCountry()...),
		Category:  append([]string{}, f.GetCategory()...),
		MinWidth:  f.GetMinWidth(),
		MaxWidth:  f.GetMaxWidth(),
		MinHeight: f.GetMinHeight(),
		MaxHeight: f.GetMaxHeight(),
		MinDepth:  f.GetMinDepth(),
		MaxDepth:  f.GetMaxDepth(),
		Materials: append([]string{}, f.GetMaterials()...),
		Colors:    append([]string{}, f.GetColors()...),
		MinPrice:  f.GetMinPrice(),
		MaxPrice:  f.GetMaxPrice(),
		SortBy:    f.GetSortBy(),
		SortOrder: f.GetSortOrder(),
		Offset:    f.GetOffset(),
		Limit:     f.GetLimit(),
		Title:     f.GetTitle(),
	}
}

func ProductFilterToRpc(f *ProductFilter) *productsRPC.ProductFilter {
	if f == nil {
		return nil
	}
	return &productsRPC.ProductFilter{
		Brand:     append([]string{}, f.Brand...),
		Country:   append([]string{}, f.Country...),
		Category:  append([]string{}, f.Category...),
		MinWidth:  f.MinWidth,
		MaxWidth:  f.MaxWidth,
		MinHeight: f.MinHeight,
		MaxHeight: f.MaxHeight,
		MinDepth:  f.MinDepth,
		MaxDepth:  f.MaxDepth,
		Materials: append([]string{}, f.Materials...),
		Colors:    append([]string{}, f.Colors...),
		MinPrice:  f.MinPrice,
		MaxPrice:  f.MaxPrice,
		SortBy:    f.SortBy,
		SortOrder: f.SortOrder,
		Offset:    f.Offset,
		Limit:     f.Limit,
		Title:     f.Title,
	}
}
