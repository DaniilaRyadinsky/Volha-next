package handlers

import (
	"context"
	"errors"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"

	"github.com/autumnterror/volha-backend/internal/gateway/model"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/autumnterror/volha-backend/internal/gateway/redis"
	log "github.com/autumnterror/volha-backend/internal/platform/logger"

	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// GetPhotosByProductAndColor godoc
// @Summary Получить фотографии продукта по ID продукта и цвета
// @Description Возвращает массив фотографий для указанного продукта и цвета
// @Tags product_color_photos
// @Accept json
// @Produce json
// @Param id body productsRPC.ProductColorPhotosId true "ID продукта и цвета"
// @Success 200 {object} []string "Список фотографий"
// @Failure 400 {object} view.ErrorResponse "Неверные данные"
// @Failure 404 {object} view.ErrorResponse "Не найдено"
// @Failure 503 {object} view.ErrorResponse "Сервис недоступен"
// @Failure 502 {object} view.ErrorResponse "Ошибка взаимодействия с сервисом"
// @Router /api/colorphotos/photos/get [post]
func (a *Apis) GetPhotosByProductAndColor(c echo.Context) error {
	const op = "handlers.GetPhotosByProductAndColor"

	var pcpId productsRPC.ProductColorPhotosId
	if err := c.Bind(&pcpId); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	photos, err := a.apiProduct.API.GetPhotosByProductAndColor(ctx, &pcpId)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	return c.JSON(http.StatusOK, photos.GetItems())
}

// GetAllDictionariesByCategory godoc
// @Summary Получить все справочники по категории
// @Description Возвращает список всех доступных справочных данных: бренды, материалы, страны и цвета и размеры по определенной категории
// @Tags dictionaries
// @Produce json
// @Param id query string true "ID категории"
// @Success 200 {object} model.DictionariesByCategory "Успешный запрос"
// @Failure 400 {object} view.ErrorResponse "Ошибка на сервере"
// @Failure 502 {object} view.ErrorResponse "Ошибка взаимодействия с сервисом"
// @Router /api/dictionaries/all/by-category [get]
func (a *Apis) GetAllDictionariesByCategory(c echo.Context) error {
	const op = "handlers.GetAllDictionariesByCategory"

	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "missing id"})
	}

	if ds, err := a.rds.GetDictionariesByCategory(c.Request().Context(), id); err == nil {
		return c.JSON(http.StatusOK, ds)
	} else if !errors.Is(err, redis.ErrCacheMiss) {
		log.Error(c.Request().Context(), "get dictionaries by category cache failed", "operation", op, "error", err)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	data, err := a.apiProduct.API.GetDictionaries(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	if len(data.GetMaterials()) == 0 {
		data.Materials = []*productsRPC.Material{}
	}
	if len(data.GetColors()) == 0 {
		data.Colors = []*productsRPC.Color{}
	}
	if len(data.GetBrands()) == 0 {
		data.Brands = []*productsRPC.Brand{}
	}
	if len(data.GetCountries()) == 0 {
		data.Countries = []*productsRPC.Country{}
	}

	response := model.DictionariesByCategory{
		Brands:           data.GetBrands(),
		Categories:       data.GetCategories(),
		Countries:        data.GetCountries(),
		Materials:        data.GetMaterials(),
		Colors:           data.GetColors(),
		MinPrice:         data.GetMinPrice(),
		MaxPrice:         data.GetMaxPrice(),
		MinWidth:         data.GetMinWidth(),
		MaxWidth:         data.GetMaxWidth(),
		MinHeight:        data.GetMinHeight(),
		MaxHeight:        data.GetMaxHeight(),
		MinDepth:         data.GetMinDepth(),
		MaxDepth:         data.GetMaxDepth(),
		FilteredCategory: filteredCategoryTitle(data.GetCategories(), id),
	}
	if err := a.rds.SetDictionariesByCategory(ctx, id, &response); err != nil {
		log.Error(ctx, "set dictionaries by category cache failed", "operation", op, "error", err)
	}

	return c.JSON(http.StatusOK, response)
}

func filteredCategoryTitle(categories []*productsRPC.Category, id string) string {
	for _, category := range categories {
		if category.GetId() == id {
			return category.GetTitle()
		}
	}
	return ""
}

// GetAllDictionaries godoc
// @Summary Получить все справочники
// @Description Возвращает список всех доступных справочных данных: бренды, категории, материалы, страны и цвета
// @Tags dictionaries
// @Produce json
// @Success 200 {object} productsRPC.Dictionaries "Успешный запрос"
// @Failure 400 {object} view.ErrorResponse "Ошибка на сервере"
// @Failure 502 {object} view.ErrorResponse "Ошибка взаимодействия с сервисом"
// @Router /api/dictionaries/all [get]
func (a *Apis) GetAllDictionaries(c echo.Context) error {
	const op = "handlers.GetAllDictionaries"

	if ds, err := a.rds.GetDictionaries(c.Request().Context()); err == nil {
		return c.JSON(http.StatusOK, ds)
	} else if !errors.Is(err, redis.ErrCacheMiss) {
		log.Error(c.Request().Context(), "get dictionaries cache failed", "operation", op, "error", err)
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	data, err := a.apiProduct.API.GetDictionaries(ctx, &productsRPC.Id{Id: notByCategory})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	if len(data.GetCategories()) == 0 {
		data.Categories = []*productsRPC.Category{}
	}
	if len(data.GetMaterials()) == 0 {
		data.Materials = []*productsRPC.Material{}
	}
	if len(data.GetColors()) == 0 {
		data.Colors = []*productsRPC.Color{}
	}
	if len(data.GetBrands()) == 0 {
		data.Brands = []*productsRPC.Brand{}
	}
	if len(data.GetCountries()) == 0 {
		data.Countries = []*productsRPC.Country{}
	}

	if err := a.rds.SetDictionaries(ctx, data); err != nil {
		log.Error(ctx, "set dictionaries cache failed", "operation", op, "error", err)
	}
	return c.JSON(http.StatusOK, data)
}

// DeleteProductColorPhotos godoc
// @Summary Удалить фотографии продукта для цвета
// @Description Удаляет записи фотографий продукта для указанного цвета
// @Tags product_color_photos
// @Accept json
// @Produce json
// @Param id body productsRPC.ProductColorPhotosId true "ID продукта и цвета"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/colorphotos [delete]
func (a *Apis) DeleteProductColorPhotos(c echo.Context) error {
	const op = "handlers.DeleteProductColorPhotos"

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	var pcpId productsRPC.ProductColorPhotosId

	if err := c.Bind(&pcpId); err != nil {
		log.Error(ctx, "bind product color photos id failed", "operation", op, "error", err)
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}

	pcp, err := a.apiProduct.API.GetPhotosByProductAndColor(ctx, &pcpId)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	if _, err := a.apiProduct.API.DeleteProductColorPhotos(ctx, &pcpId); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	for _, filename := range pcp.Items {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}

	return c.JSON(http.StatusOK, view.MessageResponse{Message: "product color photos"})
}

// DeleteProduct godoc
// @Summary Удалить продукт
// @Description Удаляет продукт по ID
// @Tags product
// @Produce json
// @Param id query string true "ID продукта"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/product [delete]
func (a *Apis) DeleteProduct(c echo.Context) error {
	const op = "handlers.DeleteProduct"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	pr, err := a.apiProduct.API.GetProduct(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	colorPhotos, err := a.productColorPhotoFiles(ctx, id)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	if _, err := a.apiProduct.API.DeleteProduct(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)

	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	for _, filename := range pr.GetPhotos() {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}
	for _, filename := range colorPhotos {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "product deleted successfully"})
}

func (a *Apis) productColorPhotoFiles(ctx context.Context, productID string) ([]string, error) {
	pcps, err := a.apiProduct.API.GetProductColorPhotosByProduct(ctx, &productsRPC.Id{Id: productID})
	if err != nil {
		return nil, err
	}
	files := make([]string, 0)
	for _, pcp := range pcps.GetItems() {
		files = append(files, pcp.GetPhotos()...)
	}
	return files, nil
}
