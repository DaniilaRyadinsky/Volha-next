package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/mocks"
	"github.com/autumnterror/volha-backend/internal/product-service/repository"
	"github.com/autumnterror/volha-backend/internal/product-service/service"
	"github.com/rs/xid"
)

type articleCatalogRepo struct {
	repository.CatalogRepo
	createArticleFn func(context.Context, *domain.Article) error
}

func (r articleCatalogRepo) CreateArticle(ctx context.Context, a *domain.Article) error {
	return r.createArticleFn(ctx, a)
}

func TestCreateArticleOwnsCreationTime(t *testing.T) {
	t.Parallel()

	before := time.Now().UTC().Unix()
	var got int64
	repo := articleCatalogRepo{
		createArticleFn: func(_ context.Context, article *domain.Article) error {
			got = article.CreationTime
			return nil
		},
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{CatalogRepo: repo})

	err := svc.CreateArticle(context.Background(), &domain.Article{
		Id:           xid.New().String(),
		Img:          "article.jpg",
		Title:        "Article",
		Text:         "Body",
		CreationTime: 1,
	})
	after := time.Now().UTC().Unix()

	if err != nil {
		t.Fatalf("CreateArticle() error = %v", err)
	}
	if got < before || got > after {
		t.Fatalf("CreationTime = %d, want between %d and %d", got, before, after)
	}
}
