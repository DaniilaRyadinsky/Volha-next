package api

import (
	"context"
	"errors"
	"testing"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/autumnterror/volha-backend/internal/product-service/config"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/mocks"
	"github.com/autumnterror/volha-backend/internal/product-service/service"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func TestBrandRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Brand{Id: id, Title: "Brand"}
	repo := mocks.CatalogRepo{GetAllBrandsFn: func(context.Context) ([]*domain.Brand, error) { return []*domain.Brand{value}, nil }, GetBrandFn: func(context.Context, string) (*domain.Brand, error) { return value, nil }, CreateBrandFn: func(context.Context, *domain.Brand) error { return nil }, UpdateBrandFn: func(context.Context, *domain.Brand) error { return nil }, DeleteBrandFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllBrands(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetBrand(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateBrand(ctx, domain.BrandToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateBrand(ctx, domain.BrandToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteBrand(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetBrand(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCategoryRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Category{Id: id, Title: "Category", Slug: "cat"}
	repo := mocks.CatalogRepo{GetAllCategoriesFn: func(context.Context) ([]*domain.Category, error) { return []*domain.Category{value}, nil }, GetCategoryFn: func(context.Context, string) (*domain.Category, error) { return value, nil }, CreateCategoryFn: func(context.Context, *domain.Category) error { return nil }, UpdateCategoryFn: func(context.Context, *domain.Category) error { return nil }, DeleteCategoryFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllCategories(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetCategory(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateCategory(ctx, domain.CategoryToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateCategory(ctx, domain.CategoryToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteCategory(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetCategory(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestCountryRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Country{Id: id, Title: "Country"}
	repo := mocks.CatalogRepo{GetAllCountriesFn: func(context.Context) ([]*domain.Country, error) { return []*domain.Country{value}, nil }, GetCountryFn: func(context.Context, string) (*domain.Country, error) { return value, nil }, CreateCountryFn: func(context.Context, *domain.Country) error { return nil }, UpdateCountryFn: func(context.Context, *domain.Country) error { return nil }, DeleteCountryFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllCountries(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetCountry(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateCountry(ctx, domain.CountryToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateCountry(ctx, domain.CountryToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteCountry(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetCountry(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestMaterialRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Material{Id: id, Title: "Wood"}
	repo := mocks.CatalogRepo{GetAllMaterialsFn: func(context.Context) ([]*domain.Material, error) { return []*domain.Material{value}, nil }, GetMaterialFn: func(context.Context, string) (*domain.Material, error) { return value, nil }, CreateMaterialFn: func(context.Context, *domain.Material) error { return nil }, UpdateMaterialFn: func(context.Context, *domain.Material) error { return nil }, DeleteMaterialFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllMaterials(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetMaterial(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateMaterial(ctx, domain.MaterialToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateMaterial(ctx, domain.MaterialToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteMaterial(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetMaterial(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestColorRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Color{Id: id, Title: "Color", Hex: "#000000"}
	repo := mocks.CatalogRepo{GetAllColorsFn: func(context.Context) ([]*domain.Color, error) { return []*domain.Color{value}, nil }, GetColorFn: func(context.Context, string) (*domain.Color, error) { return value, nil }, CreateColorFn: func(context.Context, *domain.Color) error { return nil }, UpdateColorFn: func(context.Context, *domain.Color) error { return nil }, DeleteColorFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllColors(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetColor(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateColor(ctx, domain.ColorToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateColor(ctx, domain.ColorToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteColor(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetColor(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestSlideRPC(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	value := &domain.Slide{Id: id, Title: "Slide", Link: "/", Img: "a.jpg", Img762: "b.jpg"}
	repo := mocks.CatalogRepo{GetAllSlidesFn: func(context.Context) ([]*domain.Slide, error) { return []*domain.Slide{value}, nil }, GetSlideFn: func(context.Context, string) (*domain.Slide, error) { return value, nil }, CreateSlideFn: func(context.Context, *domain.Slide) error { return nil }, UpdateSlideFn: func(context.Context, *domain.Slide) error { return nil }, DeleteSlideFn: func(context.Context, string) error { return domain.ErrNotFound }}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})}
	all, err := api.GetAllSlides(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	require.Equal(t, id, all.GetItems()[0].GetId())
	got, err := api.GetSlide(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, id, got.GetId())
	_, err = api.CreateSlide(ctx, domain.SlideToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateSlide(ctx, domain.SlideToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteSlide(ctx, &productsRPC.Id{Id: id})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = api.GetSlide(ctx, &productsRPC.Id{Id: "invalid"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func TestProductPhotoDictionaryArticleRPCs(t *testing.T) {
	ctx := context.Background()
	id := xid.New().String()
	p := &domain.Product{Id: id, Title: "Desk", Materials: []*domain.Material{}, Colors: []*domain.Color{}, Photos: []string{}, Seems: []*domain.Product{}}
	value := &domain.ProductId{Id: id, Title: "Desk", Slug: "desk", Article: "12345678", Brand: id, Category: id, Country: id, Width: 1, Height: 1, Depth: 1, Price: 10, Photos: []string{"a.jpg"}}
	photo := domain.ProductColorPhotos{ProductId: id, ColorId: id, Photos: []string{"a.jpg"}}
	article := &domain.Article{Id: id, Title: "Guide", Slug: "guide", Text: "/pdfs/guide.pdf", Img: []string{}, CreationTime: 123}
	provider := mocks.Provider{ProductRepo: mocks.ProductRepo{GetAllProductsFn: func(context.Context, int, int) ([]*domain.Product, int, error) { return []*domain.Product{p}, 1, nil }, GetProductFn: func(context.Context, string) (*domain.Product, error) { return p, nil }, CreateProductFn: func(context.Context, *domain.ProductId) error { return nil }, UpdateProductFn: func(context.Context, *domain.ProductId) error { return nil }, DeleteProductFn: func(context.Context, string) error { return nil }, FilterProductsFn: func(context.Context, *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error) {
		return []*domain.ProductCatalogItem{{Id: id, CategoryDetails: &domain.Category{Id: id}, Photos: []string{"a.jpg"}}}, 1, nil
	}},
		PcpRepo: mocks.PcpRepo{GetPhotosByProductColorFn: func(context.Context, string, string) ([]string, error) { return photo.Photos, nil }, GetAllProductColorPhotosFn: func(context.Context) ([]domain.ProductColorPhotos, error) {
			return []domain.ProductColorPhotos{photo}, nil
		}, GetProductColorPhotosByProductFn: func(context.Context, string) ([]domain.ProductColorPhotos, error) {
			return []domain.ProductColorPhotos{photo}, nil
		}, CreateProductColorPhotosFn: func(context.Context, *domain.ProductColorPhotos) error { return nil }, UpdateProductColorPhotosFn: func(context.Context, *domain.ProductColorPhotos) error { return nil }, DeleteProductColorPhotosFn: func(context.Context, string, string) error { return nil }},
		DictRepo: mocks.DictRepo{GetDictionariesFn: func(context.Context) (*domain.Dictionaries, error) { return &domain.Dictionaries{MinPrice: 1}, nil }}, FileRepo: mocks.FileRepo{IsFileReferencedFn: func(context.Context, string) (bool, error) { return true, nil }},
		ArticleRepo: mocks.ArticleRepo{GetAllArticlesFn: func(context.Context, int, int) ([]*domain.Article, int, error) {
			return []*domain.Article{article}, 1, nil
		}, GetArticleFn: func(context.Context, string) (*domain.Article, error) { return article, nil }, CreateArticleFn: func(context.Context, *domain.Article) error { return nil }, UpdateArticleFn: func(context.Context, *domain.Article) error { return nil }, DeleteArticleFn: func(context.Context, string) error { return nil }}}
	api := &ServerAPI{s: service.NewProductsService(mocks.TxRunner{}, provider)}
	list, err := api.GetAllProducts(ctx, &productsRPC.Pagination{Offset: 10, Limit: 20})
	require.NoError(t, err)
	require.Equal(t, int32(10), list.GetOffset())
	require.Len(t, list.GetItems(), 1)
	got, err := api.GetProduct(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, "Desk", got.GetTitle())
	_, err = api.CreateProduct(ctx, domain.ProductIdToRpc(value))
	require.NoError(t, err)
	_, err = api.UpdateProduct(ctx, domain.ProductIdToRpc(value))
	require.NoError(t, err)
	_, err = api.DeleteProduct(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	filtered, err := api.FilterProducts(ctx, &productsRPC.ProductFilter{Limit: 20})
	require.NoError(t, err)
	require.Equal(t, id, filtered.GetItems()[0].GetCategoryDetails().GetId())
	require.Equal(t, []string{"a.jpg"}, filtered.GetItems()[0].GetPhotos())
	dict, err := api.GetDictionaries(ctx, &productsRPC.Id{Id: "CLR"})
	require.NoError(t, err)
	require.Equal(t, int32(1), dict.GetMinPrice())
	refs, err := api.IsFileReferenced(ctx, &productsRPC.FileReferenceRequest{Filename: "a.pdf"})
	require.NoError(t, err)
	require.True(t, refs.GetReferenced())
	_, err = api.CreateProductColorPhotos(ctx, domain.ProductColorPhotosToRpc(&photo))
	require.NoError(t, err)
	_, err = api.UpdateProductColorPhotos(ctx, domain.ProductColorPhotosToRpc(&photo))
	require.NoError(t, err)
	_, err = api.DeleteProductColorPhotos(ctx, &productsRPC.ProductColorPhotosId{ProductId: id, ColorId: id})
	require.NoError(t, err)
	photos, err := api.GetPhotosByProductAndColor(ctx, &productsRPC.ProductColorPhotosId{ProductId: id, ColorId: id})
	require.NoError(t, err)
	require.Equal(t, photo.Photos, photos.GetItems())
	all, err := api.GetAllProductColorPhotos(ctx, nil)
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	all, err = api.GetProductColorPhotosByProduct(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Len(t, all.GetItems(), 1)
	articles, err := api.GetAllArticles(ctx, &productsRPC.Pagination{Limit: 10})
	require.NoError(t, err)
	require.Len(t, articles.GetItems(), 1)
	a, err := api.GetArticle(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
	require.Equal(t, int64(123), a.GetCreationTime())
	_, err = api.CreateArticle(ctx, domain.ArticleToRpc(article))
	require.NoError(t, err)
	_, err = api.UpdateArticle(ctx, domain.ArticleToRpc(article))
	require.NoError(t, err)
	_, err = api.DeleteArticle(ctx, &productsRPC.Id{Id: id})
	require.NoError(t, err)
}
func TestRPCErrorMappingAndHealth(t *testing.T) {
	for input, want := range map[error]codes.Code{domain.ErrUnknownType: codes.Unimplemented, domain.ErrInvalidType: codes.InvalidArgument, domain.ErrNotFound: codes.NotFound, domain.ErrAlreadyExists: codes.AlreadyExists, domain.ErrForeignKey: codes.FailedPrecondition, service.ErrBadServiceCheck: codes.InvalidArgument, errors.New("secret db detail"): codes.Internal} {
		_, err := handleCRUDResponse(context.Background(), "test", func() (any, error) { return nil, input })
		require.Equal(t, want, status.Code(err))
		if want == codes.Internal {
			require.NotContains(t, err.Error(), "secret db detail")
		}
	}
	good, err := (dependencyHealthServer{}).Check(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, healthpb.HealthCheckResponse_SERVING, good.GetStatus())
	_, err = (dependencyHealthServer{ready: func(context.Context) error { return errors.New("down") }}).Check(context.Background(), nil)
	require.Equal(t, codes.Unavailable, status.Code(err))
	app := New(&config.Config{Port: -1}, service.NewProductsService(nil, nil), nil, observability.Config{}, observability.NewMetrics("test", "dev", "test"))
	require.Error(t, app.Run())
	app.Stop(context.Background())
}
