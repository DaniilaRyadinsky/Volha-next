package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateBrand godoc
// @Summary Создать новый бренд
// @Description Добавляет новый бренд
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body productsRPC.Brand true "Данные нового бренда"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/brand [post]
func (a *Apis) CreateBrand(c echo.Context) error {
	const op = "handlers.CreateBrand"

	ctx, cancel := requestContext(c)
	defer cancel()

	var b productsRPC.Brand
	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	b.Id = xid.New().String()
	if _, err := a.apiProduct.API.CreateBrand(ctx, &b); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: b.Id})
}

// UpdateBrand godoc
// @Summary Обновить бренд
// @Description Обновляет информацию о бренде по ID
// @Tags brands
// @Accept json
// @Produce json
// @Param brand body productsRPC.Brand true "Данные бренда"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/brand [put]
func (a *Apis) UpdateBrand(c echo.Context) error {
	const op = "handlers.UpdateBrand"

	ctx, cancel := requestContext(c)
	defer cancel()

	var b productsRPC.Brand
	if err := c.Bind(&b); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if _, err := a.apiProduct.API.UpdateBrand(ctx, &b); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "brand updated"})
}

// DeleteBrand godoc
// @Summary Удалить бренд
// @Description Удаляет бренд по ID
// @Tags brands
// @Produce json
// @Param id query string true "ID бренда"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/brand [delete]
func (a *Apis) DeleteBrand(c echo.Context) error {
	const op = "handlers.DeleteBrand"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	if _, err := a.apiProduct.API.DeleteBrand(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "brand deleted"})
}

// GetAllBrands godoc
// @Summary Получить все бренды
// @Description Возвращает список брендов
// @Tags brands
// @Produce json
// @Success 200 {object} []productsRPC.Brand "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/brand/all [get]
func (a *Apis) GetAllBrands(c echo.Context) error {
	const op = "handlers.GetAllBrands"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllBrands(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Brand{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetBrand godoc
// @Summary Получить бренд
// @Description Возвращает бренд по id
// @Tags brands
// @Produce json
// @Param id query string true "ID бренда"
// @Success 200 {object} productsRPC.Brand "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/brand [get]
func (a *Apis) GetBrand(c echo.Context) error {
	const op = "handlers.GetBrand"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	brand, err := a.apiProduct.API.GetBrand(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, brand)
}
