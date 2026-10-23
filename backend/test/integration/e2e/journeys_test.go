//go:build integration

package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"testing"

	"github.com/autumnterror/volha-backend/test/integration/testkit"
	"github.com/stretchr/testify/require"
)

func TestE2EAdminPublishesAndVisitorBrowses(t *testing.T) {
	s := testkit.New(t)
	photo := upload(t, s, "image.png", []byte("\x89PNG\r\n\x1a\nfixture"))
	pdf := upload(t, s, "guide.pdf", []byte("%PDF-1.4\nfixture"))
	base := testkit.ID()
	brand := createID(t, s, "/api/brand", map[string]any{"title": "Brand " + base}, "brands")
	category := createID(t, s, "/api/category", map[string]any{"title": "Category " + base, "slug": "cat-" + base, "img": photo}, "categories")
	country := createID(t, s, "/api/country", map[string]any{"title": "Country " + base}, "countries")
	material := createID(t, s, "/api/material", map[string]any{"title": "Material " + base}, "materials")
	color := createID(t, s, "/api/color", map[string]any{"title": "Color " + base, "hex": "#aabbcc"}, "colors")

	status, body := s.JSON(http.MethodGet, "/api/dictionaries/all/by-category?id="+category, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	before := decode[map[string]any](t, body)
	require.Equal(t, float64(0), before["minPrice"])

	productBody := map[string]any{
		"title": "Chair " + base, "slug": "chair-" + base, "sku": testkit.SKU(t),
		"brand": brand, "category": category, "country": country,
		"width": 10, "height": 20, "depth": 30, "price": 100, "pricesale": 80,
		"materials": []string{material}, "colors": []string{color}, "photos": []string{photo}, "similarProducts": []string{},
	}
	product := createID(t, s, "/api/product", productBody, "products")
	status, body = s.JSON(http.MethodGet, "/api/dictionaries/all/by-category?id="+category, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	after := decode[map[string]any](t, body)
	require.Equal(t, float64(100), after["minPrice"], "product creation must invalidate cached price range")

	status, body = s.JSON(http.MethodGet, "/api/product?id="+product, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, product, decode[map[string]any](t, body)["id"])
	status, body = s.JSON(http.MethodGet, "/api/product/all?page=1&limit=1", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	list := decode[map[string]any](t, body)
	require.Equal(t, float64(1), list["limit"])
	require.NotEmpty(t, list["items"])
	status, body = s.JSON(http.MethodPost, "/api/product/filter", map[string]any{"category": []string{category}, "search": "Chair", "page": 1, "limit": 20}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(1), decode[map[string]any](t, body)["total"])
	status, body = s.JSON(http.MethodGet, "/api/dictionaries/all/by-category?id="+category, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(100), decode[map[string]any](t, body)["minPrice"])
	status, body = s.JSON(http.MethodGet, "/api/dictionaries/all", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))

	photoLink := map[string]any{"product_id": product, "color_id": color, "photos": []string{photo}}
	status, body = s.JSON(http.MethodPost, "/api/colorphotos", photoLink, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/colorphotos/all", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.NotEmpty(t, decode[[]any](t, body))
	status, body = s.JSON(http.MethodPost, "/api/colorphotos/photos/get", map[string]any{"product_id": product, "color_id": color}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, []string{photo}, decode[[]string](t, body))
	status, body = s.JSON(http.MethodPut, "/api/colorphotos", photoLink, true)
	require.Equal(t, http.StatusOK, status, string(body))

	slide := createID(t, s, "/api/slide", map[string]any{"link": "/api/product?id=" + product, "img": photo, "img762": photo, "title": "Featured", "subtitle": "Sale"}, "slides")
	status, body = s.JSON(http.MethodGet, "/api/slide?id="+slide, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/slide/all", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.NotEmpty(t, decode[[]any](t, body))

	article := createID(t, s, "/api/article", map[string]any{"slug": "guide-" + base, "title": "Buying guide", "img": []string{photo}, "text": "/pdfs/" + pdf}, "articles")
	status, body = s.JSON(http.MethodGet, "/api/article?id="+article, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, "/pdfs/"+pdf, decode[map[string]any](t, body)["text"])
	status, body = s.JSON(http.MethodGet, "/api/article/all?page=1&limit=20", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.NotEmpty(t, decode[map[string]any](t, body)["items"])

	// Read endpoints for each catalog family must work through the gateway.
	for _, item := range []struct{ path, id string }{
		{"brand", brand}, {"category", category}, {"country", country}, {"material", material}, {"color", color},
	} {
		status, body = s.JSON(http.MethodGet, "/api/"+item.path+"?id="+item.id, nil, false)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
		status, body = s.JSON(http.MethodGet, "/api/"+item.path+"/all", nil, false)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
	}
	status, body = s.JSON(http.MethodGet, "/images/"+photo, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/pdfs/"+pdf, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))

	// Shared files remain protected while any entity still points to them.
	status, body = s.JSON(http.MethodDelete, "/api/files?title="+photo, nil, true)
	require.Equal(t, http.StatusConflict, status, string(body))
	status, body = s.JSON(http.MethodDelete, "/api/colorphotos", map[string]any{"product_id": product, "color_id": color}, true)
	require.Equal(t, http.StatusOK, status, string(body))
	for _, item := range []struct{ path, id string }{
		{"article", article}, {"slide", slide}, {"product", product}, {"category", category},
	} {
		status, body = s.JSON(http.MethodDelete, "/api/"+item.path+"?id="+item.id, nil, true)
		require.Equal(t, http.StatusOK, status, item.path+": "+string(body))
	}
	status, body = s.JSON(http.MethodGet, "/api/product?id="+product, nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
}

func TestE2EAccountAndCartHappyPath(t *testing.T) {
	s := testkit.New(t)
	c := s.SeedCatalog()
	p := s.SeedProduct(c, "cart.jpg")
	unique := testkit.ID()
	password := "StrongPassword123"
	status, body := s.JSON(http.MethodPost, "/api/auth/register", map[string]any{"email": unique + "@example.com", "login": "u" + unique, "password": password}, false)
	require.Equal(t, http.StatusCreated, status, string(body))
	user := decode[map[string]any](t, body)
	userID := user["id"].(string)
	s.Track("users", userID)
	status, body = s.JSON(http.MethodGet, "/api/auth/me", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, userID, decode[map[string]any](t, body)["id"])

	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": p.Id, "quantity": 2}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(160), decode[map[string]any](t, body)["total_price"])
	status, body = s.JSON(http.MethodGet, "/api/cart", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(2), decode[map[string]any](t, body)["total_quantity"])
	status, body = s.JSON(http.MethodPut, "/api/cart/items/"+p.Id, map[string]any{"quantity": 3}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(3), decode[map[string]any](t, body)["total_quantity"])
	status, body = s.JSON(http.MethodDelete, "/api/cart/items/"+p.Id, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": p.Id, "quantity": 1}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodDelete, "/api/cart", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	emptyCart := decode[struct {
		TotalQuantity int32 `json:"total_quantity"`
	}](t, body)
	require.Zero(t, emptyCart.TotalQuantity)

	status, body = s.JSON(http.MethodPut, "/api/auth/change-login", map[string]any{"password": password, "new_login": "changed" + unique}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPut, "/api/auth/change-password", map[string]any{"old_password": password, "new_password": "ChangedPassword123"}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/refresh", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/logout", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/auth/me", nil, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/login", map[string]any{"identifier": "changed" + unique, "password": "ChangedPassword123"}, false)
	require.Equal(t, http.StatusOK, status, string(body))
}

func TestE2EBadPaths(t *testing.T) {
	s := testkit.New(t)
	c := s.SeedCatalog()
	p := s.SeedProduct(c, "cart.jpg")
	status, body := s.JSON(http.MethodGet, "/api/cart", nil, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/login", map[string]any{"identifier": "nobody", "password": "wrong"}, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/register", map[string]any{"email": "invalid", "login": "x", "password": "short"}, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/product?id=bad", nil, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/product?id="+testkit.ID(), nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/product/all?page=0&limit=101", nil, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/product/filter", map[string]any{"sortBy": "bogus", "limit": 20}, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodGet, "/api/dictionaries/all/by-category", nil, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/brand", map[string]any{"title": "Unauthorized"}, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/product", map[string]any{"title": "Missing dependencies"}, true)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/colorphotos", map[string]any{"product_id": p.Id, "color_id": c.Color, "photos": []string{"missing.png"}}, true)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodDelete, "/api/files?title=../secret", nil, true)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/article", map[string]any{"title": "Bad", "slug": "bad", "text": "/pdfs/missing.pdf"}, true)
	require.Equal(t, http.StatusBadRequest, status, string(body))

	unique := testkit.ID()
	userData := map[string]any{"email": unique + "@example.com", "login": "u" + unique, "password": "StrongPassword123"}
	status, body = s.JSON(http.MethodPost, "/api/auth/register", userData, false)
	require.Equal(t, http.StatusCreated, status, string(body))
	s.Track("users", decode[map[string]any](t, body)["id"].(string))
	status, body = s.JSON(http.MethodPost, "/api/auth/register", userData, false)
	require.Equal(t, http.StatusConflict, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": p.Id, "quantity": 0}, false)
	require.Equal(t, http.StatusBadRequest, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": testkit.ID(), "quantity": 1}, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": p.Id, "quantity": 99}, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/cart/items", map[string]any{"product_id": p.Id, "quantity": 1}, false)
	require.Equal(t, http.StatusTooManyRequests, status, string(body))
	status, body = s.JSON(http.MethodPut, "/api/auth/change-password", map[string]any{"old_password": "wrong", "new_password": "OtherPassword123"}, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/logout", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodPost, "/api/auth/refresh", nil, false)
	require.Equal(t, http.StatusUnauthorized, status, string(body))
}

func createID(t *testing.T, s *testkit.Stack, path string, payload any, table string) string {
	t.Helper()
	status, body := s.JSON(http.MethodPost, path, payload, true)
	require.Equal(t, http.StatusOK, status, path+": "+string(body))
	id := decode[map[string]string](t, body)["id"]
	require.NotEmpty(t, id)
	s.Track(table, id)
	return id
}

func decode[T any](t *testing.T, data []byte) T {
	t.Helper()
	var result T
	require.NoError(t, json.Unmarshal(data, &result), string(data))
	return result
}

func upload(t *testing.T, s *testkit.Stack, name string, contents []byte) string {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", name)
	require.NoError(t, err)
	_, err = part.Write(contents)
	require.NoError(t, err)
	require.NoError(t, w.Close())
	req, err := http.NewRequest(http.MethodPost, s.BaseURL+"/api/files", &body)
	require.NoError(t, err)
	req.Header.Set("Content-Type", w.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+s.AdminToken)
	response, err := s.HTTP.Do(req)
	require.NoError(t, err)
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, response.StatusCode, string(data))
	return decode[map[string]string](t, data)["name"]
}
