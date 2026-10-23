//go:build integration

package e2e

import (
	"net/http"
	"testing"

	"github.com/autumnterror/volha-backend/test/integration/testkit"
	"github.com/stretchr/testify/require"
)

// The public/admin HTTP contracts are exercised against both real gRPC services
// and PostgreSQL, including every content update endpoint.
func TestIntegrationGatewayContentMutations(t *testing.T) {
	s := testkit.New(t)
	unique := testkit.ID()
	photo := upload(t, s, "content.png", []byte("\x89PNG\r\n\x1a\nfixture"))
	pdf := upload(t, s, "content.pdf", []byte("%PDF-1.4\nfixture"))
	brand := createID(t, s, "/api/brand", map[string]any{"title": "Brand " + unique}, "brands")
	category := createID(t, s, "/api/category", map[string]any{"title": "Category " + unique, "slug": "category-" + unique, "img": photo}, "categories")
	country := createID(t, s, "/api/country", map[string]any{"title": "Country " + unique}, "countries")
	material := createID(t, s, "/api/material", map[string]any{"title": "Material " + unique}, "materials")
	color := createID(t, s, "/api/color", map[string]any{"title": "Color " + unique, "hex": "#123456"}, "colors")

	for _, item := range []struct {
		path, id string
		body     map[string]any
	}{
		{"brand", brand, map[string]any{"title": "New Brand " + unique}},
		{"category", category, map[string]any{"title": "New Category " + unique, "slug": "new-category-" + unique, "img": photo}},
		{"country", country, map[string]any{"title": "New Country " + unique, "friendly": "Friendly"}},
		{"material", material, map[string]any{"title": "New Material " + unique}},
		{"color", color, map[string]any{"title": "New Color " + unique, "hex": "#abcdef"}},
	} {
		item.body["id"] = item.id
		status, body := s.JSON(http.MethodPut, "/api/"+item.path, item.body, true)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
		status, body = s.JSON(http.MethodGet, "/api/"+item.path+"?id="+item.id, nil, false)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
		require.Equal(t, item.body["title"], decode[map[string]any](t, body)["title"])
	}

	product := map[string]any{
		"title": "Product " + unique, "slug": "product-" + unique, "sku": testkit.SKU(t),
		"brand": brand, "category": category, "country": country,
		"width": 10, "height": 20, "depth": 30, "price": 100, "pricesale": 80,
		"materials": []string{material}, "colors": []string{color}, "photos": []string{photo},
	}
	productID := createID(t, s, "/api/product", product, "products")
	product["id"] = productID
	product["price"] = 130
	status, body := s.JSON(http.MethodPut, "/api/product", product, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/product?id="+productID, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(130), decode[map[string]any](t, body)["price"])

	slide := createID(t, s, "/api/slide", map[string]any{"link": "/", "img": photo, "img762": photo}, "slides")
	status, body = s.JSON(http.MethodPut, "/api/slide", map[string]any{
		"id": slide, "link": "/api/product?id=" + productID, "img": photo, "img762": photo, "title": "Updated slide",
	}, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/slide?id="+slide, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, "Updated slide", decode[map[string]any](t, body)["title"])

	article := createID(t, s, "/api/article", map[string]any{
		"title": "Guide", "slug": "guide-" + unique, "img": []string{photo}, "text": "/pdfs/" + pdf,
	}, "articles")
	status, body = s.JSON(http.MethodPut, "/api/article", map[string]any{
		"id": article, "title": "Updated guide", "slug": "guide-" + unique, "img": []string{photo}, "text": "/pdfs/" + pdf,
	}, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/article?id="+article, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, "Updated guide", decode[map[string]any](t, body)["title"])

	for _, item := range []struct{ path, id string }{
		{"article", article}, {"slide", slide}, {"product", productID},
		{"category", category}, {"brand", brand}, {"country", country}, {"material", material}, {"color", color},
	} {
		status, body = s.JSON(http.MethodDelete, "/api/"+item.path+"?id="+item.id, nil, true)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
	}
	status, body = s.JSON(http.MethodGet, "/images/"+photo, nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
	status, body = s.JSON(http.MethodGet, "/pdfs/"+pdf, nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
	orphan := upload(t, s, "orphan.png", []byte("\x89PNG\r\n\x1a\nfixture"))
	status, body = s.JSON(http.MethodDelete, "/api/files?title="+orphan, nil, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/images/"+orphan, nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
}
