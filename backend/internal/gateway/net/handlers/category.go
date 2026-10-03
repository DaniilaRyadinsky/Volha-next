package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateCategory godoc
// @Summary Создать категорию
// @Description Добавляет новую категорию
// @Tags category
// @Accept json
// @Produce json
// @Param category body productsRPC.Category true "Новая категория"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/category [post]
func (a *Apis) CreateCategory(c echo.Context) error {
	const op = "handlers.CreateCategory"

	ctx, cancel := requestContext(c)
	defer cancel()

	var cat productsRPC.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	cat.Id = xid.New().String()
	if err := ensureImageReferences(cat.GetImg()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.CreateCategory(ctx, &cat); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: cat.Id})
}

// UpdateCategory godoc
// @Summary Обновить категорию
// @Description Обновляет информацию о категории
// @Tags category
// @Accept json
// @Produce json
// @Param category body productsRPC.Category true "Обновлённая категория"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/category [put]
func (a *Apis) UpdateCategory(c echo.Context) error {
	const op = "handlers.UpdateCategory"

	ctx, cancel := requestContext(c)
	defer cancel()

	var cat productsRPC.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if err := ensureImageReferences(cat.GetImg()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldCategory, err := a.apiProduct.API.GetCategory(ctx, &productsRPC.Id{Id: cat.GetId()})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.UpdateCategory(ctx, &cat); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	if oldCategory.GetImg() != cat.GetImg() {
		cleanupCtx, cleanupCancel := fileCleanupContext()
		defer cleanupCancel()
		a.cleanupFileIfUnreferenced(cleanupCtx, op, oldCategory.GetImg())
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "category updated successfully"})
}

// DeleteCategory godoc
// @Summary Удалить категорию
// @Description Удаляет категорию по ID
// @Tags category
// @Produce json
// @Param id query string true "ID категории"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/category [delete]
func (a *Apis) DeleteCategory(c echo.Context) error {
	const op = "handlers.DeleteCategory"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	category, err := a.apiProduct.API.GetCategory(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.DeleteCategory(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	a.cleanupFileIfUnreferenced(cleanupCtx, op, category.GetImg())
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "category deleted successfully"})
}

// GetAllCategories godoc
// @Summary Получить все категории
// @Description Возвращает список всех категорий
// @Tags category
// @Produce json
// @Success 200 {object} []productsRPC.Category "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/category/all [get]
func (a *Apis) GetAllCategories(c echo.Context) error {
	const op = "handlers.GetAllCategories"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllCategories(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Category{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetCategory godoc
// @Summary Получить категорию
// @Description Возвращает категорию по id
// @Tags category
// @Produce json
// @Param id query string true "ID категории"
// @Success 200 {object} productsRPC.Category "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/category [get]
func (a *Apis) GetCategory(c echo.Context) error {
	const op = "handlers.GetCategory"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	category, err := a.apiProduct.API.GetCategory(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, category)
}
