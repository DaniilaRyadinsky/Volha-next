package model

import productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

type DictionariesByCategory struct {
	Brands     []*productsRPC.Brand    `json:"brands"`
	Categories []*productsRPC.Category `json:"categories"`
	Countries  []*productsRPC.Country  `json:"countries"`
	Materials  []*productsRPC.Material `json:"materials"`
	Colors     []*productsRPC.Color    `json:"colors"`

	MinPrice  int32 `json:"min_price"`
	MaxPrice  int32 `json:"max_price"`
	MinWidth  int32 `json:"min_width"`
	MaxWidth  int32 `json:"max_width"`
	MinHeight int32 `json:"min_height"`
	MaxHeight int32 `json:"max_height"`
	MinDepth  int32 `json:"min_depth"`
	MaxDepth  int32 `json:"max_depth"`

	FilteredCategory string `json:"filtered_category"`
}
