package mocks

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

type CatalogRepo struct {
	GetAllBrandsFn     func(ctx context.Context) ([]*domain.Brand, error)
	GetBrandFn         func(ctx context.Context, id string) (*domain.Brand, error)
	CreateBrandFn      func(ctx context.Context, v *domain.Brand) error
	UpdateBrandFn      func(ctx context.Context, v *domain.Brand) error
	DeleteBrandFn      func(ctx context.Context, id string) error
	GetAllCategoriesFn func(ctx context.Context) ([]*domain.Category, error)
	GetCategoryFn      func(ctx context.Context, id string) (*domain.Category, error)
	CreateCategoryFn   func(ctx context.Context, v *domain.Category) error
	UpdateCategoryFn   func(ctx context.Context, v *domain.Category) error
	DeleteCategoryFn   func(ctx context.Context, id string) error
	GetAllCountriesFn  func(ctx context.Context) ([]*domain.Country, error)
	GetCountryFn       func(ctx context.Context, id string) (*domain.Country, error)
	CreateCountryFn    func(ctx context.Context, v *domain.Country) error
	UpdateCountryFn    func(ctx context.Context, v *domain.Country) error
	DeleteCountryFn    func(ctx context.Context, id string) error
	GetAllMaterialsFn  func(ctx context.Context) ([]*domain.Material, error)
	GetMaterialFn      func(ctx context.Context, id string) (*domain.Material, error)
	CreateMaterialFn   func(ctx context.Context, v *domain.Material) error
	UpdateMaterialFn   func(ctx context.Context, v *domain.Material) error
	DeleteMaterialFn   func(ctx context.Context, id string) error
	GetAllColorsFn     func(ctx context.Context) ([]*domain.Color, error)
	GetColorFn         func(ctx context.Context, id string) (*domain.Color, error)
	CreateColorFn      func(ctx context.Context, v *domain.Color) error
	UpdateColorFn      func(ctx context.Context, v *domain.Color) error
	DeleteColorFn      func(ctx context.Context, id string) error
	GetAllSlidesFn     func(ctx context.Context) ([]*domain.Slide, error)
	GetSlideFn         func(ctx context.Context, id string) (*domain.Slide, error)
	CreateSlideFn      func(ctx context.Context, v *domain.Slide) error
	UpdateSlideFn      func(ctx context.Context, v *domain.Slide) error
	DeleteSlideFn      func(ctx context.Context, id string) error
}

func (s CatalogRepo) GetAllBrands(ctx context.Context) ([]*domain.Brand, error) {
	return s.GetAllBrandsFn(ctx)
}
func (s CatalogRepo) GetBrand(ctx context.Context, id string) (*domain.Brand, error) {
	return s.GetBrandFn(ctx, id)
}
func (s CatalogRepo) CreateBrand(ctx context.Context, v *domain.Brand) error {
	return s.CreateBrandFn(ctx, v)
}
func (s CatalogRepo) UpdateBrand(ctx context.Context, v *domain.Brand) error {
	return s.UpdateBrandFn(ctx, v)
}
func (s CatalogRepo) DeleteBrand(ctx context.Context, id string) error {
	return s.DeleteBrandFn(ctx, id)
}
func (s CatalogRepo) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	return s.GetAllCategoriesFn(ctx)
}
func (s CatalogRepo) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	return s.GetCategoryFn(ctx, id)
}
func (s CatalogRepo) CreateCategory(ctx context.Context, v *domain.Category) error {
	return s.CreateCategoryFn(ctx, v)
}
func (s CatalogRepo) UpdateCategory(ctx context.Context, v *domain.Category) error {
	return s.UpdateCategoryFn(ctx, v)
}
func (s CatalogRepo) DeleteCategory(ctx context.Context, id string) error {
	return s.DeleteCategoryFn(ctx, id)
}
func (s CatalogRepo) GetAllCountries(ctx context.Context) ([]*domain.Country, error) {
	return s.GetAllCountriesFn(ctx)
}
func (s CatalogRepo) GetCountry(ctx context.Context, id string) (*domain.Country, error) {
	return s.GetCountryFn(ctx, id)
}
func (s CatalogRepo) CreateCountry(ctx context.Context, v *domain.Country) error {
	return s.CreateCountryFn(ctx, v)
}
func (s CatalogRepo) UpdateCountry(ctx context.Context, v *domain.Country) error {
	return s.UpdateCountryFn(ctx, v)
}
func (s CatalogRepo) DeleteCountry(ctx context.Context, id string) error {
	return s.DeleteCountryFn(ctx, id)
}
func (s CatalogRepo) GetAllMaterials(ctx context.Context) ([]*domain.Material, error) {
	return s.GetAllMaterialsFn(ctx)
}
func (s CatalogRepo) GetMaterial(ctx context.Context, id string) (*domain.Material, error) {
	return s.GetMaterialFn(ctx, id)
}
func (s CatalogRepo) CreateMaterial(ctx context.Context, v *domain.Material) error {
	return s.CreateMaterialFn(ctx, v)
}
func (s CatalogRepo) UpdateMaterial(ctx context.Context, v *domain.Material) error {
	return s.UpdateMaterialFn(ctx, v)
}
func (s CatalogRepo) DeleteMaterial(ctx context.Context, id string) error {
	return s.DeleteMaterialFn(ctx, id)
}
func (s CatalogRepo) GetAllColors(ctx context.Context) ([]*domain.Color, error) {
	return s.GetAllColorsFn(ctx)
}
func (s CatalogRepo) GetColor(ctx context.Context, id string) (*domain.Color, error) {
	return s.GetColorFn(ctx, id)
}
func (s CatalogRepo) CreateColor(ctx context.Context, v *domain.Color) error {
	return s.CreateColorFn(ctx, v)
}
func (s CatalogRepo) UpdateColor(ctx context.Context, v *domain.Color) error {
	return s.UpdateColorFn(ctx, v)
}
func (s CatalogRepo) DeleteColor(ctx context.Context, id string) error {
	return s.DeleteColorFn(ctx, id)
}
func (s CatalogRepo) GetAllSlides(ctx context.Context) ([]*domain.Slide, error) {
	return s.GetAllSlidesFn(ctx)
}
func (s CatalogRepo) GetSlide(ctx context.Context, id string) (*domain.Slide, error) {
	return s.GetSlideFn(ctx, id)
}
func (s CatalogRepo) CreateSlide(ctx context.Context, v *domain.Slide) error {
	return s.CreateSlideFn(ctx, v)
}
func (s CatalogRepo) UpdateSlide(ctx context.Context, v *domain.Slide) error {
	return s.UpdateSlideFn(ctx, v)
}
func (s CatalogRepo) DeleteSlide(ctx context.Context, id string) error {
	return s.DeleteSlideFn(ctx, id)
}

