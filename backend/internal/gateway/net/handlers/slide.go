package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateSlide godoc
// @Summary Создать слайд
// @Description Добавляет новый слайд
// @Tags slide
// @Accept json
// @Produce json
// @Param material body productsRPC.Slide true "Новый слайд"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/slide [post]
func (a *Apis) CreateSlide(c echo.Context) error {
	const op = "handlers.CreateSlide"

	ctx, cancel := requestContext(c)
	defer cancel()

	var s productsRPC.Slide
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	s.Id = xid.New().String()
	if err := ensureImageReferences(s.GetImg(), s.GetImg762()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.CreateSlide(ctx, &s); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, view.IDResponse{ID: s.Id})
}

// UpdateSlide godoc
// @Summary Обновить слайд
// @Description Обновляет данные существующего слайда
// @Tags slide
// @Accept json
// @Produce json
// @Param material body productsRPC.Slide true "Обновлённые данные слайда"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/slide [put]
func (a *Apis) UpdateSlide(c echo.Context) error {
	const op = "handlers.UpdateSlide"

	ctx, cancel := requestContext(c)
	defer cancel()

	var s productsRPC.Slide
	if err := c.Bind(&s); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if err := ensureImageReferences(s.GetImg(), s.GetImg762()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldSlide, err := a.apiProduct.API.GetSlide(ctx, &productsRPC.Id{Id: s.GetId()})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.UpdateSlide(ctx, &s); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	if oldSlide.GetImg() != s.GetImg() {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, oldSlide.GetImg())
	}
	if oldSlide.GetImg762() != s.GetImg762() {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, oldSlide.GetImg762())
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "slide updated successfully"})
}

// DeleteSlide godoc
// @Summary Удалить слайд
// @Description Удаляет слайд по ID
// @Tags slide
// @Produce json
// @Param id query string true "ID слайда"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/slide [delete]
func (a *Apis) DeleteSlide(c echo.Context) error {
	const op = "handlers.DeleteSlide"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	slide, err := a.apiProduct.API.GetSlide(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.DeleteSlide(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	a.cleanupFileIfUnreferenced(cleanupCtx, op, slide.GetImg())
	a.cleanupFileIfUnreferenced(cleanupCtx, op, slide.GetImg762())
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "slide deleted"})
}

// GetAllSlides godoc
// @Summary Получить все слайды
// @Description Возвращает список всех слайдов
// @Tags slide
// @Produce json
// @Success 200 {object} []productsRPC.Slide "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/slide/all [get]
func (a *Apis) GetAllSlides(c echo.Context) error {
	const op = "handlers.GetAllSlides"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllSlides(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Slide{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetSlide rand godoc
// @Summary Получить слайд
// @Description Возвращает слайд по id
// @Tags slide
// @Produce json
// @Param id query string true "ID слайд"
// @Success 200 {object} productsRPC.Slide "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/slide [get]
func (a *Apis) GetSlide(c echo.Context) error {
	const op = "handlers.GetSlide"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	slide, err := a.apiProduct.API.GetSlide(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, slide)
}
