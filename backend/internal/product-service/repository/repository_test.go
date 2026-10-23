package repository

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/lib/pq"
	"github.com/stretchr/testify/require"
)

func mockDriver(t *testing.T) (Driver, sqlmock.Sqlmock) {
	t.Helper()
	db, m, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, m.ExpectationsWereMet()); db.Close() })
	return Driver{Driver: db}, m
}

func TestBrandRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Brand{Id: "b", Title: "Brand"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM brands").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("b", "Brand"))
		items, err := d.GetAllBrands(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Brand{value}, items)
		m.ExpectQuery("FROM brands").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("b", "Brand"))
		got, err := d.GetBrand(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM brands").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetBrand(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM brands").WillReturnError(boom)
		_, err = d.GetAllBrands(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO brands").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateBrand(ctx, value))
		m.ExpectExec("UPDATE brands").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateBrand(ctx, value))
		m.ExpectExec("DELETE FROM brands").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteBrand(ctx, value.Id))
		m.ExpectExec("UPDATE brands").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateBrand(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM brands").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteBrand(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO brands").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateBrand(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM brands").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteBrand(ctx, value.Id), domain.ErrForeignKey)
	})
}

func TestCategoryRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Category{Id: "c", Slug: "cat", Title: "Category", Description: "desc", Count: 3, Img: "cat.jpg"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM categories").WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "count", "img"}).AddRow("c", "cat", "Category", "desc", 3, "cat.jpg"))
		items, err := d.GetAllCategories(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Category{value}, items)
		m.ExpectQuery("FROM categories").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "slug", "title", "description", "count", "img"}).AddRow("c", "cat", "Category", "desc", 3, "cat.jpg"))
		got, err := d.GetCategory(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM categories").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetCategory(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM categories").WillReturnError(boom)
		_, err = d.GetAllCategories(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO categories").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateCategory(ctx, value))
		m.ExpectExec("UPDATE categories").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateCategory(ctx, value))
		m.ExpectExec("DELETE FROM categories").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteCategory(ctx, value.Id))
		m.ExpectExec("UPDATE categories").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateCategory(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM categories").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteCategory(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO categories").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateCategory(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM categories").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteCategory(ctx, value.Id), domain.ErrForeignKey)
	})
}

func TestCountryRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Country{Id: "co", Title: "Country", Friendly: "friendly"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM countries").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "friendly"}).AddRow("co", "Country", "friendly"))
		items, err := d.GetAllCountries(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Country{value}, items)
		m.ExpectQuery("FROM countries").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "friendly"}).AddRow("co", "Country", "friendly"))
		got, err := d.GetCountry(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM countries").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetCountry(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM countries").WillReturnError(boom)
		_, err = d.GetAllCountries(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO countries").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateCountry(ctx, value))
		m.ExpectExec("UPDATE countries").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateCountry(ctx, value))
		m.ExpectExec("DELETE FROM countries").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteCountry(ctx, value.Id))
		m.ExpectExec("UPDATE countries").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateCountry(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM countries").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteCountry(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO countries").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateCountry(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM countries").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteCountry(ctx, value.Id), domain.ErrForeignKey)
	})
}

func TestMaterialRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Material{Id: "m", Title: "Material"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM materials").WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("m", "Material"))
		items, err := d.GetAllMaterials(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Material{value}, items)
		m.ExpectQuery("FROM materials").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "title"}).AddRow("m", "Material"))
		got, err := d.GetMaterial(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM materials").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetMaterial(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM materials").WillReturnError(boom)
		_, err = d.GetAllMaterials(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO materials").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateMaterial(ctx, value))
		m.ExpectExec("UPDATE materials").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateMaterial(ctx, value))
		m.ExpectExec("DELETE FROM materials").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteMaterial(ctx, value.Id))
		m.ExpectExec("UPDATE materials").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateMaterial(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM materials").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteMaterial(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO materials").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateMaterial(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM materials").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteMaterial(ctx, value.Id), domain.ErrForeignKey)
	})
}

func TestColorRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Color{Id: "cl", Title: "Color", Hex: "#ffffff"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM colors").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "hex"}).AddRow("cl", "Color", "#ffffff"))
		items, err := d.GetAllColors(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Color{value}, items)
		m.ExpectQuery("FROM colors").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "title", "hex"}).AddRow("cl", "Color", "#ffffff"))
		got, err := d.GetColor(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM colors").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetColor(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM colors").WillReturnError(boom)
		_, err = d.GetAllColors(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO colors").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateColor(ctx, value))
		m.ExpectExec("UPDATE colors").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateColor(ctx, value))
		m.ExpectExec("DELETE FROM colors").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteColor(ctx, value.Id))
		m.ExpectExec("UPDATE colors").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateColor(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM colors").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteColor(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO colors").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateColor(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM colors").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteColor(ctx, value.Id), domain.ErrForeignKey)
	})
}