type DictRepo struct {
	GetDictionariesFn             func(ctx context.Context) (*domain.Dictionaries, error)
	GetDictionariesByCategoryIDFn func(ctx context.Context, id string) (*domain.Dictionaries, error)
}

func (s DictRepo) GetDictionaries(ctx context.Context) (*domain.Dictionaries, error) {
	return s.GetDictionariesFn(ctx)
}
func (s DictRepo) GetDictionariesByCategoryID(ctx context.Context, id string) (*domain.Dictionaries, error) {
	return s.GetDictionariesByCategoryIDFn(ctx, id)
}

type PcpRepo struct {
	GetPhotosByProductColorFn        func(ctx context.Context, productID string, colorID string) ([]string, error)
	GetAllProductColorPhotosFn       func(ctx context.Context) ([]domain.ProductColorPhotos, error)
	GetProductColorPhotosByProductFn func(ctx context.Context, id string) ([]domain.ProductColorPhotos, error)
	CreateProductColorPhotosFn       func(ctx context.Context, v *domain.ProductColorPhotos) error
	UpdateProductColorPhotosFn       func(ctx context.Context, v *domain.ProductColorPhotos) error
	DeleteProductColorPhotosFn       func(ctx context.Context, productID string, colorID string) error
}

func (s PcpRepo) GetPhotosByProductColor(ctx context.Context, productID string, colorID string) ([]string, error) {
	return s.GetPhotosByProductColorFn(ctx, productID, colorID)
}
func (s PcpRepo) GetAllProductColorPhotos(ctx context.Context) ([]domain.ProductColorPhotos, error) {
	return s.GetAllProductColorPhotosFn(ctx)
}
func (s PcpRepo) GetProductColorPhotosByProduct(ctx context.Context, id string) ([]domain.ProductColorPhotos, error) {
	return s.GetProductColorPhotosByProductFn(ctx, id)
}
func (s PcpRepo) CreateProductColorPhotos(ctx context.Context, v *domain.ProductColorPhotos) error {
	return s.CreateProductColorPhotosFn(ctx, v)
}
func (s PcpRepo) UpdateProductColorPhotos(ctx context.Context, v *domain.ProductColorPhotos) error {
	return s.UpdateProductColorPhotosFn(ctx, v)
}
func (s PcpRepo) DeleteProductColorPhotos(ctx context.Context, productID string, colorID string) error {
	return s.DeleteProductColorPhotosFn(ctx, productID, colorID)
}

type FileRepo struct {
	IsFileReferencedFn func(ctx context.Context, filename string) (bool, error)
}

func (s FileRepo) IsFileReferenced(ctx context.Context, filename string) (bool, error) {
	return s.IsFileReferencedFn(ctx, filename)
}

type ArticleRepo struct {
	GetAllArticlesFn func(ctx context.Context, offset int, limit int) ([]*domain.Article, int, error)
	GetArticleFn     func(ctx context.Context, id string) (*domain.Article, error)
	CreateArticleFn  func(ctx context.Context, v *domain.Article) error
	UpdateArticleFn  func(ctx context.Context, v *domain.Article) error
	DeleteArticleFn  func(ctx context.Context, id string) error
}

func (s ArticleRepo) GetAllArticles(ctx context.Context, offset int, limit int) ([]*domain.Article, int, error) {
	return s.GetAllArticlesFn(ctx, offset, limit)
}
func (s ArticleRepo) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	return s.GetArticleFn(ctx, id)
}
func (s ArticleRepo) CreateArticle(ctx context.Context, v *domain.Article) error {
	return s.CreateArticleFn(ctx, v)
}
func (s ArticleRepo) UpdateArticle(ctx context.Context, v *domain.Article) error {
	return s.UpdateArticleFn(ctx, v)
}
func (s ArticleRepo) DeleteArticle(ctx context.Context, id string) error {
	return s.DeleteArticleFn(ctx, id)
}
