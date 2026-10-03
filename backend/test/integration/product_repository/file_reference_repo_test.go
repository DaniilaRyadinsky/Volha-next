//go:build integration

package integration

import (
	"context"
	"testing"
	"time"

	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/infra/psql"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
	"github.com/stretchr/testify/require"
)

func TestFileReferenceAndProductColorPhotosByProduct(t *testing.T) {
	ctx := context.Background()

	db, err := psql.NewConnect(config.Test())
	require.NoError(t, err)

	tx, err := db.Driver.Begin()
	require.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	t.Cleanup(func() {
		_ = tx.Rollback()
		_ = db.Disconnect()
	})

	brand := &domain.Brand{Id: "brand_ref", Title: "Reference Brand"}
	category := &domain.Category{Id: "cat_ref", Title: "Reference Cat", Uri: "reference-cat", Img: "category-ref.jpg"}
	country := &domain.Country{Id: "country_ref", Title: "Reference Country", Friendly: "Reference Friendly"}
	material := &domain.Material{Id: "mat_ref", Title: "Reference Material"}
	color := &domain.Color{Id: "color_ref", Title: "Reference Color", Hex: "#123456"}

	require.NoError(t, driver.CreateBrand(ctx, brand))
	require.NoError(t, driver.CreateCategory(ctx, category))
	require.NoError(t, driver.CreateCountry(ctx, country))
	require.NoError(t, driver.CreateMaterial(ctx, material))
	require.NoError(t, driver.CreateColor(ctx, color))
	require.NoError(t, driver.CreateProduct(ctx, &domain.ProductId{
		Id:          "product_ref",
		Title:       "Reference Product",
		Article:     "reference-product",
		Brand:       brand.Id,
		Category:    category.Id,
		Country:     country.Id,
		Width:       10,
		Height:      20,
		Depth:       30,
		Materials:   []string{material.Id},
		Colors:      []string{color.Id},
		Photos:      []string{"product-ref.jpg"},
		Seems:       []string{},
		Price:       100,
		Description: "reference product",
	}))
	require.NoError(t, driver.CreateProductColorPhotos(ctx, &domain.ProductColorPhotos{
		ProductId: "product_ref",
		ColorId:   color.Id,
		Photos:    []string{"pcp-ref.jpg"},
	}))
	require.NoError(t, driver.CreateSlide(ctx, &domain.Slide{
		Id:     "slide_ref",
		Link:   "https://example.com",
		Img:    "slide-ref.jpg",
		Img762: "slide-762-ref.jpg",
	}))
	require.NoError(t, driver.CreateArticle(ctx, &domain.Article{
		Id:           "article_ref",
		Img:          "article-ref.jpg",
		Title:        "Reference Article",
		Text:         "reference article",
		CreationTime: time.Now().UTC().Unix(),
	}))

	for _, filename := range []string{
		"category-ref.jpg",
		"product-ref.jpg",
		"pcp-ref.jpg",
		"slide-ref.jpg",
		"slide-762-ref.jpg",
		"article-ref.jpg",
	} {
		referenced, err := driver.IsFileReferenced(ctx, filename)
		require.NoError(t, err)
		require.True(t, referenced, filename)
	}

	referenced, err := driver.IsFileReferenced(ctx, "unused-ref.jpg")
	require.NoError(t, err)
	require.False(t, referenced)

	pcps, err := driver.GetProductColorPhotosByProduct(ctx, "product_ref")
	require.NoError(t, err)
	require.Len(t, pcps, 1)
	require.Equal(t, []string{"pcp-ref.jpg"}, pcps[0].Photos)
}
