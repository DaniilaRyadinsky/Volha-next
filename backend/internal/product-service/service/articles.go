package service

import (
	"context"
	"errors"
	"strings"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
)

func (s *ProductsService) articleRepo(ctx context.Context) (repository.ArticleRepo, error) {
	if s.repos == nil {
		return nil, errors.New("repository provider is nil")
	}
	repo := s.repos.Article(ctx)
	if repo == nil {
		return nil, errors.New("article repository is nil")
	}
	return repo, nil
}
func validateArticle(a *domain.Article) error {
	if a == nil {
		return errors.New("article is nil")
	}
	if err := validateID(a.Id); err != nil {
		return err
	}
	if strings.TrimSpace(a.Title) == "" || strings.TrimSpace(a.Slug) == "" {
		return errors.New("article title and slug are required")
	}
	if !strings.HasPrefix(a.Text, "/pdfs/") {
		return errors.New("article text must reference /pdfs/filename.pdf")
	}
	name := strings.TrimPrefix(a.Text, "/pdfs/")
	if strings.ContainsAny(name, "/\\:") || name == ".pdf" || !strings.HasSuffix(name, ".pdf") {
		return errors.New("invalid PDF link")
	}
	for _, image := range a.Img {
		if err := validateFilename(image); err != nil {
			return err
		}
	}
	return nil
}
func (s *ProductsService) GetAllArticles(ctx context.Context, offset, limit int) ([]*domain.Article, int, error) {
	if err := validateProductRange(offset, limit); err != nil {
		return nil, 0, wrapServiceCheck("articles.GetAll", err)
	}
	r, err := s.articleRepo(ctx)
	if err != nil {
		return nil, 0, wrapServiceCheck("articles.GetAll", err)
	}
	return r.GetAllArticles(ctx, offset, limit)
}
func (s *ProductsService) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	if err := validateID(id); err != nil {
		return nil, wrapServiceCheck("articles.Get", err)
	}
	r, err := s.articleRepo(ctx)
	if err != nil {
		return nil, wrapServiceCheck("articles.Get", err)
	}
	return r.GetArticle(ctx, id)
}
func (s *ProductsService) CreateArticle(ctx context.Context, a *domain.Article) error {
	if err := validateArticle(a); err != nil {
		return wrapServiceCheck("articles.Create", err)
	}
	return s.runInTx(ctx, "articles.Create", func(ctx context.Context) error {
		r, err := s.articleRepo(ctx)
		if err != nil {
			return wrapServiceCheck("articles.Create", err)
		}
		return r.CreateArticle(ctx, a)
	})
}
func (s *ProductsService) UpdateArticle(ctx context.Context, a *domain.Article) error {
	if err := validateArticle(a); err != nil {
		return wrapServiceCheck("articles.Update", err)
	}
	return s.runInTx(ctx, "articles.Update", func(ctx context.Context) error {
		r, err := s.articleRepo(ctx)
		if err != nil {
			return wrapServiceCheck("articles.Update", err)
		}
		return r.UpdateArticle(ctx, a)
	})
}
func (s *ProductsService) DeleteArticle(ctx context.Context, id string) error {
	if err := validateID(id); err != nil {
		return wrapServiceCheck("articles.Delete", err)
	}
	return s.runInTx(ctx, "articles.Delete", func(ctx context.Context) error {
		r, err := s.articleRepo(ctx)
		if err != nil {
			return wrapServiceCheck("articles.Delete", err)
		}
		return r.DeleteArticle(ctx, id)
	})
}
