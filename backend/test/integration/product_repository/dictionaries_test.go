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

func TestGetDictionaries(t *testing.T) {
	db, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)

	ctx := context.Background()

	tx, err := db.Driver.Begin()
	assert.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	defer db.Disconnect()
	t.Cleanup(func() {
		_ = tx.Rollback()
	})

	brand := &domain.Brand{Id: uniqueID(t, "brand_dic"), Title: uniqueTitle(t, "Test Brand")}
	cat := &domain.Category{Id: uniqueID(t, "cat_dic"), Title: uniqueTitle(t, "Test Cat"), Uri: uniqueTitle(t, "test-uri")}
	cat2 := &domain.Category{Id: uniqueID(t, "cat_dic"), Title: uniqueTitle(t, "Test Cat"), Uri: uniqueTitle(t, "test-uri-2")}
	cou := &domain.Country{Id: uniqueID(t, "country_dic"), Title: uniqueTitle(t, "Test Country"), Friendly: "FriendlyName"}
	cou2 := &domain.Country{Id: uniqueID(t, "country_dic"), Title: uniqueTitle(t, "Test Country2"), Friendly: "FriendlyName"}
	m := &domain.Material{Id: uniqueID(t, "mat_dic"), Title: uniqueTitle(t, "Test Material")}
	col := &domain.Color{Id: uniqueID(t, "color_dic"), Title: uniqueTitle(t, "TestColor"), Hex: "#123456"}
	pmax := &domain.ProductId{
		Id:          uniqueID(t, "pmax_dic"),
		Title:       uniqueTitle(t, "test"),
		Article:     uniqueTitle(t, "article_max"),
		Brand:       brand.Id,
		Category:    cat.Id,
		Country:     cou.Id,
		Width:       12,
		Height:      12,
		Depth:       12,
		Materials:   []string{m.Id},
		Colors:      []string{col.Id},
		Photos:      []string{"test"},
		Seems:       nil,
		Price:       12,
		Description: "test",
	}
	pmin := &domain.ProductId{
		Id:          uniqueID(t, "pmin_dic"),
		Title:       uniqueTitle(t, "test2"),
		Article:     uniqueTitle(t, "article_min"),
		Brand:       brand.Id,
		Category:    cat2.Id,
		Country:     cou2.Id,
		Width:       5,
		Height:      5,
		Depth:       5,
		Materials:   []string{m.Id},
		Colors:      []string{col.Id},
		Photos:      []string{"test"},
		Seems:       []string{pmax.Id},
		Price:       5,
		Description: "test",
	}
	err = driver.CreateBrand(ctx, brand)
	assert.NoError(t, err)
	err = driver.CreateCategory(ctx, cat)
	assert.NoError(t, err)
	err = driver.CreateCategory(ctx, cat2)
	assert.NoError(t, err)
	err = driver.CreateCountry(ctx, cou)
	assert.NoError(t, err)
	err = driver.CreateCountry(ctx, cou2)
	assert.NoError(t, err)
	err = driver.CreateMaterial(ctx, m)
	assert.NoError(t, err)
	err = driver.CreateColor(ctx, col)
	assert.NoError(t, err)
	assert.NoError(t, driver.CreateProduct(ctx, pmax))
	assert.NoError(t, driver.CreateProduct(ctx, pmin))

	d, err := driver.GetDictionaries(ctx)

	assert.NoError(t, err)
	t.Logf("get all dictionaries: brands=%d categories=%d", len(d.Brands), len(d.Categories))
	assert.Equal(t, 5, int(d.MinHeight))
	assert.Contains(t, d.Countries, cou)
	assert.Contains(t, d.Countries, cou2)

	d, err = driver.GetDictionariesByCategoryID(ctx, cat.Id)
	assert.NoError(t, err)
	t.Logf("get dictionaries by category: brands=%d categories=%d", len(d.Brands), len(d.Categories))
	assert.Equal(t, 12, int(d.MinHeight))
	assert.Contains(t, d.Countries, cou)
	assert.NotContains(t, d.Countries, cou2)
}
