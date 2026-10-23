package handlers

import (
	"net/http"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateColor godoc
// @Summary Создать цвет
// @Description Добавляет новый цвет
// @Tags color
// @Accept json
// @Produce json
// @Param color body productsRPC.Color true "Новый цвет"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/color [post]
func (a *Apis) CreateColor(c echo.Context) error {
	const op = "handlers.CreateColor"

	ctx, cancel := requestContext(c)
	defer cancel()

	var clr productsRPC.Color
	if err := c.Bind(&clr); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	clr.Id = xid.New().String()
	if _, err := a.apiProduct.API.CreateColor(ctx, &clr); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: clr.Id})
}

// UpdateColor godoc
// @Summary Обновить цвет
// @Description Обновляет данные существующего цвета
// @Tags color
// @Accept json
// @Produce json
// @Param color body productsRPC.Color true "Обновлённый цвет"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/color [put]
func (a *Apis) UpdateColor(c echo.Context) error {
	const op = "handlers.UpdateColor"

	ctx, cancel := requestContext(c)
	defer cancel()

	var clr productsRPC.Color
	if err := c.Bind(&clr); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if _, err := a.apiProduct.API.UpdateColor(ctx, &clr); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "color updated successfully"})
}

// DeleteColor godoc
// @Summary Удалить цвет
// @Description Удаляет цвет по ID
// @Tags color
// @Produce json
// @Param id query string true "ID цвета"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/color [delete]
func (a *Apis) DeleteColor(c echo.Context) error {
	const op = "handlers.DeleteColor"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	if _, err := a.apiProduct.API.DeleteColor(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "color deleted successfully"})
}

// GetAllColors godoc
// @Summary Получить все цвета
// @Description Возвращает список всех доступных цветов
// @Tags color
// @Produce json
// @Success 200 {object} []domain.Color "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/color/all [get]
func (a *Apis) GetAllColors(c echo.Context) error {
	const op = "handlers.GetAllColors"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllColors(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*domain.Color{})
	}
	return c.JSON(http.StatusOK, domain.ColorsFromRpc(list.GetItems()))
}

// GetColor godoc
// @Summary Получить цвет
// @Description Возвращает цвет по id
// @Tags color
// @Produce json
// @Param id query string true "ID цвета"
// @Success 200 {object} domain.Color "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/color [get]
func (a *Apis) GetColor(c echo.Context) error {
	const op = "handlers.GetColor"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	color, err := a.apiProduct.API.GetColor(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, domain.ColorFromRpc(color))
}
