//go:build integration

package integration

import (
	"context"
	"testing"

	pb "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/test/integration/testkit"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestIntegrationProductCatalogContentAndFiles(t *testing.T) {
	s := testkit.New(t)
	ctx := context.Background()
	c := s.SeedCatalog()
	id := testkit.ID()

	// Every reference catalog supports read, update and list over the real gRPC
	// transport and PostgreSQL. Deletes are checked after dependent products.
	brands, err := s.Products.GetAllBrands(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfBrands(brands.GetItems()), c.Brand)
	brand, err := s.Products.GetBrand(ctx, &pb.Id{Id: c.Brand})
	require.NoError(t, err)
	brand.Title = "Updated Brand " + id
	_, err = s.Products.UpdateBrand(ctx, brand)
	require.NoError(t, err)

	categories, err := s.Products.GetAllCategories(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfCategories(categories.GetItems()), c.Category)
	category, err := s.Products.GetCategory(ctx, &pb.Id{Id: c.Category})
	require.NoError(t, err)
	category.Description = "Updated category"
	_, err = s.Products.UpdateCategory(ctx, category)
	require.NoError(t, err)

	countries, err := s.Products.GetAllCountries(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfCountries(countries.GetItems()), c.Country)
	country, err := s.Products.GetCountry(ctx, &pb.Id{Id: c.Country})
	require.NoError(t, err)
	country.Friendly = "Updated country"
	_, err = s.Products.UpdateCountry(ctx, country)
	require.NoError(t, err)

	materials, err := s.Products.GetAllMaterials(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfMaterials(materials.GetItems()), c.Material)
	material, err := s.Products.GetMaterial(ctx, &pb.Id{Id: c.Material})
	require.NoError(t, err)
	material.Title = "Updated Material " + id
	_, err = s.Products.UpdateMaterial(ctx, material)
	require.NoError(t, err)

	colors, err := s.Products.GetAllColors(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfColors(colors.GetItems()), c.Color)
	color, err := s.Products.GetColor(ctx, &pb.Id{Id: c.Color})
	require.NoError(t, err)
	color.Hex = "#445566"
	_, err = s.Products.UpdateColor(ctx, color)
	require.NoError(t, err)

	product := s.SeedProduct(c, "shared.jpg")
	got, err := s.Products.GetProduct(ctx, &pb.Id{Id: product.Id})
	require.NoError(t, err)
	require.Equal(t, c.Brand, got.GetBrand().GetId())
	require.Equal(t, c.Category, got.GetCategory().GetId())
	require.Equal(t, c.Country, got.GetCountry().GetId())
	require.Len(t, got.GetMaterials(), 1)
	require.Len(t, got.GetColors(), 1)

	product.Title = "Updated Product " + id
	product.Price = 120
	_, err = s.Products.UpdateProduct(ctx, product)
	require.NoError(t, err)
	got, err = s.Products.GetProduct(ctx, &pb.Id{Id: product.Id})
	require.NoError(t, err)
	require.Equal(t, product.Title, got.GetTitle())

	list, err := s.Products.GetAllProducts(ctx, &pb.Pagination{Limit: 20})
	require.NoError(t, err)
	require.Contains(t, idsOfProducts(list.GetItems()), product.Id)
	filtered, err := s.Products.FilterProducts(ctx, &pb.ProductFilter{Category: []string{c.Category}, Materials: []string{c.Material}, Colors: []string{c.Color}, Title: "Updated Product", Limit: 20})
	require.NoError(t, err)
	require.Equal(t, int32(1), filtered.GetTotal())
	require.Equal(t, product.Id, filtered.GetItems()[0].GetId())

	full, err := s.Products.GetDictionaries(ctx, &pb.Id{Id: "CLR"})
	require.NoError(t, err)
	require.NotEmpty(t, full.GetCategories())
	byCategory, err := s.Products.GetDictionaries(ctx, &pb.Id{Id: c.Category})
	require.NoError(t, err)
	require.Equal(t, int32(120), byCategory.GetMinPrice())
	require.Equal(t, int32(120), byCategory.GetMaxPrice())
	similar := s.SeedProduct(c, "similar.jpg")
	product.Seems = []string{similar.Id}
	_, err = s.Products.UpdateProduct(ctx, product)
	require.NoError(t, err)
	got, err = s.Products.GetProduct(ctx, &pb.Id{Id: product.Id})
	require.NoError(t, err)
	require.Contains(t, idsOfProducts(got.GetSeems()), similar.Id)

	photos := &pb.ProductColorPhotos{ProductId: product.Id, ColorId: c.Color, Photos: []string{"shared.jpg", "detail.jpg"}}
	_, err = s.Products.CreateProductColorPhotos(ctx, photos)
	require.NoError(t, err)
	photoList, err := s.Products.GetAllProductColorPhotos(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.NotEmpty(t, photoList.GetItems())
	byProduct, err := s.Products.GetProductColorPhotosByProduct(ctx, &pb.Id{Id: product.Id})
	require.NoError(t, err)
	require.Len(t, byProduct.GetItems(), 1)
	photoKey := &pb.ProductColorPhotosId{ProductId: product.Id, ColorId: c.Color}
	photoNames, err := s.Products.GetPhotosByProductAndColor(ctx, photoKey)
	require.NoError(t, err)
	require.Equal(t, photos.Photos, photoNames.GetItems())
	photos.Photos = []string{"detail.jpg"}
	_, err = s.Products.UpdateProductColorPhotos(ctx, photos)
	require.NoError(t, err)
	photoNames, err = s.Products.GetPhotosByProductAndColor(ctx, photoKey)
	require.NoError(t, err)
	require.Equal(t, photos.Photos, photoNames.GetItems())

	slide := &pb.Slide{Id: testkit.ID(), Link: "/", Img: "shared.jpg", Img762: "detail.jpg"}
	_, err = s.Products.CreateSlide(ctx, slide)
	require.NoError(t, err)
	s.Track("slides", slide.Id)
	slides, err := s.Products.GetAllSlides(ctx, &emptypb.Empty{})
	require.NoError(t, err)
	require.Contains(t, idsOfSlides(slides.GetItems()), slide.Id)
	slide.Title = "New title"
	_, err = s.Products.UpdateSlide(ctx, slide)
	require.NoError(t, err)
	loadedSlide, err := s.Products.GetSlide(ctx, &pb.Id{Id: slide.Id})
	require.NoError(t, err)
	require.Equal(t, "New title", loadedSlide.GetTitle())

	article := &pb.Article{Id: testkit.ID(), Slug: "guide-" + id, Title: "Guide", Img: []string{"shared.jpg"}, Text: "/pdfs/guide.pdf"}
	_, err = s.Products.CreateArticle(ctx, article)
	require.NoError(t, err)
	s.Track("articles", article.Id)
	loadedArticle, err := s.Products.GetArticle(ctx, &pb.Id{Id: article.Id})
	require.NoError(t, err)
	require.Positive(t, loadedArticle.GetCreationTime())
	articles, err := s.Products.GetAllArticles(ctx, &pb.Pagination{Limit: 20})
	require.NoError(t, err)
	require.Contains(t, idsOfArticles(articles.GetItems()), article.Id)
	article.Title = "Updated Guide"
	_, err = s.Products.UpdateArticle(ctx, article)
	require.NoError(t, err)
	loadedArticle, err = s.Products.GetArticle(ctx, &pb.Id{Id: article.Id})
	require.NoError(t, err)
	require.Equal(t, "Updated Guide", loadedArticle.GetTitle())

	referenced, err := s.Products.IsFileReferenced(ctx, &pb.FileReferenceRequest{Filename: "shared.jpg"})
	require.NoError(t, err)
	require.True(t, referenced.GetReferenced())
	referenced, err = s.Products.IsFileReferenced(ctx, &pb.FileReferenceRequest{Filename: "guide.pdf"})
	require.NoError(t, err)
	require.True(t, referenced.GetReferenced())

	_, err = s.Products.DeleteProductColorPhotos(ctx, photoKey)
	require.NoError(t, err)
	_, err = s.Products.DeleteArticle(ctx, &pb.Id{Id: article.Id})
	require.NoError(t, err)
	_, err = s.Products.DeleteSlide(ctx, &pb.Id{Id: slide.Id})
	require.NoError(t, err)
	_, err = s.Products.DeleteProduct(ctx, &pb.Id{Id: product.Id})
	require.NoError(t, err)
	_, err = s.Products.DeleteProduct(ctx, &pb.Id{Id: similar.Id})
	require.NoError(t, err)
	for _, del := range []struct {
		id string
		fn func(context.Context, *pb.Id, ...grpc.CallOption) (*emptypb.Empty, error)
	}{
		{c.Brand, s.Products.DeleteBrand}, {c.Category, s.Products.DeleteCategory},
		{c.Country, s.Products.DeleteCountry}, {c.Material, s.Products.DeleteMaterial}, {c.Color, s.Products.DeleteColor},
	} {
		_, err = del.fn(ctx, &pb.Id{Id: del.id})
		require.NoError(t, err)
	}
	referenced, err = s.Products.IsFileReferenced(ctx, &pb.FileReferenceRequest{Filename: "shared.jpg"})
	require.NoError(t, err)
	require.False(t, referenced.GetReferenced())
}

func TestIntegrationAuthSessionsAndCart(t *testing.T) {
	s := testkit.New(t)
	ctx := context.Background()
	c := s.SeedCatalog()
	p := s.SeedProduct(c, "cart.jpg")
	identity := testkit.ID()
	password := "InitialPassword123"
	pair, err := s.Auth.Register(ctx, &pb.RegisterRequest{Email: identity + "@example.com", Login: "u" + identity, Password: password})
	require.NoError(t, err)
	userID := pair.GetUser().GetId()
	s.Track("users", userID)
	require.NotEmpty(t, pair.GetAccessToken())

	_, err = s.Auth.Register(ctx, &pb.RegisterRequest{Email: identity + "@example.com", Login: "other" + identity, Password: password})
	require.Equal(t, codes.AlreadyExists, status.Code(err))
	_, err = s.Auth.Login(ctx, &pb.LoginRequest{Identifier: identity + "@example.com", Password: "wrong"})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	otherSession, err := s.Auth.Login(ctx, &pb.LoginRequest{Identifier: "u" + identity, Password: password})
	require.NoError(t, err)
	principal, err := s.Auth.Authenticate(ctx, &pb.AuthenticateRequest{AccessToken: pair.AccessToken})
	require.NoError(t, err)
	require.Equal(t, userID, principal.GetUser().GetId())

	cart, err := s.Auth.AddCartItem(ctx, &pb.CartItemRequest{UserId: userID, ProductId: p.Id, Quantity: 2})
	require.NoError(t, err)
	require.Equal(t, int32(2), cart.GetTotalQuantity())
	require.Equal(t, int64(160), cart.GetTotalPrice())
	cart, err = s.Auth.AddCartItem(ctx, &pb.CartItemRequest{UserId: userID, ProductId: p.Id, Quantity: 1})
	require.NoError(t, err)
	require.Equal(t, int32(3), cart.GetTotalQuantity())
	cart, err = s.Auth.GetCart(ctx, &pb.CartRequest{UserId: userID})
	require.NoError(t, err)
	require.Equal(t, p.Slug, cart.GetItems()[0].GetSlug())
	_, err = s.Auth.AddCartItem(ctx, &pb.CartItemRequest{UserId: userID, ProductId: p.Id, Quantity: 99})
	require.Equal(t, codes.ResourceExhausted, status.Code(err))
	cart, err = s.Auth.UpdateCartItem(ctx, &pb.CartItemRequest{UserId: userID, ProductId: p.Id, Quantity: 5})
	require.NoError(t, err)
	require.Equal(t, int32(5), cart.GetTotalQuantity())
	cart, err = s.Auth.RemoveCartItem(ctx, &pb.CartItemKey{UserId: userID, ProductId: p.Id})
	require.NoError(t, err)
	require.Empty(t, cart.GetItems())
	_, err = s.Auth.AddCartItem(ctx, &pb.CartItemRequest{UserId: userID, ProductId: p.Id, Quantity: 1})
	require.NoError(t, err)
	cart, err = s.Auth.ClearCart(ctx, &pb.CartRequest{UserId: userID})
	require.NoError(t, err)
	require.Zero(t, cart.GetTotalPrice())

	fresh, err := s.Auth.Refresh(ctx, &pb.RefreshRequest{RefreshToken: pair.RefreshToken})
	require.NoError(t, err)
	_, err = s.Auth.Refresh(ctx, &pb.RefreshRequest{RefreshToken: pair.RefreshToken})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	_, err = s.Auth.ChangePassword(ctx, &pb.ChangePasswordRequest{UserId: userID, SessionId: principal.SessionId, OldPassword: password, NewPassword: "ChangedPassword123"})
	require.NoError(t, err)
	_, err = s.Auth.Authenticate(ctx, &pb.AuthenticateRequest{AccessToken: otherSession.AccessToken})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
	_, err = s.Auth.ChangeLogin(ctx, &pb.ChangeLoginRequest{UserId: userID, SessionId: principal.SessionId, Password: "ChangedPassword123", NewLogin: "changed" + identity})
	require.NoError(t, err)
	_, err = s.Auth.Login(ctx, &pb.LoginRequest{Identifier: "changed" + identity, Password: "ChangedPassword123"})
	require.NoError(t, err)
	_, err = s.Auth.Logout(ctx, &pb.LogoutRequest{RefreshToken: fresh.RefreshToken})
	require.NoError(t, err)
	_, err = s.Auth.Authenticate(ctx, &pb.AuthenticateRequest{AccessToken: fresh.AccessToken})
	require.Equal(t, codes.Unauthenticated, status.Code(err))
}

func TestIntegrationBadPathsAndTransactionRollback(t *testing.T) {
	s := testkit.New(t)
	ctx := context.Background()
	c := s.SeedCatalog()
	valid := s.SeedProduct(c, "ok.jpg")
	missing := testkit.ID()
	_, err := s.Products.GetProduct(ctx, &pb.Id{Id: missing})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = s.Products.GetDictionaries(ctx, &pb.Id{Id: missing})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = s.Products.FilterProducts(ctx, &pb.ProductFilter{Limit: 1, SortBy: "drop table products"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = s.Products.CreateProduct(ctx, &pb.ProductId{Id: testkit.ID(), Title: "Invalid", Article: "bad", Brand: c.Brand, Category: c.Category, Country: c.Country, Price: 1})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	invalidFK := proto.Clone(valid).(*pb.ProductId)
	invalidFK.Id = testkit.ID()
	invalidFK.Title = "Broken " + invalidFK.Id
	invalidFK.Slug = "broken-" + invalidFK.Id
	invalidFK.Article = testkit.SKU(t)
	invalidFK.Brand = missing
	_, err = s.Products.CreateProduct(ctx, invalidFK)
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = s.Products.GetProduct(ctx, &pb.Id{Id: invalidFK.Id})
	require.Equal(t, codes.NotFound, status.Code(err), "failed insert must roll back the whole product")
	_, err = s.Products.DeleteBrand(ctx, &pb.Id{Id: c.Brand})
	require.Equal(t, codes.FailedPrecondition, status.Code(err), "linked brand must not be deleted")
	_, err = s.Products.CreateProductColorPhotos(ctx, &pb.ProductColorPhotos{ProductId: valid.Id, ColorId: missing, Photos: []string{"x.jpg"}})
	require.Equal(t, codes.FailedPrecondition, status.Code(err))
	_, err = s.Products.GetPhotosByProductAndColor(ctx, &pb.ProductColorPhotosId{ProductId: valid.Id, ColorId: c.Color})
	require.Equal(t, codes.NotFound, status.Code(err))
	_, err = s.Products.CreateArticle(ctx, &pb.Article{Id: testkit.ID(), Title: "Bad", Slug: "bad", Text: "https://example.com/file.pdf"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
	_, err = s.Products.IsFileReferenced(ctx, &pb.FileReferenceRequest{Filename: "../secret"})
	require.Equal(t, codes.InvalidArgument, status.Code(err))
}

func idsOfBrands(v []*pb.Brand) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfCategories(v []*pb.Category) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfCountries(v []*pb.Country) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfMaterials(v []*pb.Material) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfColors(v []*pb.Color) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfSlides(v []*pb.Slide) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfProducts(v []*pb.Product) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
func idsOfArticles(v []*pb.Article) []string {
	out := make([]string, 0, len(v))
	for _, x := range v {
		out = append(out, x.Id)
	}
	return out
}
