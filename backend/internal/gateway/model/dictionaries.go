package model

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type DictionariesByCategory struct {
	Brands     []*productsRPC.Brand    `json:"brands"`
	Categories []*productsRPC.Category `json:"categories"`
	Countries  []*productsRPC.Country  `json:"countries"`
	Materials  []*productsRPC.Material `json:"materials"`
	Colors     []*productsRPC.Color    `json:"colors"`

	MinPrice  int32 `json:"minPrice"`
	MaxPrice  int32 `json:"maxPrice"`
	MinWidth  int32 `json:"minWidth"`
	MaxWidth  int32 `json:"maxWidth"`
	MinHeight int32 `json:"minHeight"`
	MaxHeight int32 `json:"maxHeight"`
	MinDepth  int32 `json:"minDepth"`
	MaxDepth  int32 `json:"maxDepth"`

	FilteredCategory string `json:"filtered_category"`
}
