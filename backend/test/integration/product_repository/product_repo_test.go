//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/infra/psql"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
	"github.com/stretchr/testify/assert"
)

func TestProductGood(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)

	tx, err := db.Driver.Begin()
	assert.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	defer t.Cleanup(func() {
		_ = tx.Rollback()
		db.Disconnect()
	})

	brand := &domain.Brand{Id: uniqueID(t, "brand"), Title: uniqueTitle(t, "TestBrand")}
	category := &domain.Category{Id: uniqueID(t, "cat"), Title: uniqueTitle(t, "TestCategory"), Slug: uniqueTitle(t, "test-cat"), Img: "cat.jpg"}
	country := &domain.Country{Id: uniqueID(t, "country"), Title: uniqueTitle(t, "TestCountry"), Friendly: "Friendly"}
	material := &domain.Material{Id: uniqueID(t, "mat"), Title: uniqueTitle(t, "Leather")}
	color := &domain.Color{Id: uniqueID(t, "color"), Title: uniqueTitle(t, "Red"), Hex: "#FF0000"}
	assert.NoError(t, driver.CreateBrand(ctx, brand))
	assert.NoError(t, driver.CreateCategory(ctx, category))
	assert.NoError(t, driver.CreateCountry(ctx, country))
	assert.NoError(t, driver.CreateMaterial(ctx, material))
	assert.NoError(t, driver.CreateColor(ctx, color))

	productID := uniqueID(t, "product")
	updatedTitle := uniqueTitle(t, "new_test")
	assert.NoError(t, driver.CreateProduct(ctx, &domain.ProductId{
		Id:          productID,
		Title:       uniqueTitle(t, "test"),
		Slug:        uniqueID(t, "slug"),
		Article:     uniqueTitle(t, "article"),
		Brand:       brand.Id,
		Category:    category.Id,
		Country:     country.Id,
		Width:       1,
		Height:      2,
		Depth:       3,
		Materials:   []string{material.Id},
		Colors:      []string{color.Id},
		Photos:      []string{"photo1", "photo2"},
		Seems:       []string{},
		Price:       100,
		Description: "test",
	}))

	p, total, err := driver.GetAllProducts(ctx, 0, 10)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, p)
		assert.NotZero(t, total)
		t.Logf("products after create: %d total: %d", len(p), total)
	}

	assert.NoError(t, driver.UpdateProduct(ctx, &domain.ProductId{
		Id:          productID,
		Title:       updatedTitle,
		Slug:        uniqueID(t, "slug"),
		Article:     uniqueTitle(t, "new_article"),
		Brand:       brand.Id,
		Category:    category.Id,
		Country:     country.Id,
		Width:       2,
		Height:      3,
		Depth:       4,
		Materials:   []string{material.Id},
		Colors:      []string{color.Id},
		Photos:      []string{"new photo1", "new photo2"},
		Seems:       []string{},
		Price:       200,
		Description: "updated",
	}))

	pr, err := driver.GetProduct(ctx, productID)
	if assert.NoError(t, err) {
		assert.Equal(t, updatedTitle, pr.Title)
		t.Logf("product after update: %s", pr.Id)
	}

	assert.NoError(t, driver.DeleteProduct(ctx, productID))
}

func TestProductFilters(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	db, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)

	tx, err := db.Driver.Begin()
	assert.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	defer db.Disconnect()
	t.Cleanup(func() { _ = tx.Rollback() })

	brand := &domain.Brand{Id: uniqueID(t, "brand_filter"), Title: uniqueTitle(t, "TestBrand")}
	category := &domain.Category{Id: uniqueID(t, "cat_filter"), Title: uniqueTitle(t, "TestCategory"), Slug: uniqueTitle(t, "test-cat-filter"), Img: "cat.jpg"}
	country := &domain.Country{Id: uniqueID(t, "country_filter"), Title: uniqueTitle(t, "TestCountry"), Friendly: "Friendly"}
	material := &domain.Material{Id: uniqueID(t, "mat_filter"), Title: uniqueTitle(t, "Leather")}
	color := &domain.Color{Id: uniqueID(t, "color_filter"), Title: uniqueTitle(t, "Red"), Hex: "#FF0000"}
	assert.NoError(t, driver.CreateBrand(ctx, brand))
	assert.NoError(t, driver.CreateCategory(ctx, category))
	assert.NoError(t, driver.CreateCountry(ctx, country))
	assert.NoError(t, driver.CreateMaterial(ctx, material))
	assert.NoError(t, driver.CreateColor(ctx, color))

	product := &domain.ProductId{
		Id:          uniqueID(t, "prod_filter"),
		Title:       uniqueTitle(t, "FilterProduct"),
		Slug:        uniqueID(t, "slug_filter"),
		Article:     uniqueTitle(t, "F-001"),
		Brand:       brand.Id,
		Category:    category.Id,
		Country:     country.Id,
		Width:       60,
		Height:      110,
		Depth:       40,
		Materials:   []string{material.Id},
		Colors:      []string{color.Id},
		Photos:      []string{"x.jpg"},
		Seems:       []string{},
		Price:       888,
		Description: "Filtered",
	}
	assert.NoError(t, driver.CreateProduct(ctx, product))

	tests := []struct {
		name   string
		filter domain.ProductFilter
	}{
		{"by brand", domain.ProductFilter{Brand: []string{brand.Id}, Limit: 10}},
		{"by category", domain.ProductFilter{Category: []string{category.Id}, Limit: 10}},
		{"by country", domain.ProductFilter{Country: []string{country.Id}, Limit: 10}},
		{"by material", domain.ProductFilter{Materials: []string{material.Id}, Limit: 10}},
		{"by color", domain.ProductFilter{Colors: []string{color.Id}, Limit: 10}},
		{"by size", domain.ProductFilter{MinWidth: 50, MaxWidth: 70, MinHeight: 100, MaxHeight: 120, MinDepth: 30, MaxDepth: 50, Limit: 10}},
		{"by price", domain.ProductFilter{MinPrice: 500, MaxPrice: 1000, Limit: 10}},
		{"with sort asc", domain.ProductFilter{SortBy: "price", SortOrder: "asc", Limit: 10}},
		{"with sort desc", domain.ProductFilter{SortBy: "width", SortOrder: "desc", Limit: 10}},
		{"with limit and offset", domain.ProductFilter{Limit: 1, Offset: 0}},
		{"by title", domain.ProductFilter{Title: "terPro", Limit: 10}},
		{"complex", domain.ProductFilter{
			Brand:     []string{brand.Id},
			Category:  []string{category.Id},
			Country:   []string{country.Id},
			Materials: []string{material.Id},
			Colors:    []string{color.Id},
			MinWidth:  50, MaxWidth: 70,
			MinHeight: 100, MaxHeight: 120,
			MinDepth: 30, MaxDepth: 50,
			MinPrice: 800, MaxPrice: 1000,
			SortBy: "title", SortOrder: "asc",
			Limit: 10,
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, total, err := driver.FilterProducts(ctx, &tt.filter)
			assert.NoError(t, err)
			assert.NotEmpty(t, result)
			assert.NotEqual(t, 0, total)
			t.Logf("Filter %q returned %d product(s)", tt.name, len(result))
		})
	}
}
