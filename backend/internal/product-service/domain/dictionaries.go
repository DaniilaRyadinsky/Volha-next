package domain

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type Dictionaries struct {
	Brands     []*Brand    `json:"brands"`
	Categories []*Category `json:"categories"`
	Countries  []*Country  `json:"countries"`
	Materials  []*Material `json:"materials"`
	Colors     []*Color    `json:"colors"`

	MinPrice  int32 `json:"minPrice"`
	MaxPrice  int32 `json:"maxPrice"`
	MinWidth  int32 `json:"minWidth"`
	MaxWidth  int32 `json:"maxWidth"`
	MinHeight int32 `json:"minHeight"`
	MaxHeight int32 `json:"maxHeight"`
	MinDepth  int32 `json:"minDepth"`
	MaxDepth  int32 `json:"maxDepth"`
}

func DictionariesFromRpc(d *productsRPC.Dictionaries) *Dictionaries {
	if d == nil {
		return nil
	}
	return &Dictionaries{
		Brands:     BrandsFromRpc(d.GetBrands()),
		Categories: CategoriesFromRpc(d.GetCategories()),
		Countries:  CountriesFromRpc(d.GetCountries()),
		Materials:  MaterialsFromRpc(d.GetMaterials()),
		Colors:     ColorsFromRpc(d.GetColors()),

		MinPrice:  d.GetMinPrice(),
		MaxPrice:  d.GetMaxPrice(),
		MinWidth:  d.GetMinWidth(),
		MaxWidth:  d.GetMaxWidth(),
		MinHeight: d.GetMinHeight(),
		MaxHeight: d.GetMaxHeight(),
		MinDepth:  d.GetMinDepth(),
		MaxDepth:  d.GetMaxDepth(),
	}
}

func DictionariesToRpc(d *Dictionaries) *productsRPC.Dictionaries {
	if d == nil {
		return nil
	}
	return &productsRPC.Dictionaries{
		Brands:     BrandsToRpc(d.Brands),
		Categories: CategoriesToRpc(d.Categories),
		Countries:  CountriesToRpc(d.Countries),
		Materials:  MaterialsToRpc(d.Materials),
		Colors:     ColorsToRpc(d.Colors),

		MinPrice:  d.MinPrice,
		MaxPrice:  d.MaxPrice,
		MinWidth:  d.MinWidth,
		MaxWidth:  d.MaxWidth,
		MinHeight: d.MinHeight,
		MaxHeight: d.MaxHeight,
		MinDepth:  d.MinDepth,
		MaxDepth:  d.MaxDepth,
	}
}
