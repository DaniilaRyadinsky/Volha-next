//go:build integration

package integration

import (
	"context"
	"testing"

	"github.com/autumnterror/volha-backend/internal/product-service/infra/psql"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"

	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/stretchr/testify/assert"
)

func TestGoodProductColorPhotos(t *testing.T) {
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
	category := &domain.Category{Id: uniqueID(t, "cat"), Title: uniqueTitle(t, "TestCategory"), Slug: uniqueTitle(t, "test-cat")}
	country := &domain.Country{Id: uniqueID(t, "country"), Title: uniqueTitle(t, "TestCountry"), Friendly: "Friendly"}
	material := &domain.Material{Id: uniqueID(t, "mat"), Title: uniqueTitle(t, "Leather")}
	color := &domain.Color{Id: uniqueID(t, "color"), Title: uniqueTitle(t, "Red"), Hex: "#FF0000"}
	assert.NoError(t, driver.CreateBrand(ctx, brand))
	assert.NoError(t, driver.CreateCategory(ctx, category))
	assert.NoError(t, driver.CreateCountry(ctx, country))
	assert.NoError(t, driver.CreateMaterial(ctx, material))
	assert.NoError(t, driver.CreateColor(ctx, color))

	productID := uniqueID(t, "product")
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

	c := &domain.ProductColorPhotos{ProductId: productID, ColorId: color.Id, Photos: []string{"photo1", "photo2"}}
	updated := &domain.ProductColorPhotos{ProductId: productID, ColorId: color.Id, Photos: []string{"new photo1", "new photo2"}}

	assert.NoError(t, driver.CreateProductColorPhotos(ctx, c))

	all, err := driver.GetAllProductColorPhotos(ctx)
	assert.NoError(t, err)
	t.Logf("product color photos after create: %d", len(all))

	assert.NoError(t, driver.UpdateProductColorPhotos(ctx, updated))

	ph, err := driver.GetPhotosByProductColor(ctx, c.ProductId, c.ColorId)
	assert.NoError(t, err)
	t.Logf("get photos after update: %d", len(ph))

	assert.NoError(t, driver.DeleteProductColorPhotos(ctx, c.ProductId, c.ColorId))
}

func TestBadProductColorPhotos(t *testing.T) {
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

	pcp := &domain.ProductColorPhotos{ProductId: "product_test_bad", ColorId: "color_test_bad", Photos: []string{"photo1", "photo2"}}
	tx1, err := db.Driver.Begin()
	assert.NoError(t, err)
	{
		driver := repository.Driver{Driver: tx1}
		assert.ErrorIs(t, driver.CreateProductColorPhotos(ctx, pcp), domain.ErrForeignKey)
	}
	_ = tx1.Rollback()

	assert.ErrorIs(
		t,
		driver.UpdateProductColorPhotos(
			ctx,
			&domain.ProductColorPhotos{ProductId: "not exist", ColorId: "not exist", Photos: []string{"photo1", "photo2"}},
		),
		domain.ErrNotFound,
	)

	assert.ErrorIs(t, driver.DeleteProductColorPhotos(ctx, pcp.ProductId, pcp.ColorId), domain.ErrNotFound)

}