func TestSlideRepository(t *testing.T) {
	ctx := context.Background()
	value := &domain.Slide{Id: "s", Link: "/catalog", Img: "a.jpg", Img762: "b.jpg", Title: "Slide", Description: "desc", Badge: "sub"}
	t.Run("read list and record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM slides").WillReturnRows(sqlmock.NewRows([]string{"id", "link", "img", "img762", "title", "description", "badge"}).AddRow("s", "/catalog", "a.jpg", "b.jpg", "Slide", "desc", "sub"))
		items, err := d.GetAllSlides(ctx)
		require.NoError(t, err)
		require.Equal(t, []*domain.Slide{value}, items)
		m.ExpectQuery("FROM slides").WithArgs(value.Id).WillReturnRows(sqlmock.NewRows([]string{"id", "link", "img", "img762", "title", "description", "badge"}).AddRow("s", "/catalog", "a.jpg", "b.jpg", "Slide", "desc", "sub"))
		got, err := d.GetSlide(ctx, value.Id)
		require.NoError(t, err)
		require.Equal(t, value, got)
	})
	t.Run("not found and database error", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectQuery("FROM slides").WithArgs(value.Id).WillReturnError(sql.ErrNoRows)
		_, err := d.GetSlide(ctx, value.Id)
		require.ErrorIs(t, err, domain.ErrNotFound)
		boom := errors.New("disconnected")
		m.ExpectQuery("FROM slides").WillReturnError(boom)
		_, err = d.GetAllSlides(ctx)
		require.ErrorIs(t, err, boom)
	})
	t.Run("writes and missing record", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO slides").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.CreateSlide(ctx, value))
		m.ExpectExec("UPDATE slides").WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.UpdateSlide(ctx, value))
		m.ExpectExec("DELETE FROM slides").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 1))
		require.NoError(t, d.DeleteSlide(ctx, value.Id))
		m.ExpectExec("UPDATE slides").WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.UpdateSlide(ctx, value), domain.ErrNotFound)
		m.ExpectExec("DELETE FROM slides").WithArgs(value.Id).WillReturnResult(sqlmock.NewResult(0, 0))
		require.ErrorIs(t, d.DeleteSlide(ctx, value.Id), domain.ErrNotFound)
	})
	t.Run("unique and foreign key", func(t *testing.T) {
		d, m := mockDriver(t)
		m.ExpectExec("INSERT INTO slides").WillReturnError(&pq.Error{Code: "23505"})
		require.ErrorIs(t, d.CreateSlide(ctx, value), domain.ErrAlreadyExists)
		m.ExpectExec("DELETE FROM slides").WillReturnError(&pq.Error{Code: "23503"})
		require.ErrorIs(t, d.DeleteSlide(ctx, value.Id), domain.ErrForeignKey)
	})
}

