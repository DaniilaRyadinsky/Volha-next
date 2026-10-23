package domain

import (
	"encoding/json"
	"reflect"
	"testing"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/stretchr/testify/require"
)

func jsonKeys(t *testing.T, v any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(v)
	require.NoError(t, err)
	var m map[string]any
	require.NoError(t, json.Unmarshal(raw, &m))
	return m
}
func TestProductJSONContract(t *testing.T) {
	p := ProductFromRpc(&productsRPC.Product{Id: "p", Article: "12345678", Seems: []*productsRPC.Product{{Id: "similar"}}, Photos: []string{"hidden.jpg"}})
	m := jsonKeys(t, p)
	require.Equal(t, "12345678", m["sku"])
	require.Nil(t, m["brand"])
	require.Nil(t, m["category"])
	require.Nil(t, m["country"])
	require.Equal(t, []any{}, m["materials"])
	require.Equal(t, []any{}, m["colors"])
	require.Len(t, m["similarProducts"], 1)
	for _, key := range []string{"article", "seems", "photos", "pricesale", "isFeatured"} {
		require.NotContains(t, m, key)
	}
	require.Len(t, m, 15)
	d := jsonKeys(t, DictionariesFromRpc(&productsRPC.Dictionaries{}))
	require.Contains(t, d, "minPrice")
	require.NotContains(t, d, "min_price")
	require.Equal(t, []any{}, d["brands"])
	require.NotContains(t, jsonKeys(t, &Country{Id: "c", Friendly: "hidden"}), "friendly")
	require.Equal(t, "Subtitle", jsonKeys(t, &Slide{Badge: "Subtitle"})["subtitle"])
	require.Equal(t, float64(123), jsonKeys(t, &Article{CreationTime: 123})["creationTime"])
}
func TestModelRoundTrips(t *testing.T) {
	p := &Product{Id: "p", Title: "Desk", Article: "12345678", Slug: "desk", Brand: &Brand{Id: "b", Title: "Brand"}, Category: &Category{Id: "c", Slug: "cat", Count: 4}, Country: &Country{Id: "co", Friendly: "Country"}, Width: 10, Height: 20, Depth: 30, Materials: []*Material{{Id: "m"}}, Colors: []*Color{{Id: "clr", Hex: "#000000"}}, Photos: []string{"a.jpg"}, Seems: []*Product{}, Price: 100, PriceSale: 90, IsFeatured: true}
	require.Equal(t, p, ProductFromRpc(ProductToRpc(p)))
	id := &ProductId{Id: "p", Title: "Desk", Article: "12345678", Materials: []string{}, Colors: []string{}, Photos: []string{}, Seems: []string{}}
	require.Equal(t, id, ProductIdFromRpc(ProductIdToRpc(id)))
	filter := &ProductFilter{Brand: []string{}, Country: []string{}, Category: []string{}, Materials: []string{}, Colors: []string{}, Limit: 20, Title: "Desk"}
	require.Equal(t, filter, ProductFilterFromRpc(ProductFilterToRpc(filter)))
	require.Nil(t, ProductFromRpc(nil))
	require.Nil(t, ProductToRpc(nil))
	require.Nil(t, ProductIdFromRpc(nil))
	require.Nil(t, ProductIdToRpc(nil))
	require.Nil(t, ProductFilterFromRpc(nil))
	require.Nil(t, ProductFilterToRpc(nil))
	require.NotNil(t, NewEmptyProduct().Photos)
	require.Empty(t, ProductsFromRpc([]*productsRPC.Product{nil}))
	require.Empty(t, ProductsToRpc([]*Product{nil}))
	require.Empty(t, ProductCatalogItemsToRpc([]*ProductCatalogItem{nil}))
	for _, fn := range []any{BrandFromRpc, CategoryFromRpc, CountryFromRpc, MaterialFromRpc, ColorFromRpc, SlideFromRpc, ProductColorPhotosFromRpc, DictionariesFromRpc, ArticleFromRpc} {
		require.True(t, reflect.ValueOf(fn).Call([]reflect.Value{reflect.Zero(reflect.TypeOf(fn).In(0))})[0].IsNil())
	}
}

