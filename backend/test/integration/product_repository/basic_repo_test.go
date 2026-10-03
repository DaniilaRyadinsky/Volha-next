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
	"github.com/stretchr/testify/assert"
)

func TestCatalogRepoGood(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)

	tx, err := db.Driver.Begin()
	assert.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	t.Cleanup(func() {
		_ = tx.Rollback()
		_ = db.Disconnect()
	})

	brand := &domain.Brand{Id: uniqueID(t, "brand"), Title: uniqueTitle(t, "Test Brand")}
	assert.NoError(t, driver.CreateBrand(ctx, brand))
	brands, err := driver.GetAllBrands(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, brands)
	assert.NoError(t, driver.UpdateBrand(ctx, &domain.Brand{Id: brand.Id, Title: uniqueTitle(t, "Updated Brand")}))
	_, err = driver.GetBrand(ctx, brand.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteBrand(ctx, brand.Id))

	category := &domain.Category{Id: uniqueID(t, "cat"), Title: uniqueTitle(t, "Test Cat"), Uri: uniqueTitle(t, "test-uri")}
	assert.NoError(t, driver.CreateCategory(ctx, category))
	_, err = driver.GetAllCategories(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateCategory(ctx, &domain.Category{Id: category.Id, Title: uniqueTitle(t, "Updated"), Uri: uniqueTitle(t, "upd-uri")}))
	_, err = driver.GetCategory(ctx, category.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteCategory(ctx, category.Id))

	country := &domain.Country{Id: uniqueID(t, "country"), Title: uniqueTitle(t, "Test Country"), Friendly: "FriendlyName"}
	assert.NoError(t, driver.CreateCountry(ctx, country))
	_, err = driver.GetAllCountries(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateCountry(ctx, &domain.Country{Id: country.Id, Title: uniqueTitle(t, "Updated"), Friendly: "UpdatedFriendly"}))
	_, err = driver.GetCountry(ctx, country.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteCountry(ctx, country.Id))

	material := &domain.Material{Id: uniqueID(t, "mat"), Title: uniqueTitle(t, "Test Material")}
	assert.NoError(t, driver.CreateMaterial(ctx, material))
	_, err = driver.GetAllMaterials(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateMaterial(ctx, &domain.Material{Id: material.Id, Title: uniqueTitle(t, "Updated Material")}))
	_, err = driver.GetMaterial(ctx, material.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteMaterial(ctx, material.Id))

	color := &domain.Color{Id: uniqueID(t, "color"), Title: uniqueTitle(t, "TestColor"), Hex: "#123456"}
	assert.NoError(t, driver.CreateColor(ctx, color))
	_, err = driver.GetAllColors(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateColor(ctx, &domain.Color{Id: color.Id, Title: uniqueTitle(t, "UpdatedColor"), Hex: "#654321"}))
	_, err = driver.GetColor(ctx, color.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteColor(ctx, color.Id))

	slide := &domain.Slide{Id: uniqueID(t, "slide"), Link: "http://example.com", Img: "img1.jpg", Img762: "img762.jpg"}
	assert.NoError(t, driver.CreateSlide(ctx, slide))
	_, err = driver.GetAllSlides(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateSlide(ctx, &domain.Slide{Id: slide.Id, Link: "http://new.example.com", Img: "newimg1.jpg", Img762: "newimg762.jpg"}))
	_, err = driver.GetSlide(ctx, slide.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteSlide(ctx, slide.Id))

	article := &domain.Article{Id: uniqueID(t, "article"), Title: uniqueTitle(t, "Title"), Img: "img.jpg", Text: "# heading", CreationTime: time.Now().Unix()}
	assert.NoError(t, driver.CreateArticle(ctx, article))
	_, err = driver.GetAllArticles(ctx)
	assert.NoError(t, err)
	assert.NoError(t, driver.UpdateArticle(ctx, &domain.Article{Id: article.Id, Title: uniqueTitle(t, "newTitle"), Img: "newimg.jpg", Text: "# newheading", CreationTime: time.Now().Unix()}))
	_, err = driver.GetArticle(ctx, article.Id)
	assert.NoError(t, err)
	assert.NoError(t, driver.DeleteArticle(ctx, article.Id))
}

func TestCatalogRepoNotFound(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	db, err := psql.NewConnect(config.Test())
	assert.NoError(t, err)

	tx, err := db.Driver.Begin()
	assert.NoError(t, err)

	driver := repository.Driver{Driver: tx}
	t.Cleanup(func() {
		_ = tx.Rollback()
		_ = db.Disconnect()
	})

	assert.ErrorIs(t, driver.UpdateBrand(ctx, &domain.Brand{Id: "not_exist", Title: "Test Brand"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteBrand(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateCategory(ctx, &domain.Category{Id: "not_exist", Title: "Test Cat", Uri: "test-uri"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteCategory(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateCountry(ctx, &domain.Country{Id: "not_exist", Title: "Test Country", Friendly: "FriendlyName"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteCountry(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateMaterial(ctx, &domain.Material{Id: "not_exist", Title: "Test Material"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteMaterial(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateColor(ctx, &domain.Color{Id: "not_exist", Title: "TestColor", Hex: "#123456"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteColor(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateSlide(ctx, &domain.Slide{Id: "not_exist"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteSlide(ctx, "not_exist"), domain.ErrNotFound)
	assert.ErrorIs(t, driver.UpdateArticle(ctx, &domain.Article{Id: "not_exist"}), domain.ErrNotFound)
	assert.ErrorIs(t, driver.DeleteArticle(ctx, "not_exist"), domain.ErrNotFound)
}
