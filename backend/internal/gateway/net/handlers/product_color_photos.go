package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/types/known/emptypb"
)

// CreateProductColorPhotos godoc
// @Summary Создать фотографии продукта для цвета
// @Description Добавляет массив фотографий для конкретного продукта и цвета
// @Tags product_color_photos
// @Accept json
// @Produce json
// @Param photos body productsRPC.ProductColorPhotos true "Новый набор фотографий"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/colorphotos [post]
func (a *Apis) CreateProductColorPhotos(c echo.Context) error {
	const op = "handlers.CreateProductColorPhotos"

	ctx, cancel := requestContext(c)
	defer cancel()

	var pcp productsRPC.ProductColorPhotos
	if err := c.Bind(&pcp); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if err := ensureImageReferences(pcp.GetPhotos()...); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.CreateProductColorPhotos(ctx, &pcp); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "product color photos created successfully"})
}

// UpdateProductColorPhotos godoc
// @Summary Обновить фотографии продукта для цвета
// @Description Обновляет массив фотографий для конкретного продукта и цвета
// @Tags product_color_photos
// @Accept json
// @Produce json
// @Param photos body productsRPC.ProductColorPhotos true "Измененный набор фотографий"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/colorphotos [put]
func (a *Apis) UpdateProductColorPhotos(c echo.Context) error {
	const op = "handlers.UpdateProductColorPhotos"

	ctx, cancel := requestContext(c)
	defer cancel()

	var pcp productsRPC.ProductColorPhotos
	if err := c.Bind(&pcp); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if err := ensureImageReferences(pcp.GetPhotos()...); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldPhotos, err := a.apiProduct.API.GetPhotosByProductAndColor(ctx, &productsRPC.ProductColorPhotosId{
		ProductId: pcp.GetProductId(),
		ColorId:   pcp.GetColorId(),
	})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.UpdateProductColorPhotos(ctx, &pcp); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	for _, filename := range removedStrings(oldPhotos.GetItems(), pcp.GetPhotos()) {
		a.cleanupFileIfUnreferenced(cleanupCtx, op, filename)
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "product color photos updated successfully"})
}

// GetAllProductColorPhotos godoc
// @Summary Получить все фотографии продуктов по цветам
// @Description Возвращает все записи product_color_photos
// @Tags product_color_photos
// @Produce json
// @Success 200 {object} []productsRPC.ProductColorPhotos "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/colorphotos/all [get]
func (a *Apis) GetAllProductColorPhotos(c echo.Context) error {
	const op = "handlers.GetAllProductColorPhotos"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllProductColorPhotos(ctx, &emptypb.Empty{})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.ProductColorPhotos{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}
