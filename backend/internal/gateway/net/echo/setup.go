package echo

import (
	"context"
	"errors"
	"fmt"

	"github.com/autumnterror/volha-backend/internal/gateway/config"
	"github.com/autumnterror/volha-backend/internal/gateway/grpc/products"
	"github.com/autumnterror/volha-backend/internal/gateway/net/handlers"
	"github.com/autumnterror/volha-backend/internal/gateway/net/mw"
	"github.com/autumnterror/volha-backend/internal/gateway/redis"
	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/autumnterror/volha-backend/internal/platform/health"
	"github.com/autumnterror/volha-backend/internal/platform/observability"
	"github.com/autumnterror/volha-backend/internal/platform/requestid"

	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
)

type Echo struct {
	e          *echo.Echo
	cfg        *config.Config
	rds        *redis.Client
	apiProduct *products.Client
}

const (
	ImagesDir      = "./images"
	MaxUploadBytes = 10 << 20 // 10 MB
)

func New(rds *redis.Client, a *products.Client, cfg *config.Config, obsCfg observability.Config, metrics *observability.Metrics) *Echo {
	e := echo.New()

	h := handlers.New(a, rds, cfg)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health/live", func(c echo.Context) error {
		return c.JSON(http.StatusOK, health.Live())
	})
	e.GET("/health/ready", func(c echo.Context) error {
		if !a.Ready(c.Request().Context()) {
			status, code := health.Ready(errors.New("product-service is not ready"))
			return c.JSON(code, status)
		}
		status, code := health.Ready(nil)
		return c.JSON(code, status)
	})
	e.Use(middleware.Recover())
	e.Use(requestid.EchoMiddleware())
	e.Use(otelecho.Middleware(obsCfg.Service, otelecho.WithSkipper(skipObservabilityRoute)))
	e.Use(observability.EchoMiddleware(obsCfg, metrics))
	e.Use(middleware.BodyLimit(fmt.Sprintf("%d", MaxUploadBytes)))
	e.Use(middleware.CORS())
	e.Static("/images", ImagesDir)

	userApi := e.Group("/api", mw.CheckId())
	{
		p := userApi.Group("/product")
		{
			p.POST("/filter", h.FilterProducts)
			p.GET("/all", h.GetAllProducts)
			p.GET("", h.GetProduct)
		}

		b := userApi.Group("/brand")
		{
			b.GET("/all", h.GetAllBrands)
			b.GET("", h.GetBrand)
		}

		c := userApi.Group("/category")
		{
			c.GET("/all", h.GetAllCategories)
			c.GET("", h.GetCategory)
		}

		co := userApi.Group("/color")
		{
			co.GET("/all", h.GetAllColors)
			co.GET("", h.GetColor)
		}

		m := userApi.Group("/material")
		{
			m.GET("/all", h.GetAllMaterials)
			m.GET("", h.GetMaterial)
		}

		ct := userApi.Group("/country")
		{
			ct.GET("/all", h.GetAllCountries)
			ct.GET("", h.GetCountry)
		}

		dict := userApi.Group("/dictionaries")
		{
			dict.GET("/all", h.GetAllDictionaries)
			dict.GET("/all/by-category", h.GetAllDictionariesByCategory)
		}

		cp := userApi.Group("/colorphotos")
		{
			cp.GET("/all", h.GetAllProductColorPhotos)
			cp.POST("/photos/get", h.GetPhotosByProductAndColor)
		}

		sl := userApi.Group("/slide")
		{
			sl.GET("/all", h.GetAllSlides)
			sl.GET("", h.GetSlide)
		}
	}

	adminApi := e.Group("/api", mw.CheckId(), mw.DevAdminAuth(cfg))
	{
		f := adminApi.Group("/files")
		{
			f.POST("", h.UploadFile)
			f.DELETE("", h.DeleteFile)
		}
		p := adminApi.Group("/product")
		{
			p.POST("", h.CreateProduct)
			p.PUT("", h.UpdateProduct)
			p.DELETE("", h.DeleteProduct)
		}

		b := adminApi.Group("/brand")
		{
			b.POST("", h.CreateBrand)
			b.PUT("", h.UpdateBrand)
			b.DELETE("", h.DeleteBrand)
		}

		c := adminApi.Group("/category")
		{
			c.POST("", h.CreateCategory)
			c.PUT("", h.UpdateCategory)
			c.DELETE("", h.DeleteCategory)
		}

		co := adminApi.Group("/color")
		{
			co.POST("", h.CreateColor)
			co.PUT("", h.UpdateColor)
			co.DELETE("", h.DeleteColor)
		}

		m := adminApi.Group("/material")
		{
			m.POST("", h.CreateMaterial)
			m.PUT("", h.UpdateMaterial)
			m.DELETE("", h.DeleteMaterial)
		}

		ct := adminApi.Group("/country")
		{
			ct.POST("", h.CreateCountry)
			ct.PUT("", h.UpdateCountry)
			ct.DELETE("", h.DeleteCountry)
		}
		cp := adminApi.Group("/colorphotos")
		{
			cp.POST("", h.CreateProductColorPhotos)
			cp.PUT("", h.UpdateProductColorPhotos)
			cp.DELETE("", h.DeleteProductColorPhotos)
		}

		sl := adminApi.Group("/slide")
		{
			sl.POST("", h.CreateSlide)
			sl.PUT("", h.UpdateSlide)
			sl.DELETE("", h.DeleteSlide)
		}
	}

	return &Echo{
		e:          e,
		cfg:        cfg,
		rds:        rds,
		apiProduct: a,
	}
}

func skipObservabilityRoute(c echo.Context) bool {
	path := c.Path()
	return path == "/health/live" ||
		path == "/health/ready" ||
		path == "/metrics" ||
		strings.HasPrefix(path, "/debug/pprof")
}

func (e *Echo) MustRun() {
	const op = "echo.Run"

	if err := e.Run(); err != nil {
		e.e.Logger.Error(format.Error(op, err))
	}
}

func (e *Echo) Run() error {
	const op = "echo.Run"
	if err := e.e.Start(fmt.Sprintf(":%d", e.cfg.Port)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return format.Error(op, err)
	}
	return nil
}

func (e *Echo) Stop(ctx context.Context) error {
	const op = "echo.Stop"

	if err := e.e.Shutdown(ctx); err != nil {
		return format.Error(op, err)
	}
	return nil
}
