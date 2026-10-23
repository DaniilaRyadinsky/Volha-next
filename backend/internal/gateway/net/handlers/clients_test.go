package handlers

import (
	"context"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type productStub struct {
	respond func(string, any) (any, error)
}

func (s productStub) CreateArticle(ctx context.Context, in *productsRPC.Article, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateArticle", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateArticle(ctx context.Context, in *productsRPC.Article, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateArticle", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteArticle(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteArticle", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetArticle(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Article, error) {
	v, err := s.respond("GetArticle", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Article), err
}
func (s productStub) GetAllArticles(ctx context.Context, in *productsRPC.Pagination, opts ...grpc.CallOption) (*productsRPC.ArticleList, error) {
	v, err := s.respond("GetAllArticles", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ArticleList), err
}
func (s productStub) GetDictionaries(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Dictionaries, error) {
	v, err := s.respond("GetDictionaries", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Dictionaries), err
}
func (s productStub) CreateProduct(ctx context.Context, in *productsRPC.ProductId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateProduct", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateProduct(ctx context.Context, in *productsRPC.ProductId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateProduct", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteProduct(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteProduct", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllProducts(ctx context.Context, in *productsRPC.Pagination, opts ...grpc.CallOption) (*productsRPC.ProductList, error) {
	v, err := s.respond("GetAllProducts", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ProductList), err
}
func (s productStub) GetProduct(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Product, error) {
	v, err := s.respond("GetProduct", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Product), err
}
func (s productStub) FilterProducts(ctx context.Context, in *productsRPC.ProductFilter, opts ...grpc.CallOption) (*productsRPC.ProductCatalogList, error) {
	v, err := s.respond("FilterProducts", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ProductCatalogList), err
}
func (s productStub) CreateBrand(ctx context.Context, in *productsRPC.Brand, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateBrand", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateBrand(ctx context.Context, in *productsRPC.Brand, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateBrand", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteBrand(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteBrand", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllBrands(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.BrandList, error) {
	v, err := s.respond("GetAllBrands", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.BrandList), err
}
func (s productStub) GetBrand(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Brand, error) {
	v, err := s.respond("GetBrand", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Brand), err
}
func (s productStub) CreateCategory(ctx context.Context, in *productsRPC.Category, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateCategory", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateCategory(ctx context.Context, in *productsRPC.Category, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateCategory", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteCategory(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteCategory", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllCategories(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.CategoryList, error) {
	v, err := s.respond("GetAllCategories", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.CategoryList), err
}
func (s productStub) GetCategory(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Category, error) {
	v, err := s.respond("GetCategory", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Category), err
}
func (s productStub) CreateCountry(ctx context.Context, in *productsRPC.Country, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateCountry", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateCountry(ctx context.Context, in *productsRPC.Country, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateCountry", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteCountry(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteCountry", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllCountries(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.CountryList, error) {
	v, err := s.respond("GetAllCountries", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.CountryList), err
}
func (s productStub) GetCountry(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Country, error) {
	v, err := s.respond("GetCountry", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Country), err
}
func (s productStub) CreateMaterial(ctx context.Context, in *productsRPC.Material, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateMaterial", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateMaterial(ctx context.Context, in *productsRPC.Material, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateMaterial", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteMaterial(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteMaterial", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllMaterials(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.MaterialList, error) {
	v, err := s.respond("GetAllMaterials", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.MaterialList), err
}
func (s productStub) GetMaterial(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Material, error) {
	v, err := s.respond("GetMaterial", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Material), err
}
func (s productStub) CreateColor(ctx context.Context, in *productsRPC.Color, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateColor", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateColor(ctx context.Context, in *productsRPC.Color, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateColor", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteColor(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteColor", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllColors(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.ColorList, error) {
	v, err := s.respond("GetAllColors", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ColorList), err
}
func (s productStub) GetColor(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Color, error) {
	v, err := s.respond("GetColor", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Color), err
}
func (s productStub) CreateSlide(ctx context.Context, in *productsRPC.Slide, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateSlide", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateSlide(ctx context.Context, in *productsRPC.Slide, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateSlide", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteSlide(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteSlide", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllSlides(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.SlideList, error) {
	v, err := s.respond("GetAllSlides", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.SlideList), err
}
func (s productStub) GetSlide(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.Slide, error) {
	v, err := s.respond("GetSlide", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Slide), err
}
func (s productStub) CreateProductColorPhotos(ctx context.Context, in *productsRPC.ProductColorPhotos, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("CreateProductColorPhotos", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) UpdateProductColorPhotos(ctx context.Context, in *productsRPC.ProductColorPhotos, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("UpdateProductColorPhotos", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) DeleteProductColorPhotos(ctx context.Context, in *productsRPC.ProductColorPhotosId, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("DeleteProductColorPhotos", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s productStub) GetAllProductColorPhotos(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*productsRPC.ProductColorPhotosList, error) {
	v, err := s.respond("GetAllProductColorPhotos", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ProductColorPhotosList), err
}
func (s productStub) GetProductColorPhotosByProduct(ctx context.Context, in *productsRPC.Id, opts ...grpc.CallOption) (*productsRPC.ProductColorPhotosList, error) {
	v, err := s.respond("GetProductColorPhotosByProduct", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.ProductColorPhotosList), err
}
func (s productStub) GetPhotosByProductAndColor(ctx context.Context, in *productsRPC.ProductColorPhotosId, opts ...grpc.CallOption) (*productsRPC.PhotoList, error) {
	v, err := s.respond("GetPhotosByProductAndColor", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.PhotoList), err
}
func (s productStub) IsFileReferenced(ctx context.Context, in *productsRPC.FileReferenceRequest, opts ...grpc.CallOption) (*productsRPC.FileReferenceResponse, error) {
	v, err := s.respond("IsFileReferenced", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.FileReferenceResponse), err
}

type authStub struct {
	respond func(string, any) (any, error)
}

func (s authStub) Register(ctx context.Context, in *productsRPC.RegisterRequest, opts ...grpc.CallOption) (*productsRPC.AuthResponse, error) {
	v, err := s.respond("Register", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthResponse), err
}
func (s authStub) Login(ctx context.Context, in *productsRPC.LoginRequest, opts ...grpc.CallOption) (*productsRPC.AuthResponse, error) {
	v, err := s.respond("Login", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthResponse), err
}
func (s authStub) Refresh(ctx context.Context, in *productsRPC.RefreshRequest, opts ...grpc.CallOption) (*productsRPC.AuthResponse, error) {
	v, err := s.respond("Refresh", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthResponse), err
}
func (s authStub) Logout(ctx context.Context, in *productsRPC.LogoutRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	v, err := s.respond("Logout", in)
	if v == nil {
		return nil, err
	}
	return v.(*emptypb.Empty), err
}
func (s authStub) Authenticate(ctx context.Context, in *productsRPC.AuthenticateRequest, opts ...grpc.CallOption) (*productsRPC.AuthenticatedUser, error) {
	v, err := s.respond("Authenticate", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthenticatedUser), err
}
func (s authStub) ChangePassword(ctx context.Context, in *productsRPC.ChangePasswordRequest, opts ...grpc.CallOption) (*productsRPC.AuthUser, error) {
	v, err := s.respond("ChangePassword", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthUser), err
}
func (s authStub) ChangeLogin(ctx context.Context, in *productsRPC.ChangeLoginRequest, opts ...grpc.CallOption) (*productsRPC.AuthUser, error) {
	v, err := s.respond("ChangeLogin", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.AuthUser), err
}
func (s authStub) GetCart(ctx context.Context, in *productsRPC.CartRequest, opts ...grpc.CallOption) (*productsRPC.Cart, error) {
	v, err := s.respond("GetCart", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Cart), err
}
func (s authStub) AddCartItem(ctx context.Context, in *productsRPC.CartItemRequest, opts ...grpc.CallOption) (*productsRPC.Cart, error) {
	v, err := s.respond("AddCartItem", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Cart), err
}
func (s authStub) UpdateCartItem(ctx context.Context, in *productsRPC.CartItemRequest, opts ...grpc.CallOption) (*productsRPC.Cart, error) {
	v, err := s.respond("UpdateCartItem", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Cart), err
}
func (s authStub) RemoveCartItem(ctx context.Context, in *productsRPC.CartItemKey, opts ...grpc.CallOption) (*productsRPC.Cart, error) {
	v, err := s.respond("RemoveCartItem", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Cart), err
}
func (s authStub) ClearCart(ctx context.Context, in *productsRPC.CartRequest, opts ...grpc.CallOption) (*productsRPC.Cart, error) {
	v, err := s.respond("ClearCart", in)
	if v == nil {
		return nil, err
	}
	return v.(*productsRPC.Cart), err
}
