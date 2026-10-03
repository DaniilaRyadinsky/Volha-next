package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateMaterial godoc
// @Summary Создать материал
// @Description Добавляет новый материал
// @Tags material
// @Accept json
// @Produce json
// @Param material body productsRPC.Material true "Новый материал"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/material [post]
func (a *Apis) CreateMaterial(c echo.Context) error {
	const op = "handlers.CreateMaterial"

	ctx, cancel := requestContext(c)
	defer cancel()

	var m productsRPC.Material
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	m.Id = xid.New().String()
	if _, err := a.apiProduct.API.CreateMaterial(ctx, &m); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.IDResponse{ID: m.Id})
}

// UpdateMaterial godoc
// @Summary Обновить материал
// @Description Обновляет данные существующего материала
// @Tags material
// @Accept json
// @Produce json
// @Param material body productsRPC.Material true "Обновлённые данные материала"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/material [put]
func (a *Apis) UpdateMaterial(c echo.Context) error {
	const op = "handlers.UpdateMaterial"

	ctx, cancel := requestContext(c)
	defer cancel()

	var m productsRPC.Material
	if err := c.Bind(&m); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if _, err := a.apiProduct.API.UpdateMaterial(ctx, &m); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "material updated successfully"})
}

// DeleteMaterial godoc
// @Summary Удалить материал
// @Description Удаляет материал по ID
// @Tags material
// @Produce json
// @Param id query string true "ID материала"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/material [delete]
func (a *Apis) DeleteMaterial(c echo.Context) error {
	const op = "handlers.DeleteMaterial"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	if _, err := a.apiProduct.API.DeleteMaterial(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	a.cleanDictionariesAfterMutation(ctx, op)
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "material deleted successfully"})
}

// GetAllMaterials godoc
// @Summary Получить все материалы
// @Description Возвращает список всех материалов
// @Tags material
// @Produce json
// @Success 200 {object} []productsRPC.Material "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/material/all [get]
func (a *Apis) GetAllMaterials(c echo.Context) error {
	const op = "handlers.GetAllMaterials"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllMaterials(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Material{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetMaterial rand godoc
// @Summary Получить материал
// @Description Возвращает материал по id
// @Tags material
// @Produce json
// @Param id query string true "ID материал"
// @Success 200 {object} productsRPC.Material "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/material [get]
func (a *Apis) GetMaterial(c echo.Context) error {
	const op = "handlers.GetMaterial"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	material, err := a.apiProduct.API.GetMaterial(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, material)
}
