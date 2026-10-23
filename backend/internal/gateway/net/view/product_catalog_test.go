package view

import (
	"encoding/json"
	"testing"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/stretchr/testify/require"
)

func TestCatalogViewsPreserveNullableCategoryAndEmptyPhotos(t *testing.T) {
	for _, item := range []ProductCatalogItem{CatalogItemFromRPC(&productsRPC.ProductCatalogItem{Id: "p"}), CatalogItemFromProduct(&productsRPC.Product{Id: "p"})} {
		raw, err := json.Marshal(item)
		require.NoError(t, err)
		require.JSONEq(t, `{"id":"p","slug":"","title":"","category":null,"photos":[],"price":0}`, string(raw))
	}
	got := CatalogItemFromRPC(&productsRPC.ProductCatalogItem{CategoryDetails: &productsRPC.Category{Id: "c"}, Photos: []string{"a.jpg", "b.jpg"}})
	require.Equal(t, "c", got.Category.Id)
	require.Len(t, got.Photos, 2)
}
