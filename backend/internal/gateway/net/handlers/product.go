package handlers

import (
	"net/http"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateProduct godoc
// @Summary Создать продукт
// @Description Добавляет новый продукт
// @Tags product
// @Accept json
// @Produce json
// @Param product body domain.ProductId true "Новый продукт"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/product [post]
func (a *Apis) CreateProduct(c echo.Context) error {
	const op = "handlers.CreateProduct"

	ctx, cancel := requestContext(c)
	defer cancel()

	var body domain.ProductId
	var p productsRPC.ProductId
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	p = *domain.ProductIdToRpc(&body)
	p.Id = xid.New().String()
	if err := ensureImageReferences(p.GetPhotos()...); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.CreateProduct(ctx, &p); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: p.Id})
}

// UpdateProduct godoc
// @Summary Обновить продукт
// @Description Обновляет данные существующего продукта
// @Tags product
// @Accept json
// @Produce json
// @Param product body domain.ProductId true "Обновлённые данные продукта"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/product [put]
func (a *Apis) UpdateProduct(c echo.Context) error {
	const op = "handlers.UpdateProduct"

	ctx, cancel := requestContext(c)
	defer cancel()

	var body domain.ProductId
	var p productsRPC.ProductId
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	p = *domain.ProductIdToRpc(&body)
	if err := ensureImageReferences(p.GetPhotos()...); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldProduct, err := a.apiProduct.API.GetProduct(ctx, &productsRPC.Id{Id: p.GetId()})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldColorPhotos, err := a.productColorPhotoFiles(ctx, p.GetId())
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.UpdateProduct(ctx, &p); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	for _, filename := range removedStrings(oldProduct.GetPhotos(), p.GetPhotos()) {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}
	for _, filename := range oldColorPhotos {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "product updated successfully"})
}

// GetProduct rand godoc
// @Summary Получить продукт
// @Description Возвращает продукт по id
// @Tags product
// @Produce json
// @Param id query string true "ID продукт"
// @Success 200 {object} domain.Product "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/product [get]
func (a *Apis) GetProduct(c echo.Context) error {
	const op = "handlers.GetProduct"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	product, err := a.apiProduct.API.GetProduct(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, domain.ProductFromRpc(product))
}

// FilterProducts godoc
// @Summary Получить продукты по фильтру
// @Description Возвращает список продуктов, соответствующих заданным критериям фильтрации
// @Tags product
// @Accept json
// @Produce json
// @Param filter body domain.ProductFilter true "Параметры фильтрации"
// @Success 200 {object} view.ProductCatalogList "Успешный запрос"
// @Failure 400 {object} view.ErrorResponse "Неверный формат запроса"
// @Failure 502 {object} view.ErrorResponse "Ошибка взаимодействия с сервисом"
// @Router /api/product/filter [post]
func (a *Apis) FilterProducts(c echo.Context) error {
	const op = "handlers.FilterProducts"
	var f domain.ProductFilter
	if err := c.Bind(&f); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	offset, limit, page, err := pageOffset(int(f.Page), int(f.Limit))
	if err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: err.Error()})
	}
	f.Offset = int32(offset)
	f.Limit = int32(limit)
	if f.Search != "" {
		f.Title = f.Search
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	list, err := a.apiProduct.API.FilterProducts(ctx, domain.ProductFilterToRpc(&f))
	if err != nil {
		code, body := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, body)
	}
	items := make([]view.ProductCatalogItem, 0, len(list.GetItems()))
	for _, item := range list.GetItems() {
		if item != nil {
			items = append(items, view.CatalogItemFromRPC(item))
		}
	}
	return c.JSON(http.StatusOK, view.ProductCatalogList{Items: items, Total: list.GetTotal(), Page: int32(page), Limit: int32(limit)})
}

// GetAllProducts godoc
// @Summary Получить все продукты
// @Description Возвращает список всех продуктов
// @Tags product
// @Produce json
// @Param page query int false "Page starting at 1"
// @Param limit query int true  "limit > 0"
// @Success 200 {object} view.ProductCatalogList
// @Failure 400 {object} view.ErrorResponse
// @Failure 502 {object} view.ErrorResponse "Ошибка взаимодействия с сервисом"
// @Router /api/product/all [get]
func (a *Apis) GetAllProducts(c echo.Context) error {
	const op = "handlers.GetAllProducts"

	ctx, cancel := requestContext(c)
	defer cancel()

	offset, limit, code, res := pag(c)
	if code != successPag {
		return c.JSON(code, res)
	}

	list, err := a.apiProduct.API.GetAllProducts(ctx, &productsRPC.Pagination{
		Offset: int32(offset),
		Limit:  int32(limit),
	})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}

	items := make([]view.ProductCatalogItem, 0, len(list.GetItems()))
	for _, item := range list.GetItems() {
		if item != nil {
			items = append(items, view.CatalogItemFromProduct(item))
		}
	}
	return c.JSON(http.StatusOK, view.ProductCatalogList{Items: items, Total: list.GetTotal(), Page: int32(offset/limit + 1), Limit: int32(limit)})
}
