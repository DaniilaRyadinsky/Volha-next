package handlers

import (
	"net/http"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

// CreateArticle godoc
// @Summary Создать новую статью
// @Description Добавляет новую статью. Id и CreationTime создаются автоматически. Можно передавать значения по умолчанию
// @Tags blog
// @Accept json
// @Produce json
// @Param article body productsRPC.Article true "Данные новой статьи"
// @Success 200 {object} view.IDResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 409 {object} view.ErrorResponse "уже существует"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/article [post]
func (a *Apis) CreateArticle(c echo.Context) error {
	const op = "handlers.CreateArticle"

	ctx, cancel := requestContext(c)
	defer cancel()

	var article productsRPC.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	article.Id = xid.New().String()
	if err := ensureImageReferences(article.GetImg()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.CreateArticle(ctx, &article); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, view.IDResponse{ID: article.Id})
}

// UpdateArticle godoc
// @Summary Обновить статью
// @Description Обновляет информацию о статье по ID
// @Tags blog
// @Accept json
// @Produce json
// @Param article body productsRPC.Article true "Данные статьи"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/article [put]
func (a *Apis) UpdateArticle(c echo.Context) error {
	const op = "handlers.UpdateArticle"

	ctx, cancel := requestContext(c)
	defer cancel()

	var article productsRPC.Article
	if err := c.Bind(&article); err != nil {
		return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "bad JSON"})
	}
	if err := ensureImageReferences(article.GetImg()); err != nil {
		code, swg := imageReferenceErrorResponse(ctx, op, err)
		return c.JSON(code, swg)
	}
	oldArticle, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: article.GetId()})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.UpdateArticle(ctx, &article); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if oldArticle.GetImg() != article.GetImg() {
		cleanupCtx, cleanupCancel := fileCleanupContext()
		defer cleanupCancel()
		a.cleanupFileIfUnreferenced(cleanupCtx, op, oldArticle.GetImg())
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "article updated successfully"})
}

// DeleteArticle godoc
// @Summary Удалить статью
// @Description Удаляет статью по ID
// @Tags blog
// @Produce json
// @Param id query string true "ID статью"
// @Success 200 {object} view.MessageResponse "успех"
// @Failure 400 {object} view.ErrorResponse "неправильный ввод"
// @Failure 404 {object} view.ErrorResponse "не найдено"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/article [delete]
func (a *Apis) DeleteArticle(c echo.Context) error {
	const op = "handlers.DeleteArticle"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	article, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if _, err := a.apiProduct.API.DeleteArticle(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	cleanupCtx, cleanupCancel := fileCleanupContext()
	defer cleanupCancel()
	a.cleanupFileIfUnreferenced(cleanupCtx, op, article.GetImg())
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "article deleted"})
}

// GetAllArticles godoc
// @Summary Получить все статьи
// @Description Возвращает список статей
// @Tags blog
// @Produce json
// @Success 200 {object} []productsRPC.Article "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/article/all [get]
func (a *Apis) GetAllArticles(c echo.Context) error {
	const op = "handlers.GetAllArticles"

	ctx, cancel := requestContext(c)
	defer cancel()

	list, err := a.apiProduct.API.GetAllArticles(ctx, nil)
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	if len(list.GetItems()) == 0 {
		return c.JSON(http.StatusOK, []*productsRPC.Article{})
	}
	return c.JSON(http.StatusOK, list.GetItems())
}

// GetArticle godoc
// @Summary Получить статью
// @Description Возвращает статью по id
// @Tags blog
// @Produce json
// @Param id query string true "ID статью"
// @Success 200 {object} productsRPC.Article "успех"
// @Failure 502 {object} view.ErrorResponse "ошибка в микросервисе"
// @Router /api/article [get]
func (a *Apis) GetArticle(c echo.Context) error {
	const op = "handlers.GetArticle"

	id, code, swg := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, swg)
	}
	ctx, cancel := requestContext(c)
	defer cancel()

	article, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, swg := a.mapGRPCError(ctx, op, err)
		return c.JSON(code, swg)
	}
	return c.JSON(http.StatusOK, article)
}