func fullProductRows() *sqlmock.Rows {
	cols := make([]string, 26)
	for i := range cols {
		cols[i] = "column"
	}
	return sqlmock.NewRows(cols).AddRow("p", "Desk", "12345678", "desk", 10, 20, 30, `{"one.jpg","two.jpg"}`, 100, "Description", true, 90, "b", "Brand", "c", "cat", "Category", "Category desc", 2, "cat.jpg", "co", "Country", "friendly", `[{"id":"m","title":"Wood"}]`, `[{"id":"cl","title":"Black","hex":"#000000"}]`, `[]`)
}
func TestProductRepositoryReadsAndWrites(t *testing.T) {
	ctx := context.Background()
	d, m := mockDriver(t)
	m.ExpectQuery("WHERE p.id").WithArgs("p").WillReturnRows(fullProductRows())
	p, err := d.GetProduct(ctx, "p")
	require.NoError(t, err)
	require.Equal(t, "12345678", p.Article)
	require.Equal(t, []string{"one.jpg", "two.jpg"}, p.Photos)
	require.Equal(t, "Wood", p.Materials[0].Title)
	require.Equal(t, int32(2), p.Category.Count)
	require.NoError(t, decodeProductJSON(p, productJSONFields{
		materials: []byte(`[]`), colors: []byte(`[]`),
		similar: []byte(`[{"id":"similar","title":"Chair","article":"87654321","photos":["chair.jpg"],"isFeatured":true,"pricesale":20}]`),
	}))
	require.Equal(t, "87654321", p.Seems[0].Article)
	require.Equal(t, []string{"chair.jpg"}, p.Seems[0].Photos)
	require.Equal(t, int32(20), p.Seems[0].PriceSale)
	require.True(t, p.Seems[0].IsFeatured)
	m.ExpectQuery("COUNT").WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
	m.ExpectQuery("LIMIT").WithArgs(10, 20).WillReturnRows(fullProductRows())
	items, total, err := d.GetAllProducts(ctx, 20, 10)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
	m.ExpectQuery("WHERE p.id").WillReturnError(sql.ErrNoRows)
	_, err = d.GetProduct(ctx, "missing")
	require.ErrorIs(t, err, domain.ErrNotFound)
	value := &domain.ProductId{Id: "p", Title: "Desk", Article: "12345678", Slug: "desk"}
	m.ExpectExec("WITH").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.CreateProduct(ctx, value))
	m.ExpectQuery("WITH").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	require.NoError(t, d.UpdateProduct(ctx, value))
	m.ExpectExec("DELETE FROM products").WithArgs("p").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.DeleteProduct(ctx, "p"))
	m.ExpectExec("DELETE FROM products").WillReturnResult(sqlmock.NewResult(0, 0))
	require.ErrorIs(t, d.DeleteProduct(ctx, "p"), domain.ErrNotFound)
	m.ExpectExec("WITH").WillReturnError(&pq.Error{Code: "23503"})
	require.ErrorIs(t, d.CreateProduct(ctx, value), domain.ErrForeignKey)
}
func TestFilterRepositoryReturnsCategoryAndAllPhotos(t *testing.T) {
	d, m := mockDriver(t)
	filter := &domain.ProductFilter{Title: "Desk", Limit: 10, Offset: 20, Brand: []string{"b"}, MinPrice: 1, MaxPrice: 200, SortBy: "price", SortOrder: "desc"}
	m.ExpectQuery("COUNT").WillReturnRows(sqlmock.NewRows([]string{"total"}).AddRow(1))
	cols := make([]string, 15)
	for i := range cols {
		cols[i] = "col"
	}
	m.ExpectQuery("ORDER BY p.price DESC, p.id ASC").WillReturnRows(sqlmock.NewRows(cols).AddRow("p", "desk", "Category", "Desk", 100, "one.jpg", 90, true, "c", "cat", "Category", "desc", 2, "cat.jpg", `{"one.jpg","two.jpg"}`))
	items, total, err := d.FilterProducts(context.Background(), filter)
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Equal(t, "c", items[0].CategoryDetails.Id)
	require.Equal(t, []string{"one.jpg", "two.jpg"}, items[0].Photos)
	conditions, args := buildProductFilter(&domain.ProductFilter{Title: "' OR 1=1 --", Country: []string{"co"}, Category: []string{"c"}, Materials: []string{"m"}, Colors: []string{"cl"}, MinWidth: 1, MaxWidth: 10, MinHeight: 1, MaxHeight: 10, MinDepth: 1, MaxDepth: 10})
	require.NotEmpty(t, conditions)
	require.NotEmpty(t, args)
	for _, condition := range conditions {
		require.NotContains(t, condition, "OR 1=1")
	}
	require.Equal(t, " ORDER BY p.title ASC, p.id ASC", productOrderBy("injection", "desc"))
}
func TestDictionariesRepository(t *testing.T) {
	d, m := mockDriver(t)
	ctx := context.Background()
	rows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"type", "id", "a", "b", "c", "d", "e"}).AddRow("brand", "b", "Brand", "", "", "", "").AddRow("category", "c", "Category", "cat", "cat.jpg", "desc", "3").AddRow("country", "co", "Country", "friendly", "", "", "").AddRow("material", "m", "Wood", "", "", "", "").AddRow("color", "cl", "Black", "#000000", "", "", "").AddRow("stats", "", "10", "100", "1,10,2,20,3,30", "", "")
	}
	m.ExpectQuery("SELECT").WillReturnRows(rows())
	got, err := d.GetDictionaries(ctx)
	require.NoError(t, err)
	require.Equal(t, int32(10), got.MinPrice)
	require.Equal(t, int32(30), got.MaxDepth)
	require.Len(t, got.Brands, 1)
	require.Equal(t, int32(3), got.Categories[0].Count)
	m.ExpectQuery("SELECT").WithArgs("c").WillReturnRows(rows())
	got, err = d.GetDictionariesByCategoryID(ctx, "c")
	require.NoError(t, err)
	require.Len(t, got.Colors, 1)
	m.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"type", "id", "a", "b", "c", "d", "e"}).AddRow("category", "c", "Category", "cat", "x.jpg", "desc", "not a number"))
	_, err = d.GetDictionaries(ctx)
	require.Error(t, err)
}
func TestPhotoAndFileRepository(t *testing.T) {
	ctx := context.Background()
	d, m := mockDriver(t)
	p := &domain.ProductColorPhotos{ProductId: "p", ColorId: "c", Photos: []string{"a.jpg"}}
	m.ExpectQuery("SELECT photos").WithArgs("p", "c").WillReturnRows(sqlmock.NewRows([]string{"photos"}).AddRow(`{"a.jpg"}`))
	photos, err := d.GetPhotosByProductColor(ctx, "p", "c")
	require.NoError(t, err)
	require.Equal(t, p.Photos, photos)
	rows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"product_id", "color_id", "photos"}).AddRow("p", "c", `{"a.jpg"}`)
	}
	m.ExpectQuery("SELECT product_id").WillReturnRows(rows())
	all, err := d.GetAllProductColorPhotos(ctx)
	require.NoError(t, err)
	require.Equal(t, []domain.ProductColorPhotos{*p}, all)
	m.ExpectQuery("SELECT product_id").WithArgs("p").WillReturnRows(rows())
	all, err = d.GetProductColorPhotosByProduct(ctx, "p")
	require.NoError(t, err)
	require.Len(t, all, 1)
	m.ExpectExec("INSERT INTO product_color_photos").WithArgs("p", "c", pq.Array(p.Photos)).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.CreateProductColorPhotos(ctx, p))
	m.ExpectExec("UPDATE product_color_photos").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.UpdateProductColorPhotos(ctx, p))
	m.ExpectExec("DELETE FROM product_color_photos").WillReturnResult(sqlmock.NewResult(0, 0))
	require.ErrorIs(t, d.DeleteProductColorPhotos(ctx, "p", "c"), domain.ErrNotFound)
	m.ExpectExec("DELETE FROM product_color_photos").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.DeleteProductColorPhotos(ctx, "p", "c"))
	m.ExpectQuery("SELECT EXISTS").WithArgs("a.pdf").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	yes, err := d.IsFileReferenced(ctx, "a.pdf")
	require.NoError(t, err)
	require.True(t, yes)
}
func TestArticleRepositoryPreservesCreationTime(t *testing.T) {
	d, m := mockDriver(t)
	ctx := context.Background()
	created := time.UnixMilli(1720000000123).UTC()
	article := &domain.Article{Id: "a", Slug: "guide", Title: "Guide", Img: []string{"cover.jpg"}, Text: "/pdfs/guide.pdf"}
	m.ExpectQuery("INSERT INTO articles").WithArgs(article.Id, article.Slug, article.Title, pq.Array(article.Img), article.Text).WillReturnRows(sqlmock.NewRows([]string{"created_at"}).AddRow(created))
	require.NoError(t, d.CreateArticle(ctx, article))
	require.Equal(t, created.UnixMilli(), article.CreationTime)
	rows := func() *sqlmock.Rows {
		return sqlmock.NewRows([]string{"id", "slug", "title", "img", "pdf", "created_at"}).AddRow("a", "guide", "Guide", `{"cover.jpg"}`, "/pdfs/guide.pdf", created)
	}
	m.ExpectQuery("SELECT id").WithArgs("a").WillReturnRows(rows())
	got, err := d.GetArticle(ctx, "a")
	require.NoError(t, err)
	require.Equal(t, article, got)
	m.ExpectQuery("COUNT").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	m.ExpectQuery("ORDER BY created_at").WithArgs(10, 0).WillReturnRows(rows())
	items, total, err := d.GetAllArticles(ctx, 0, 10)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, 1, total)
	m.ExpectExec("UPDATE articles SET slug").WithArgs("a", "guide", "Guide", pq.Array(article.Img), article.Text).WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.UpdateArticle(ctx, article))
	m.ExpectExec("DELETE FROM articles").WithArgs("a").WillReturnResult(sqlmock.NewResult(0, 0))
	require.ErrorIs(t, d.DeleteArticle(ctx, "a"), domain.ErrNotFound)
	m.ExpectExec("DELETE FROM articles").WithArgs("a").WillReturnResult(sqlmock.NewResult(0, 1))
	require.NoError(t, d.DeleteArticle(ctx, "a"))
	m.ExpectQuery("SELECT id").WillReturnError(sql.ErrNoRows)
	_, err = d.GetArticle(ctx, "gone")
	require.ErrorIs(t, err, domain.ErrNotFound)
}
func TestSQLFailureMapping(t *testing.T) {
	for code, want := range map[string]error{"23503": domain.ErrForeignKey, "23505": domain.ErrAlreadyExists, "23514": domain.ErrInvalidType} {
		require.ErrorIs(t, mapPostgresError("test", &pq.Error{Code: pq.ErrorCode(code)}), want)
	}
	boom := errors.New("io")
	require.ErrorIs(t, mapPostgresError("test", boom), boom)
}
