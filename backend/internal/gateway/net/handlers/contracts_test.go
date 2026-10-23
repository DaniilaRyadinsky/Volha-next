package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/alicebob/miniredis/v2"
	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/config"
	authgrpc "github.com/autumnterror/volha-backend/internal/gateway/grpc/auth"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	"github.com/autumnterror/volha-backend/internal/gateway/net/authcookie"
	cacheredis "github.com/autumnterror/volha-backend/internal/gateway/redis"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func gatewayFixture(t *testing.T) *Apis {
	t.Helper()
	t.Chdir(t.TempDir())
	require.NoError(t, os.MkdirAll("images", 0755))
	require.NoError(t, os.MkdirAll("pdfs", 0755))
	for _, name := range []string{"one.jpg", "two.jpg", "cat.jpg", "cover.jpg"} {
		require.NoError(t, os.WriteFile(filepath.Join("images", name), []byte("image"), 0600))
	}
	require.NoError(t, os.WriteFile("pdfs/guide.pdf", []byte("%PDF-1.7 test"), 0600))
	server := miniredis.RunT(t)
	cfg := &config.Config{RedisAddr: server.Addr(), RedisDictTTL: 1000000000}
	cache, err := cacheredis.New(cfg, nil)
	require.NoError(t, err)
	t.Cleanup(func() { cache.Close() })
	return New(&products.Client{API: productStub{respond: fixtureResponse}}, &authgrpc.Client{API: authStub{respond: authResponse}}, cache, cfg)
}
func request(t *testing.T, handler echo.HandlerFunc, method, url, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	require.NoError(t, handler(echo.New().NewContext(req, rec)))
	return rec
}
func fixtureResponse(method string, in any) (any, error) {
	switch method {
	case "GetAllProducts":
		return &productsRPC.ProductList{Items: []*productsRPC.Product{{Id: "p", Slug: "desk", Title: "Desk", Category: &productsRPC.Category{Id: "c", Slug: "cat", Title: "Category"}, Photos: []string{"one.jpg", "two.jpg"}, Price: 10}}, Total: 1}, nil
	case "GetProduct":
		return &productsRPC.Product{Id: "p", Article: "12345678", Slug: "desk", Title: "Desk", Photos: []string{"one.jpg", "two.jpg"}, Price: 10}, nil
	case "FilterProducts":
		return &productsRPC.ProductCatalogList{Items: []*productsRPC.ProductCatalogItem{{Id: "p", Title: "Desk", CategoryDetails: &productsRPC.Category{Id: "c", Title: "Category"}, Photos: []string{"one.jpg", "two.jpg"}}}, Total: 1}, nil
	case "GetDictionaries":
		return &productsRPC.Dictionaries{Countries: []*productsRPC.Country{{Id: "co", Title: "Country", Friendly: "hidden"}}, MaxPrice: 10}, nil
	case "GetAllProductColorPhotos", "GetProductColorPhotosByProduct":
		return &productsRPC.ProductColorPhotosList{}, nil
	case "GetPhotosByProductAndColor":
		return &productsRPC.PhotoList{Items: []string{"one.jpg"}}, nil
	case "IsFileReferenced":
		return &productsRPC.FileReferenceResponse{Referenced: false}, nil
	case "GetArticle":
		return &productsRPC.Article{Id: "a", Title: "Guide", Slug: "guide", Img: []string{"cover.jpg"}, Text: "/pdfs/guide.pdf", CreationTime: 1720000000123}, nil
	case "GetAllArticles":
		a, _ := fixtureResponse("GetArticle", nil)
		return &productsRPC.ArticleList{Items: []*productsRPC.Article{a.(*productsRPC.Article)}, Total: 1}, nil

	case "GetBrand":
		return &productsRPC.Brand{Id: "b", Title: "Brand"}, nil
	case "GetAllBrands":
		return &productsRPC.BrandList{Items: []*productsRPC.Brand{{Id: "b", Title: "Brand"}}}, nil
	case "GetCategory":
		return &productsRPC.Category{Id: "c", Slug: "cat", Title: "Category", Img: "cat.jpg"}, nil
	case "GetAllCategories":
		return &productsRPC.CategoryList{Items: []*productsRPC.Category{{Id: "c", Slug: "cat", Title: "Category", Img: "cat.jpg"}}}, nil
	case "GetCountry":
		return &productsRPC.Country{Id: "co", Title: "Country", Friendly: "hidden"}, nil
	case "GetAllCountries":
		return &productsRPC.CountryList{Items: []*productsRPC.Country{{Id: "co", Title: "Country", Friendly: "hidden"}}}, nil
	case "GetMaterial":
		return &productsRPC.Material{Id: "m", Title: "Wood"}, nil
	case "GetAllMaterials":
		return &productsRPC.MaterialList{Items: []*productsRPC.Material{{Id: "m", Title: "Wood"}}}, nil
	case "GetColor":
		return &productsRPC.Color{Id: "cl", Title: "Black", Hex: "#000000"}, nil
	case "GetAllColors":
		return &productsRPC.ColorList{Items: []*productsRPC.Color{{Id: "cl", Title: "Black", Hex: "#000000"}}}, nil
	case "GetSlide":
		return &productsRPC.Slide{Id: "s", Title: "Slide", Badge: "Subtitle", Link: "/", Img: "one.jpg", Img762: "two.jpg"}, nil
	case "GetAllSlides":
		return &productsRPC.SlideList{Items: []*productsRPC.Slide{{Id: "s", Title: "Slide", Badge: "Subtitle", Link: "/", Img: "one.jpg", Img762: "two.jpg"}}}, nil
	}
	return &emptypb.Empty{}, nil
}
func authResponse(method string, in any) (any, error) {
	switch method {
	case "Register", "Login", "Refresh":
		return &productsRPC.AuthResponse{User: &productsRPC.AuthUser{Id: "u", Email: "u@e.com", Login: "user"}, AccessToken: "access-secret", RefreshToken: "refresh-secret", AccessExpiresIn: 900, RefreshExpiresIn: 3600}, nil
	case "ChangePassword", "ChangeLogin":
		return &productsRPC.AuthUser{Id: "u", Login: "new"}, nil
	case "GetCart", "AddCartItem", "UpdateCartItem", "RemoveCartItem", "ClearCart":
		return &productsRPC.Cart{Items: []*productsRPC.CartItem{}, TotalPrice: 20}, nil
	}
	return &emptypb.Empty{}, nil
}
func namedHandler(a *Apis, name string) echo.HandlerFunc {
	return echo.HandlerFunc(reflect.ValueOf(a).MethodByName(name).Interface().(func(echo.Context) error))
}
func TestCatalogHTTP(t *testing.T) {
	for _, kind := range []struct{ single, plural, route string }{{"Brand", "Brands", "brand"}, {"Category", "Categories", "category"}, {"Country", "Countries", "country"}, {"Material", "Materials", "material"}, {"Color", "Colors", "color"}, {"Slide", "Slides", "slide"}} {
		t.Run(kind.single, func(t *testing.T) {
			a := gatewayFixture(t)
			id := xid.New().String()
			url := "/api/" + kind.route
			body := `{"id":"id","title":"Title","slug":"slug","link":"/catalog","img":"one.jpg","img762":"two.jpg","hex":"#000000","subtitle":"new subtitle"}`
			rec := request(t, namedHandler(a, "Get"+kind.single), "GET", url+"?id="+id, "")
			require.Equal(t, 200, rec.Code)
			var value map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &value))
			require.Contains(t, value, "id")
			require.Contains(t, value, "title")
			if kind.single == "Country" {
				require.Len(t, value, 2)
			}
			if kind.single == "Slide" {
				require.Equal(t, "Subtitle", value["subtitle"])
				require.NotContains(t, value, "badge")
			}
			rec = request(t, namedHandler(a, "GetAll"+kind.plural), "GET", url+"/all", "")
			require.Equal(t, 200, rec.Code)
			var items []map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &items))
			require.Len(t, items, 1)
			for _, op := range []string{"Create", "Update"} {
				require.Equal(t, 200, request(t, namedHandler(a, op+kind.single), "POST", url, body).Code)
			}
			require.Equal(t, 200, request(t, namedHandler(a, "Delete"+kind.single), "DELETE", url+"?id="+id, "").Code)
			require.Equal(t, 400, request(t, namedHandler(a, "Get"+kind.single), "GET", url, "").Code)
			a.apiProduct.API = productStub{respond: func(string, any) (any, error) { return nil, status.Error(codes.Unavailable, "down") }}
			require.Equal(t, 503, request(t, namedHandler(a, "GetAll"+kind.plural), "GET", url+"/all", "").Code)
		})
	}
}
func TestProductHTTP(t *testing.T) {
	a := gatewayFixture(t)
	rec := request(t, a.GetProduct, "GET", "/api/product?id=p", "")
	require.Equal(t, 200, rec.Code)
	var p map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &p))
	require.Equal(t, "12345678", p["sku"])
	require.Equal(t, []any{}, p["similarProducts"])
	require.NotContains(t, p, "article")
	require.NotContains(t, p, "photos")
	require.Nil(t, p["brand"])
	a.apiProduct.API = productStub{respond: func(method string, in any) (any, error) {
		if method == "FilterProducts" {
			f := in.(*productsRPC.ProductFilter)
			require.Equal(t, int32(10), f.GetOffset())
			require.Equal(t, "desk", f.GetTitle())
			require.Equal(t, int32(5), f.GetMinPrice())
		}
		if method == "CreateProduct" {
			require.Equal(t, "12345678", in.(*productsRPC.ProductId).GetArticle())
		}
		return fixtureResponse(method, in)
	}}
	rec = request(t, a.FilterProducts, "POST", "/api/product/filter", `{"page":2,"limit":10,"search":"desk","minPrice":5,"sortBy":"price"}`)
	require.Equal(t, 200, rec.Code)
	var list map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &list))
	require.Equal(t, float64(2), list["page"])
	item := list["items"].([]any)[0].(map[string]any)
	require.Len(t, item["photos"], 2)
	require.IsType(t, map[string]any{}, item["category"])
	require.NotContains(t, item, "image")
	require.Equal(t, 200, request(t, a.GetAllProducts, "GET", "/api/product/all?page=2&limit=10", "").Code)
	for _, body := range []string{`{"page":-1}`, `{"limit":101}`, `{`} {
		require.Equal(t, 400, request(t, a.FilterProducts, "POST", "/api/product/filter", body).Code)
	}
	body := `{"id":"p","sku":"12345678","title":"Desk","photos":["one.jpg"]}`
	require.Equal(t, 200, request(t, a.CreateProduct, "POST", "/api/product", body).Code)
	require.Equal(t, 200, request(t, a.UpdateProduct, "PUT", "/api/product", body).Code)
	require.Equal(t, 200, request(t, a.DeleteProduct, "DELETE", "/api/product?id=p", "").Code)
}
func TestDictionaryHTTPFromCacheAndService(t *testing.T) {
	a := gatewayFixture(t)
	for _, handler := range []echo.HandlerFunc{a.GetAllDictionaries, a.GetAllDictionariesByCategory} {
		for i := 0; i < 2; i++ {
			rec := request(t, handler, "GET", "/api/dictionaries?id=cat", "")
			require.Equal(t, 200, rec.Code)
			var d map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &d))
			require.Contains(t, d, "minPrice")
			require.Contains(t, d, "maxDepth")
			require.NotContains(t, d, "min_price")
			require.Equal(t, []any{}, d["brands"])
			require.Len(t, d["countries"].([]any)[0], 2)
		}
	}
}
func TestArticleHTTP(t *testing.T) {
	a := gatewayFixture(t)
	rec := request(t, a.GetArticle, "GET", "/api/article?id=a", "")
	require.Equal(t, 200, rec.Code)
	var v map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &v))
	require.Equal(t, float64(1720000000123), v["creationTime"])
	require.Equal(t, "/pdfs/guide.pdf", v["text"])
	require.Equal(t, 200, request(t, a.GetAllArticles, "GET", "/api/article/all?page=1&limit=10", "").Code)
	body := `{"id":"a","slug":"guide","title":"Guide","img":["cover.jpg"],"text":"/pdfs/guide.pdf"}`
	require.Equal(t, 200, request(t, a.CreateArticle, "POST", "/api/article", body).Code)
	require.Equal(t, 200, request(t, a.UpdateArticle, "PUT", "/api/article", body).Code)
	require.Equal(t, 400, request(t, a.CreateArticle, "POST", "/api/article", `{"text":"/pdfs/missing.pdf"}`).Code)
	require.Equal(t, 200, request(t, a.DeleteArticle, "DELETE", "/api/article?id=a", "").Code)
	_, err := os.Stat("pdfs/guide.pdf")
	require.True(t, os.IsNotExist(err))
}
func TestPDFUploadAndDelete(t *testing.T) {
	a := gatewayFixture(t)
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "guide.pdf")
	require.NoError(t, err)
	_, err = part.Write([]byte("%PDF-1.7\nexample PDF"))
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	req := httptest.NewRequest("POST", "/api/files", &body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	require.NoError(t, a.UploadFile(echo.New().NewContext(req, rec)))
	require.Equal(t, 200, rec.Code)
	var result struct {
		Name string `json:"name"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	require.True(t, strings.HasSuffix(result.Name, ".pdf"))
	data, err := os.ReadFile(filepath.Join("pdfs", result.Name))
	require.NoError(t, err)
	require.Contains(t, string(data), "%PDF-1.7")
	require.Equal(t, 200, request(t, a.DeleteFile, "DELETE", "/api/files?title="+result.Name, "").Code)
	for _, link := range []string{"../escape.pdf", "/pdfs/../escape.pdf", "/pdfs/dir/escape.pdf", "/pdfs/no.txt"} {
		_, err := normalizePDFLink(link)
		require.Error(t, err)
	}
}
func TestAuthCookiesAndProtectedHandlers(t *testing.T) {
	a := gatewayFixture(t)
	a.cfg.CookieSecure = true
	for _, handler := range []echo.HandlerFunc{a.Register, a.Login} {
		rec := request(t, handler, "POST", "/api/auth", `{"login":"user","password":"Password123"}`)
		require.Contains(t, []int{200, 201}, rec.Code)
		require.NotContains(t, rec.Body.String(), "secret")
		require.Len(t, rec.Result().Cookies(), 2)
		for _, cookie := range rec.Result().Cookies() {
			require.True(t, cookie.HttpOnly)
			require.True(t, cookie.Secure)
			require.Equal(t, http.SameSiteLaxMode, cookie.SameSite)
		}
	}
	req := httptest.NewRequest("POST", "/api/auth/refresh", nil)
	req.AddCookie(&http.Cookie{Name: authcookie.RefreshTokenName, Value: "refresh-secret"})
	rec := httptest.NewRecorder()
	require.NoError(t, a.Refresh(echo.New().NewContext(req, rec)))
	require.Equal(t, 200, rec.Code)
	require.Equal(t, 401, request(t, a.Refresh, "POST", "/api/auth/refresh", "").Code)
	require.Equal(t, 200, request(t, a.Logout, "POST", "/api/auth/logout", "").Code)
	for _, handler := range []echo.HandlerFunc{a.Me, a.GetCart, a.ClearCart, a.AddCartItem, a.UpdateCartItem, a.RemoveCartItem, a.ChangeLogin, a.ChangePassword} {
		req := httptest.NewRequest("POST", "/api/cart", strings.NewReader(`{"product_id":"p","quantity":1,"password":"Password123","new_login":"new","old_password":"Password123","new_password":"NewPassword123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)
		c.Set(authcookie.PrincipalContextKey, &productsRPC.AuthenticatedUser{User: &productsRPC.AuthUser{Id: "u"}, SessionId: "sid"})
		c.SetParamNames("product_id")
		c.SetParamValues("p")
		require.NoError(t, handler(c))
		require.Equal(t, 200, rec.Code)
	}
}
func TestPhotoHandlersAndPagination(t *testing.T) {
	a := gatewayFixture(t)
	body := `{"product_id":"p","color_id":"c","photos":["one.jpg"]}`
	for _, handler := range []echo.HandlerFunc{a.GetAllProductColorPhotos, a.GetPhotosByProductAndColor, a.CreateProductColorPhotos, a.UpdateProductColorPhotos, a.DeleteProductColorPhotos} {
		require.Equal(t, 200, request(t, handler, "POST", "/api/colorphotos", body).Code)
	}
	for _, url := range []string{"/?page=bad", "/?page=-1", "/?limit=101", "/?limit=bad", "/?offset=-1"} {
		c := echo.New().NewContext(httptest.NewRequest("GET", url, nil), httptest.NewRecorder())
		_, _, code, _ := pag(c)
		require.Equal(t, 400, code)
	}
}

func TestCatalogFailureAndEmptyResponsePaths(t *testing.T) {
	a := gatewayFixture(t)
	for _, method := range []string{"GetAllBrands", "GetAllCategories", "GetAllCountries", "GetAllMaterials", "GetAllColors", "GetAllSlides"} {
		base := fixtureResponse
		a.apiProduct.API = productStub{respond: func(name string, in any) (any, error) {
			value, err := base(name, in)
			if name == method {
				reflect.ValueOf(value).Elem().FieldByName("Items").Set(reflect.Zero(reflect.ValueOf(value).Elem().FieldByName("Items").Type()))
			}
			return value, err
		}}
		rec := request(t, namedHandler(a, method), "GET", "/", "")
		require.JSONEq(t, `[]`, rec.Body.String())
	}
	a.apiProduct.API = productStub{respond: func(string, any) (any, error) { return nil, status.Error(codes.Unavailable, "backend down") }}
	for _, name := range []string{"GetProduct", "GetAllProducts", "GetArticle", "GetAllArticles", "GetAllProductColorPhotos", "GetPhotosByProductAndColor", "GetAllBrands", "GetBrand", "GetAllCategories", "GetCategory", "GetAllCountries", "GetCountry", "GetAllMaterials", "GetMaterial", "GetAllColors", "GetColor", "GetAllSlides", "GetSlide", "DeleteProduct", "DeleteArticle", "DeleteBrand", "DeleteCategory", "DeleteCountry", "DeleteMaterial", "DeleteColor", "DeleteSlide", "DeleteProductColorPhotos"} {
		rec := request(t, namedHandler(a, name), "POST", "/?id=p", `{"product_id":"p","color_id":"c"}`)
		require.Equal(t, 503, rec.Code, name)
	}
	for _, name := range []string{"CreateBrand", "CreateCategory", "CreateCountry", "CreateMaterial", "CreateColor", "CreateSlide", "CreateProduct", "CreateProductColorPhotos", "UpdateBrand", "UpdateCategory", "UpdateCountry", "UpdateMaterial", "UpdateColor", "UpdateSlide", "UpdateProduct", "UpdateProductColorPhotos", "CreateArticle", "UpdateArticle", "Register", "Login", "ChangePassword", "ChangeLogin", "AddCartItem", "UpdateCartItem"} {
		rec := request(t, namedHandler(a, name), "POST", "/", `{`)
		require.Equal(t, 400, rec.Code, name)
	}
	for _, name := range []string{"Me", "GetCart", "ClearCart", "RemoveCartItem", "AddCartItem", "UpdateCartItem", "ChangePassword", "ChangeLogin"} {
		rec := request(t, namedHandler(a, name), "POST", "/", `{}`)
		require.Equal(t, 401, rec.Code, name)
	}
}
func TestReferencedPDFCannotBeDeleted(t *testing.T) {
	a := gatewayFixture(t)
	a.apiProduct.API = productStub{respond: func(string, any) (any, error) { return &productsRPC.FileReferenceResponse{Referenced: true}, nil }}
	rec := request(t, a.DeleteFile, "DELETE", "/api/files?title=guide.pdf", "")
	require.Equal(t, 409, rec.Code)
	_, err := os.Stat("pdfs/guide.pdf")
	require.NoError(t, err)
}
