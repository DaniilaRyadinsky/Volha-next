package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/mocks"
	"github.com/autumnterror/volha-backend/internal/product-service/service"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
)

func TestBrandServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Brand{Id: id, Title: "Brand"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllBrandsFn: func(context.Context) ([]*domain.Brand, error) { calls++; return []*domain.Brand{value}, repoErr },
				GetBrandFn: func(_ context.Context, got string) (*domain.Brand, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateBrandFn: func(_ context.Context, got *domain.Brand) error { require.Same(t, value, got); calls++; return repoErr },
				UpdateBrandFn: func(_ context.Context, got *domain.Brand) error { require.Same(t, value, got); calls++; return repoErr },
				DeleteBrandFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllBrands(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Brand{value}, items)
			}
			got, err := svc.GetBrand(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateBrand(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateBrand(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteBrand(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetBrand(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateBrand(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateBrand(ctx, &domain.Brand{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteBrand(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllBrands(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestCategoryServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Category{Id: id, Title: "Category", Slug: "category", Img: "photo.jpg"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllCategoriesFn: func(context.Context) ([]*domain.Category, error) { calls++; return []*domain.Category{value}, repoErr },
				GetCategoryFn: func(_ context.Context, got string) (*domain.Category, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateCategoryFn: func(_ context.Context, got *domain.Category) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				UpdateCategoryFn: func(_ context.Context, got *domain.Category) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				DeleteCategoryFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllCategories(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Category{value}, items)
			}
			got, err := svc.GetCategory(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateCategory(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateCategory(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteCategory(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetCategory(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateCategory(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateCategory(ctx, &domain.Category{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteCategory(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllCategories(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestCountryServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Country{Id: id, Title: "Country"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllCountriesFn: func(context.Context) ([]*domain.Country, error) { calls++; return []*domain.Country{value}, repoErr },
				GetCountryFn: func(_ context.Context, got string) (*domain.Country, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateCountryFn: func(_ context.Context, got *domain.Country) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				UpdateCountryFn: func(_ context.Context, got *domain.Country) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				DeleteCountryFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllCountries(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Country{value}, items)
			}
			got, err := svc.GetCountry(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateCountry(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateCountry(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteCountry(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetCountry(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateCountry(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateCountry(ctx, &domain.Country{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteCountry(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllCountries(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestMaterialServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Material{Id: id, Title: "Material"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllMaterialsFn: func(context.Context) ([]*domain.Material, error) { calls++; return []*domain.Material{value}, repoErr },
				GetMaterialFn: func(_ context.Context, got string) (*domain.Material, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateMaterialFn: func(_ context.Context, got *domain.Material) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				UpdateMaterialFn: func(_ context.Context, got *domain.Material) error {
					require.Same(t, value, got)
					calls++
					return repoErr
				},
				DeleteMaterialFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllMaterials(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Material{value}, items)
			}
			got, err := svc.GetMaterial(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateMaterial(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateMaterial(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteMaterial(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetMaterial(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateMaterial(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateMaterial(ctx, &domain.Material{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteMaterial(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllMaterials(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestColorServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Color{Id: id, Title: "Color", Hex: "#000000"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllColorsFn: func(context.Context) ([]*domain.Color, error) { calls++; return []*domain.Color{value}, repoErr },
				GetColorFn: func(_ context.Context, got string) (*domain.Color, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateColorFn: func(_ context.Context, got *domain.Color) error { require.Same(t, value, got); calls++; return repoErr },
				UpdateColorFn: func(_ context.Context, got *domain.Color) error { require.Same(t, value, got); calls++; return repoErr },
				DeleteColorFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllColors(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Color{value}, items)
			}
			got, err := svc.GetColor(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateColor(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateColor(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteColor(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetColor(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateColor(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateColor(ctx, &domain.Color{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteColor(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllColors(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestSlideServiceOperations(t *testing.T) {
	id := xid.New().String()
	value := &domain.Slide{Id: id, Title: "Slide", Link: "/catalog", Img: "one.jpg", Img762: "two.jpg"}
	ctx := context.Background()
	boom := errors.New("database unavailable")
	for _, failure := range []bool{false, true} {
		t.Run(map[bool]string{false: "success", true: "repository error"}[failure], func(t *testing.T) {
			var repoErr error
			if failure {
				repoErr = boom
			}
			calls := 0
			repo := mocks.CatalogRepo{
				GetAllSlidesFn: func(context.Context) ([]*domain.Slide, error) { calls++; return []*domain.Slide{value}, repoErr },
				GetSlideFn: func(_ context.Context, got string) (*domain.Slide, error) {
					require.Equal(t, id, got)
					calls++
					return value, repoErr
				},
				CreateSlideFn: func(_ context.Context, got *domain.Slide) error { require.Same(t, value, got); calls++; return repoErr },
				UpdateSlideFn: func(_ context.Context, got *domain.Slide) error { require.Same(t, value, got); calls++; return repoErr },
				DeleteSlideFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return repoErr }}
			svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})
			items, err := svc.GetAllSlides(ctx)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Equal(t, []*domain.Slide{value}, items)
			}
			got, err := svc.GetSlide(ctx, id)
			require.ErrorIs(t, err, repoErr)
			if !failure {
				require.Same(t, value, got)
			}
			require.ErrorIs(t, svc.CreateSlide(ctx, value), repoErr)
			require.ErrorIs(t, svc.UpdateSlide(ctx, value), repoErr)
			require.ErrorIs(t, svc.DeleteSlide(ctx, id), repoErr)
			require.Equal(t, 5, calls)
		})
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{})
	_, err := svc.GetSlide(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.CreateSlide(ctx, nil), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.UpdateSlide(ctx, &domain.Slide{Id: id}), service.ErrBadServiceCheck)
	require.ErrorIs(t, svc.DeleteSlide(ctx, "bad"), service.ErrBadServiceCheck)
	_, err = svc.GetAllSlides(ctx)
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
}

func TestProductWritesAndFilters(t *testing.T) {
	id := xid.New().String()
	other := xid.New().String()
	ctx := context.Background()
	p := &domain.ProductId{Id: id, Title: "Desk", Slug: "desk", Article: "12345678", Brand: other, Category: other, Country: other, Width: 1, Height: 2, Depth: 3, Price: 100, Photos: []string{"one.jpg"}}
	calls := 0
	repo := mocks.ProductRepo{
		CreateProductFn: func(_ context.Context, got *domain.ProductId) error { require.Equal(t, p, got); calls++; return nil },
		UpdateProductFn: func(_ context.Context, got *domain.ProductId) error { calls++; return domain.ErrNotFound },
		DeleteProductFn: func(_ context.Context, got string) error { require.Equal(t, id, got); calls++; return nil },
		FilterProductsFn: func(_ context.Context, f *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error) {
			require.Equal(t, int32(20), f.Limit)
			calls++
			return []*domain.ProductCatalogItem{{Id: id}}, 1, nil
		}}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{ProductRepo: repo})
	require.NoError(t, svc.CreateProduct(ctx, p))
	require.ErrorIs(t, svc.UpdateProduct(ctx, p), domain.ErrNotFound)
	require.NoError(t, svc.DeleteProduct(ctx, id))
	items, total, err := svc.FilterProducts(ctx, &domain.ProductFilter{Limit: 20})
	require.NoError(t, err)
	require.Equal(t, 1, total)
	require.Len(t, items, 1)
	require.Equal(t, 4, calls)
	mutations := []func(*domain.ProductId){func(p *domain.ProductId) { p.Id = "bad" }, func(p *domain.ProductId) { p.Title = "" }, func(p *domain.ProductId) { p.Slug = "" }, func(p *domain.ProductId) { p.Article = "wrong" }, func(p *domain.ProductId) { p.Brand = "bad" }, func(p *domain.ProductId) { p.Category = "bad" }, func(p *domain.ProductId) { p.Country = "bad" }, func(p *domain.ProductId) { p.Width = 0 }, func(p *domain.ProductId) { p.Price = 0 }, func(p *domain.ProductId) { p.PriceSale = -1 }, func(p *domain.ProductId) { p.Photos = nil }, func(p *domain.ProductId) { p.Materials = []string{"bad"} }, func(p *domain.ProductId) { p.Colors = []string{other, other} }, func(p *domain.ProductId) { p.Seems = []string{id} }}
	for _, change := range mutations {
		copy := *p
		change(&copy)
		require.ErrorIs(t, svc.CreateProduct(ctx, &copy), service.ErrBadServiceCheck)
	}
	for _, f := range []*domain.ProductFilter{nil, {Limit: 101}, {Limit: 1, Offset: -1}, {Limit: 1, MinPrice: 10, MaxPrice: 1}, {Limit: 1, SortBy: "DROP TABLE products"}, {Limit: 1, SortOrder: "sideways"}, {Limit: 1, Brand: []string{"bad"}}} {
		_, _, err := svc.FilterProducts(ctx, f)
		require.ErrorIs(t, err, service.ErrBadServiceCheck)
	}
}
func TestDictionariesAndFileReferences(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Dictionaries{MinPrice: 10}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: mocks.CatalogRepo{GetCategoryFn: func(context.Context, string) (*domain.Category, error) { return &domain.Category{Id: id}, nil }}, DictRepo: mocks.DictRepo{GetDictionariesFn: func(context.Context) (*domain.Dictionaries, error) { return value, nil }, GetDictionariesByCategoryIDFn: func(_ context.Context, got string) (*domain.Dictionaries, error) {
		require.Equal(t, id, got)
		return value, nil
	}}, FileRepo: mocks.FileRepo{IsFileReferencedFn: func(_ context.Context, name string) (bool, error) {
		require.Equal(t, "photo.jpg", name)
		return true, nil
	}}})
	for _, key := range []string{"", "CLR", id} {
		got, err := svc.GetDictionaries(ctx, key)
		require.NoError(t, err)
		require.Same(t, value, got)
	}
	_, err := svc.GetDictionaries(ctx, "bad")
	require.ErrorIs(t, err, service.ErrBadServiceCheck)
	yes, err := svc.IsFileReferenced(ctx, "photo.jpg")
	require.NoError(t, err)
	require.True(t, yes)
	for _, name := range []string{"", "../escape.jpg", "dir/name.jpg"} {
		_, err = svc.IsFileReferenced(ctx, name)
		require.ErrorIs(t, err, service.ErrBadServiceCheck)
	}
}
func TestPhotoService(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	color := xid.New().String()
	p := &domain.ProductColorPhotos{ProductId: id, ColorId: color, Photos: []string{"a.jpg"}}
	repo := mocks.PcpRepo{
		GetPhotosByProductColorFn: func(_ context.Context, pid, cid string) ([]string, error) {
			require.Equal(t, id, pid)
			require.Equal(t, color, cid)
			return p.Photos, nil
		},
		GetAllProductColorPhotosFn: func(context.Context) ([]domain.ProductColorPhotos, error) {
			return []domain.ProductColorPhotos{*p}, nil
		},
		GetProductColorPhotosByProductFn: func(context.Context, string) ([]domain.ProductColorPhotos, error) {
			return []domain.ProductColorPhotos{*p}, nil
		},
		CreateProductColorPhotosFn: func(context.Context, *domain.ProductColorPhotos) error { return nil }, UpdateProductColorPhotosFn: func(context.Context, *domain.ProductColorPhotos) error { return domain.ErrNotFound }, DeleteProductColorPhotosFn: func(context.Context, string, string) error { return nil }}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{PcpRepo: repo})
	photos, err := svc.GetPhotosByProductColor(ctx, id, color)
	require.NoError(t, err)
	require.Equal(t, p.Photos, photos)
	all, err := svc.GetAllProductColorPhotos(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	all, err = svc.GetProductColorPhotosByProduct(ctx, id)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.NoError(t, svc.CreateProductColorPhotos(ctx, p))
	require.ErrorIs(t, svc.UpdateProductColorPhotos(ctx, p), domain.ErrNotFound)
	require.NoError(t, svc.DeleteProductColorPhotos(ctx, id, color))
	require.Error(t, svc.CreateProductColorPhotos(ctx, nil))
	require.Error(t, svc.UpdateProductColorPhotos(ctx, &domain.ProductColorPhotos{ProductId: id, ColorId: color}))
	require.Error(t, svc.DeleteProductColorPhotos(ctx, "bad", color))
	_, err = svc.GetProductColorPhotosByProduct(ctx, "bad")
	require.Error(t, err)
	_, err = svc.GetPhotosByProductColor(ctx, id, "bad")
	require.Error(t, err)
}
func TestArticleService(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	article := &domain.Article{Id: id, Slug: "guide", Title: "Guide", Text: "/pdfs/guide.pdf", Img: []string{}}
	repo := mocks.ArticleRepo{GetAllArticlesFn: func(_ context.Context, offset, limit int) ([]*domain.Article, int, error) {
		require.Equal(t, 20, offset)
		require.Equal(t, 10, limit)
		return []*domain.Article{article}, 1, nil
	}, GetArticleFn: func(context.Context, string) (*domain.Article, error) { return article, nil }, CreateArticleFn: func(context.Context, *domain.Article) error { return nil }, UpdateArticleFn: func(context.Context, *domain.Article) error { return domain.ErrNotFound }, DeleteArticleFn: func(context.Context, string) error { return nil }}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{ArticleRepo: repo})
	items, total, err := svc.GetAllArticles(ctx, 20, 10)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, 1, total)
	got, err := svc.GetArticle(ctx, id)
	require.NoError(t, err)
	require.Same(t, article, got)
	require.NoError(t, svc.CreateArticle(ctx, article))
	require.ErrorIs(t, svc.UpdateArticle(ctx, article), domain.ErrNotFound)
	require.NoError(t, svc.DeleteArticle(ctx, id))
	for _, link := range []string{"", "https://other/x.pdf", "/pdfs/../secret.pdf", "/pdfs/name.txt"} {
		copy := *article
		copy.Text = link
		require.ErrorIs(t, svc.CreateArticle(ctx, &copy), service.ErrBadServiceCheck)
	}
	require.Error(t, svc.CreateArticle(ctx, nil))
	require.Error(t, svc.DeleteArticle(ctx, "bad"))
	_, err = svc.GetArticle(ctx, "bad")
	require.Error(t, err)
	_, _, err = svc.GetAllArticles(ctx, -1, 10)
	require.Error(t, err)
}
