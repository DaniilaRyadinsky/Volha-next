package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/lib/pq"
)

type ArticleRepo interface {
	GetAllArticles(context.Context, int, int) ([]*domain.Article, int, error)
	GetArticle(context.Context, string) (*domain.Article, error)
	CreateArticle(context.Context, *domain.Article) error
	UpdateArticle(context.Context, *domain.Article) error
	DeleteArticle(context.Context, string) error
}

func scanArticle(row rowScanner) (*domain.Article, error) {
	a := &domain.Article{Img: []string{}}
	var created time.Time
	if err := row.Scan(&a.Id, &a.Slug, &a.Title, pq.Array(&a.Img), &a.Text, &created); err != nil {
		return nil, err
	}
	a.CreationTime = created.UnixMilli()
	return a, nil
}
func (d Driver) GetAllArticles(ctx context.Context, offset, limit int) ([]*domain.Article, int, error) {
	var total int
	if err := d.Driver.QueryRowContext(ctx, "SELECT COUNT(*) FROM articles").Scan(&total); err != nil {
		return nil, 0, format.Error("articles.Count", err)
	}
	rows, err := d.Driver.QueryContext(ctx, "SELECT id,slug,title,img,pdf,created_at FROM articles ORDER BY created_at DESC,id ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, 0, format.Error("articles.GetAll", err)
	}
	defer rows.Close()
	items := make([]*domain.Article, 0)
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}
func (d Driver) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	a, err := scanArticle(d.Driver.QueryRowContext(ctx, "SELECT id,slug,title,img,pdf,created_at FROM articles WHERE id=$1", id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	return a, err
}
func (d Driver) CreateArticle(ctx context.Context, a *domain.Article) error {
	var created time.Time
	err := d.Driver.QueryRowContext(ctx, "INSERT INTO articles(id,slug,title,img,pdf) VALUES($1,$2,$3,$4,$5) RETURNING created_at", a.Id, a.Slug, a.Title, pq.Array(append([]string{}, a.Img...)), a.Text).Scan(&created)
	if err != nil {
		return mapCatalogExecError("articles.Create", err)
	}
	a.CreationTime = created.UnixMilli()
	return nil
}
func (d Driver) UpdateArticle(ctx context.Context, a *domain.Article) error {
	return d.execUpdate("articles.Update", ctx, "UPDATE articles SET slug=$2,title=$3,img=$4,pdf=$5 WHERE id=$1", a.Id, a.Slug, a.Title, pq.Array(append([]string{}, a.Img...)), a.Text)
}
func (d Driver) DeleteArticle(ctx context.Context, id string) error {
	return d.execUpdate("articles.Delete", ctx, "DELETE FROM articles WHERE id=$1", id)
}
