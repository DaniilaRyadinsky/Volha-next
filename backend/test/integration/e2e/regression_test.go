//go:build integration

package e2e

import (
	"net/http"
	"testing"
	"time"

	"github.com/autumnterror/volha-backend/test/integration/testkit"
	"github.com/stretchr/testify/require"
)

// A cached, empty category must reflect product creation, price edits and deletion.
func TestRegressionDictionaryCacheAfterProductMutations(t *testing.T) {
	s := testkit.New(t)
	c := s.SeedCatalog()
	photo := upload(t, s, "cache.png", []byte("\x89PNG\r\n\x1a\nfixture"))
	path := "/api/dictionaries/all/by-category?id=" + c.Category
	price := func() float64 {
		status, body := s.JSON(http.MethodGet, path, nil, false)
		require.Equal(t, http.StatusOK, status, string(body))
		return decode[map[string]any](t, body)["minPrice"].(float64)
	}
	require.Zero(t, price())
	product := map[string]any{
		"title": "Cache " + testkit.ID(), "slug": "cache-" + testkit.ID(), "sku": testkit.SKU(t),
		"brand": c.Brand, "category": c.Category, "country": c.Country,
		"width": 10, "height": 20, "depth": 30, "price": 100, "pricesale": 80,
		"materials": []string{c.Material}, "colors": []string{c.Color}, "photos": []string{photo},
	}
	id := createID(t, s, "/api/product", product, "products")
	require.Equal(t, float64(100), price())
	product["id"] = id
	product["price"] = 145
	status, body := s.JSON(http.MethodPut, "/api/product", product, true)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Equal(t, float64(145), price())
	status, body = s.JSON(http.MethodDelete, "/api/product?id="+id, nil, true)
	require.Equal(t, http.StatusOK, status, string(body))
	require.Zero(t, price())
}

// Removing one owner of a file must not break the other; the last removal cleans it up.
func TestRegressionSharedImageReferenceCleanup(t *testing.T) {
	s := testkit.New(t)
	c := s.SeedCatalog()
	photo := upload(t, s, "shared.png", []byte("\x89PNG\r\n\x1a\nfixture"))
	product := createID(t, s, "/api/product", map[string]any{
		"title": "Shared " + testkit.ID(), "slug": "shared-" + testkit.ID(), "sku": testkit.SKU(t),
		"brand": c.Brand, "category": c.Category, "country": c.Country,
		"width": 10, "height": 20, "depth": 30, "price": 100, "pricesale": 80,
		"materials": []string{c.Material}, "colors": []string{c.Color}, "photos": []string{photo},
	}, "products")
	slide := createID(t, s, "/api/slide", map[string]any{
		"link": "/api/product?id=" + product, "img": photo, "img762": photo,
	}, "slides")
	status, body := s.JSON(http.MethodDelete, "/api/product?id="+product, nil, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/images/"+photo, nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodDelete, "/api/slide?id="+slide, nil, true)
	require.Equal(t, http.StatusOK, status, string(body))
	status, body = s.JSON(http.MethodGet, "/images/"+photo, nil, false)
	require.Equal(t, http.StatusNotFound, status, string(body))
}

// An old refresh cookie must be rejected after a successful rotation.
func TestRegressionRefreshTokenReplay(t *testing.T) {
	s := testkit.New(t)
	unique := testkit.ID()
	status, body := s.JSON(http.MethodPost, "/api/auth/register", map[string]any{
		"email": unique + "@example.com", "login": "u" + unique, "password": "StrongPassword123",
	}, false)
	require.Equal(t, http.StatusCreated, status, string(body))
	s.Track("users", decode[map[string]any](t, body)["id"].(string))
	base, err := http.NewRequest(http.MethodPost, s.BaseURL+"/api/auth/refresh", nil)
	require.NoError(t, err)
	cookies := s.HTTP.Jar.Cookies(base.URL)
	var oldRefresh *http.Cookie
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			oldRefresh = cookie
		}
	}
	require.NotNil(t, oldRefresh)
	status, body = s.JSON(http.MethodPost, "/api/auth/refresh", nil, false)
	require.Equal(t, http.StatusOK, status, string(body))
	replay, err := http.NewRequest(http.MethodPost, s.BaseURL+"/api/auth/refresh", nil)
	require.NoError(t, err)
	replay.AddCookie(oldRefresh)
	response, err := (&http.Client{Timeout: 5 * time.Second}).Do(replay)
	require.NoError(t, err)
	defer response.Body.Close()
	require.Equal(t, http.StatusUnauthorized, response.StatusCode)
}