func TestBrandConversion(t *testing.T) {
	v := &Brand{Id: "id"}
	require.Equal(t, v, BrandFromRpc(BrandToRpc(v)))
	require.Empty(t, BrandsFromRpc([]*productsRPC.Brand{nil}))
	require.Empty(t, BrandsToRpc([]*Brand{nil}))
	require.Len(t, BrandsFromRpc(BrandsToRpc([]*Brand{v})), 1)
	require.Nil(t, BrandToRpc(nil))
}

func TestCategoryConversion(t *testing.T) {
	v := &Category{Id: "id"}
	require.Equal(t, v, CategoryFromRpc(CategoryToRpc(v)))
	require.Empty(t, CategoriesFromRpc([]*productsRPC.Category{nil}))
	require.Empty(t, CategoriesToRpc([]*Category{nil}))
	require.Len(t, CategoriesFromRpc(CategoriesToRpc([]*Category{v})), 1)
	require.Nil(t, CategoryToRpc(nil))
}

func TestCountryConversion(t *testing.T) {
	v := &Country{Id: "id"}
	require.Equal(t, v, CountryFromRpc(CountryToRpc(v)))
	require.Empty(t, CountriesFromRpc([]*productsRPC.Country{nil}))
	require.Empty(t, CountriesToRpc([]*Country{nil}))
	require.Len(t, CountriesFromRpc(CountriesToRpc([]*Country{v})), 1)
	require.Nil(t, CountryToRpc(nil))
}

func TestMaterialConversion(t *testing.T) {
	v := &Material{Id: "id"}
	require.Equal(t, v, MaterialFromRpc(MaterialToRpc(v)))
	require.Empty(t, MaterialsFromRpc([]*productsRPC.Material{nil}))
	require.Empty(t, MaterialsToRpc([]*Material{nil}))
	require.Len(t, MaterialsFromRpc(MaterialsToRpc([]*Material{v})), 1)
	require.Nil(t, MaterialToRpc(nil))
}

func TestColorConversion(t *testing.T) {
	v := &Color{Id: "id"}
	require.Equal(t, v, ColorFromRpc(ColorToRpc(v)))
	require.Empty(t, ColorsFromRpc([]*productsRPC.Color{nil}))
	require.Empty(t, ColorsToRpc([]*Color{nil}))
	require.Len(t, ColorsFromRpc(ColorsToRpc([]*Color{v})), 1)
	require.Nil(t, ColorToRpc(nil))
}

func TestSlideConversion(t *testing.T) {
	v := &Slide{Id: "id"}
	require.Equal(t, v, SlideFromRpc(SlideToRpc(v)))
	require.Empty(t, SlidesFromRpc([]*productsRPC.Slide{nil}))
	require.Empty(t, SlidesToRpc([]*Slide{nil}))
	require.Len(t, SlidesFromRpc(SlidesToRpc([]*Slide{v})), 1)
	require.Nil(t, SlideToRpc(nil))
}

func TestArticleConversion(t *testing.T) {
	v := &Article{Id: "id", Img: []string{}}
	require.Equal(t, v, ArticleFromRpc(ArticleToRpc(v)))
	require.Empty(t, ArticlesFromRpc([]*productsRPC.Article{nil}))
	require.Empty(t, ArticlesToRpc([]*Article{nil}))
	require.Len(t, ArticlesFromRpc(ArticlesToRpc([]*Article{v})), 1)
	require.Nil(t, ArticleToRpc(nil))
}

func TestCompositeConversions(t *testing.T) {
	p := &ProductColorPhotos{ProductId: "p", ColorId: "c", Photos: []string{"a.jpg"}}
	require.Equal(t, p, ProductColorPhotosFromRpc(ProductColorPhotosToRpc(p)))
	require.Len(t, ProductColorPhotosListToRpc([]ProductColorPhotos{*p}), 1)
	require.Nil(t, ProductColorPhotosToRpc(nil))
	d := DictionariesFromRpc(&productsRPC.Dictionaries{})
	require.Equal(t, d, DictionariesFromRpc(DictionariesToRpc(d)))
	require.Nil(t, DictionariesToRpc(nil))
}
