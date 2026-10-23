package handlers

import (
	"errors"
	"net/http"
	"strings"

	productsRPC "github.com/autumnterror/volha-backend/api/proto/gen"
	"github.com/autumnterror/volha-backend/internal/gateway/net/view"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/labstack/echo/v4"
	"github.com/rs/xid"
)

type ArticleList struct {
	Items []*domain.Article `json:"items"`
	Total int32             `json:"total"`
	Page  int32             `json:"page"`
	Limit int32             `json:"limit"`
}

func bindArticle(c echo.Context) (*domain.Article, error) {
	var article domain.Article
	if err := c.Bind(&article); err != nil {
		return nil, err
	}
	link, err := normalizePDFLink(article.Text)
	if err != nil {
		return nil, err
	}
	article.Text = link
	if err := ensurePDFReference(link); err != nil {
		return nil, err
	}
	if err := ensureImageReferences(article.Img...); err != nil {
		return nil, err
	}
	return &article, nil
}

// CreateArticle godoc
// @Summary Создать статью со ссылкой на загруженный PDF
// @Tags articles
// @Accept json
// @Produce json
// @Param article body domain.Article true "Статья; text содержит /pdfs/filename.pdf"
// @Success 200 {object} view.IDResponse
// @Failure 400 {object} view.ErrorResponse
// @Router /api/article [post]
func (a *Apis) CreateArticle(c echo.Context) error {
	article, err := bindArticle(c)
	if err != nil {
		return articleInputError(c, err)
	}
	article.Id = xid.New().String()
	ctx, cancel := requestContext(c)
	defer cancel()
	if _, err := a.apiProduct.API.CreateArticle(ctx, domain.ArticleToRpc(article)); err != nil {
		code, body := a.mapGRPCError(ctx, "CreateArticle", err)
		return c.JSON(code, body)
	}
	return c.JSON(http.StatusOK, view.IDResponse{ID: article.Id})
}

// UpdateArticle godoc
// @Summary Обновить статью, сохранив время создания
// @Tags articles
// @Accept json
// @Produce json
// @Param article body domain.Article true "Статья"
// @Success 200 {object} view.MessageResponse
// @Failure 400 {object} view.ErrorResponse
// @Router /api/article [put]
func (a *Apis) UpdateArticle(c echo.Context) error {
	article, err := bindArticle(c)
	if err != nil {
		return articleInputError(c, err)
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	old, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: article.Id})
	if err != nil {
		code, body := a.mapGRPCError(ctx, "UpdateArticle.Get", err)
		return c.JSON(code, body)
	}
	if _, err := a.apiProduct.API.UpdateArticle(ctx, domain.ArticleToRpc(article)); err != nil {
		code, body := a.mapGRPCError(ctx, "UpdateArticle", err)
		return c.JSON(code, body)
	}
	cleanupCtx, done := fileCleanupContext()
	defer done()
	for _, name := range removedStrings(old.GetImg(), article.Img) {
		a.cleanupFileIfUnreferenced(cleanupCtx, "UpdateArticle", name)
	}
	if old.GetText() != article.Text {
		a.cleanupFileIfUnreferenced(cleanupCtx, "UpdateArticle", strings.TrimPrefix(old.GetText(), "/pdfs/"))
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "article updated"})
}

func articleInputError(c echo.Context, err error) error {
	if errors.Is(err, ErrImageStorage) {
		return c.JSON(http.StatusInternalServerError, view.ErrorResponse{Error: "failed to validate article files"})
	}
	return c.JSON(http.StatusBadRequest, view.ErrorResponse{Error: "invalid article or missing files"})
}

// DeleteArticle godoc
// @Summary Удалить статью и неиспользуемые файлы
// @Tags articles
// @Produce json
// @Param id query string true "ID статьи"
// @Success 200 {object} view.MessageResponse
// @Router /api/article [delete]
func (a *Apis) DeleteArticle(c echo.Context) error {
	id, code, body := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, body)
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	old, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, body := a.mapGRPCError(ctx, "DeleteArticle.Get", err)
		return c.JSON(code, body)
	}
	if _, err := a.apiProduct.API.DeleteArticle(ctx, &productsRPC.Id{Id: id}); err != nil {
		code, body := a.mapGRPCError(ctx, "DeleteArticle", err)
		return c.JSON(code, body)
	}
	cleanupCtx, done := fileCleanupContext()
	defer done()
	for _, name := range append(old.GetImg(), strings.TrimPrefix(old.GetText(), "/pdfs/")) {
		a.cleanupFileIfUnreferenced(cleanupCtx, "DeleteArticle", name)
	}
	return c.JSON(http.StatusOK, view.MessageResponse{Message: "article deleted"})
}

// GetArticle godoc
// @Summary Получить статью
// @Tags articles
// @Produce json
// @Param id query string true "ID статьи"
// @Success 200 {object} domain.Article
// @Router /api/article [get]
func (a *Apis) GetArticle(c echo.Context) error {
	id, code, body := queryID(c)
	if code != http.StatusOK {
		return c.JSON(code, body)
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	article, err := a.apiProduct.API.GetArticle(ctx, &productsRPC.Id{Id: id})
	if err != nil {
		code, body := a.mapGRPCError(ctx, "GetArticle", err)
		return c.JSON(code, body)
	}
	return c.JSON(http.StatusOK, domain.ArticleFromRpc(article))
}

// GetAllArticles godoc
// @Summary Получить список статей
// @Tags articles
// @Produce json
// @Param page query int false "Страница с 1"
// @Param limit query int false "Размер страницы, от 1 до 100"
// @Success 200 {object} ArticleList
// @Router /api/article/all [get]
func (a *Apis) GetAllArticles(c echo.Context) error {
	offset, limit, code, body := pag(c)
	if code != successPag {
		return c.JSON(code, body)
	}
	ctx, cancel := requestContext(c)
	defer cancel()
	list, err := a.apiProduct.API.GetAllArticles(ctx, &productsRPC.Pagination{Offset: int32(offset), Limit: int32(limit)})
	if err != nil {
		code, body := a.mapGRPCError(ctx, "GetAllArticles", err)
		return c.JSON(code, body)
	}
	return c.JSON(http.StatusOK, ArticleList{Items: domain.ArticlesFromRpc(list.GetItems()), Total: list.GetTotal(), Page: int32(offset/limit + 1), Limit: int32(limit)})
}
