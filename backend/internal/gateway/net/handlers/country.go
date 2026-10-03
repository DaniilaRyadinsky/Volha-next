package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateCountry godoc
// @Summary Создать страну
// @Description Добавляет новую страну
// @Tags country
// @Accept json
// @Produce json
// @Param country body productsRPC.Country true "Новая страна"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/country [post]
func (a *Apis) CreateCountry(c echo.Context) error {
	const op = "handlers.CreateCountry"

	ctx, cancel := requestContext(c)
	defer cancel()

	var ctr productsRPC.Country
	if err := c.Bind(&ctr); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	ctr.Id = xid.New().String()
	if _, err := a.apiProduct.API.CreateCountry(ctx, &ctr); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: ctr.Id})
}

// UpdateCountry godoc
// @Summary Обновить страну
// @Description Обновляет данные существующей страны
// @Tags country
// @Accept json
// @Produce json
// @Param country body productsRPC.Country true "Обновлённая страна"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/country [put]
func (a *Apis) UpdateCountry(c echo.Context) error {
	const op = "handlers.UpdateCountry"

	ctx, cancel := requestContext(c)
	defer cancel()

	var ctr productsRPC.Country
	if err := c.Bind(&ctr); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if _, err := a.apiProduct.API.UpdateCountry(ctx, &ctr); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "country updated successfully"})
}

// DeleteCountry godoc
// @Summary Удалить страну
// @Description Удаляет страну по ID
// @Tags country
// @Produce json
// @Param id query string true "ID страны"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/country [delete]
func (a *Apis) DeleteCountry(c echo.Context) error {
	const op = "handlers.DeleteCountry"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	if _, err := a.apiProduct.API.DeleteCountry(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "country deleted successfully"})
}

// GetAllCountries godoc
// @Summary Получить все страны
// @Description Возвращает список всех стран
// @Tags country
// @Produce json
// @Success 200 {object} []productsRPC.Country "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/country/all [get]
func (a *Apis) GetAllCountries(c echo.Context) error {
	const op = "handlers.GetAllCountries"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllCountries(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Country{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetCountry godoc
// @Summary Получить страну
// @Description Возвращает страну по id
// @Tags country
// @Produce json
// @Param id query string true "ID страны"
// @Success 200 {object} productsRPC.Country "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/country [get]
func (a *Apis) GetCountry(c echo.Context) error {
	const op = "handlers.GetCountry"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	country, err := a.apiProduct.API.GetCountry(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, country)
}
