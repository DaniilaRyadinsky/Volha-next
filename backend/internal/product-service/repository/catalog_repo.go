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

type CatalogRepo interface {
	GetAllBrands(ctx context.Context) ([]*domain.Brand, error)
	GetBrand(ctx context.Context, id string) (*domain.Brand, error)
	CreateBrand(ctx context.Context, b *domain.Brand) error
	UpdateBrand(ctx context.Context, b *domain.Brand) error
	DeleteBrand(ctx context.Context, id string) error

	GetAllCategories(ctx context.Context) ([]*domain.Category, error)
	GetCategory(ctx context.Context, id string) (*domain.Category, error)
	CreateCategory(ctx context.Context, c *domain.Category) error
	UpdateCategory(ctx context.Context, c *domain.Category) error
	DeleteCategory(ctx context.Context, id string) error

	GetAllCountries(ctx context.Context) ([]*domain.Country, error)
	GetCountry(ctx context.Context, id string) (*domain.Country, error)
	CreateCountry(ctx context.Context, c *domain.Country) error
	UpdateCountry(ctx context.Context, c *domain.Country) error
	DeleteCountry(ctx context.Context, id string) error

	GetAllMaterials(ctx context.Context) ([]*domain.Material, error)
	GetMaterial(ctx context.Context, id string) (*domain.Material, error)
	CreateMaterial(ctx context.Context, m *domain.Material) error
	UpdateMaterial(ctx context.Context, m *domain.Material) error
	DeleteMaterial(ctx context.Context, id string) error

	GetAllColors(ctx context.Context) ([]*domain.Color, error)
	GetColor(ctx context.Context, id string) (*domain.Color, error)
	CreateColor(ctx context.Context, c *domain.Color) error
	UpdateColor(ctx context.Context, c *domain.Color) error
	DeleteColor(ctx context.Context, id string) error

	GetAllSlides(ctx context.Context) ([]*domain.Slide, error)
	GetSlide(ctx context.Context, id string) (*domain.Slide, error)
	CreateSlide(ctx context.Context, s *domain.Slide) error
	UpdateSlide(ctx context.Context, s *domain.Slide) error
	DeleteSlide(ctx context.Context, id string) error

	GetAllArticles(ctx context.Context) ([]*domain.Article, error)
	GetArticle(ctx context.Context, id string) (*domain.Article, error)
	CreateArticle(ctx context.Context, a *domain.Article) error
	UpdateArticle(ctx context.Context, a *domain.Article) error
	DeleteArticle(ctx context.Context, id string) error
}

func (d Driver) GetAllBrands(ctx context.Context) ([]*domain.Brand, error) {
	const op = "psql.GetAllBrands"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title FROM brands ORDER BY title ASC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Brand, 0)
	for rows.Next() {
		var b domain.Brand
		if err := rows.Scan(&b.Id, &b.Title); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &b)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetBrand(ctx context.Context, id string) (*domain.Brand, error) {
	const op = "psql.GetBrand"
	var b domain.Brand
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, title FROM brands WHERE id = $1`, id).Scan(&b.Id, &b.Title); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &b, nil
}

func (d Driver) CreateBrand(ctx context.Context, b *domain.Brand) error {
	const op = "psql.CreateBrand"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO brands (id, title) VALUES ($1, $2)`, b.Id, b.Title)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateBrand(ctx context.Context, b *domain.Brand) error {
	const op = "psql.UpdateBrand"
	return d.execUpdate(op, ctx, `UPDATE brands SET title = $2 WHERE id = $1`, b.Id, b.Title)
}

func (d Driver) DeleteBrand(ctx context.Context, id string) error {
	const op = "psql.DeleteBrand"
	return d.execUpdate(op, ctx, `DELETE FROM brands WHERE id = $1`, id)
}

func (d Driver) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	const op = "psql.GetAllCategories"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title, uri, img FROM categories ORDER BY title ASC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Category, 0)
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.Id, &c.Title, &c.Uri, &c.Img); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetCategory(ctx context.Context, id string) (*domain.Category, error) {
	const op = "psql.GetCategory"
	var c domain.Category
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, title, uri, img FROM categories WHERE id = $1`, id).Scan(&c.Id, &c.Title, &c.Uri, &c.Img); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &c, nil
}

func (d Driver) CreateCategory(ctx context.Context, c *domain.Category) error {
	const op = "psql.CreateCategory"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO categories (id, title, uri, img) VALUES ($1, $2, $3, $4)`, c.Id, c.Title, c.Uri, c.Img)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateCategory(ctx context.Context, c *domain.Category) error {
	const op = "psql.UpdateCategory"
	return d.execUpdate(op, ctx, `UPDATE categories SET title = $2, uri = $3, img = $4 WHERE id = $1`, c.Id, c.Title, c.Uri, c.Img)
}

func (d Driver) DeleteCategory(ctx context.Context, id string) error {
	const op = "psql.DeleteCategory"
	return d.execUpdate(op, ctx, `DELETE FROM categories WHERE id = $1`, id)
}

func (d Driver) GetAllCountries(ctx context.Context) ([]*domain.Country, error) {
	const op = "psql.GetAllCountries"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title, friendly FROM countries ORDER BY title ASC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Country, 0)
	for rows.Next() {
		var c domain.Country
		if err := rows.Scan(&c.Id, &c.Title, &c.Friendly); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetCountry(ctx context.Context, id string) (*domain.Country, error) {
	const op = "psql.GetCountry"
	var c domain.Country
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, title, friendly FROM countries WHERE id = $1`, id).Scan(&c.Id, &c.Title, &c.Friendly); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &c, nil
}

func (d Driver) CreateCountry(ctx context.Context, c *domain.Country) error {
	const op = "psql.CreateCountry"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO countries (id, title, friendly) VALUES ($1, $2, $3)`, c.Id, c.Title, c.Friendly)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateCountry(ctx context.Context, c *domain.Country) error {
	const op = "psql.UpdateCountry"
	return d.execUpdate(op, ctx, `UPDATE countries SET title = $2, friendly = $3 WHERE id = $1`, c.Id, c.Title, c.Friendly)
}

func (d Driver) DeleteCountry(ctx context.Context, id string) error {
	const op = "psql.DeleteCountry"
	return d.execUpdate(op, ctx, `DELETE FROM countries WHERE id = $1`, id)
}

func (d Driver) GetAllMaterials(ctx context.Context) ([]*domain.Material, error) {
	const op = "psql.GetAllMaterials"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title FROM materials ORDER BY title ASC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Material, 0)
	for rows.Next() {
		var m domain.Material
		if err := rows.Scan(&m.Id, &m.Title); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &m)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetMaterial(ctx context.Context, id string) (*domain.Material, error) {
	const op = "psql.GetMaterial"
	var m domain.Material
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, title FROM materials WHERE id = $1`, id).Scan(&m.Id, &m.Title); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &m, nil
}

func (d Driver) CreateMaterial(ctx context.Context, m *domain.Material) error {
	const op = "psql.CreateMaterial"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO materials (id, title) VALUES ($1, $2)`, m.Id, m.Title)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateMaterial(ctx context.Context, m *domain.Material) error {
	const op = "psql.UpdateMaterial"
	return d.execUpdate(op, ctx, `UPDATE materials SET title = $2 WHERE id = $1`, m.Id, m.Title)
}

func (d Driver) DeleteMaterial(ctx context.Context, id string) error {
	const op = "psql.DeleteMaterial"
	return d.execUpdate(op, ctx, `DELETE FROM materials WHERE id = $1`, id)
}

func (d Driver) GetAllColors(ctx context.Context) ([]*domain.Color, error) {
	const op = "psql.GetAllColors"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title, hex FROM colors ORDER BY title ASC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Color, 0)
	for rows.Next() {
		var c domain.Color
		if err := rows.Scan(&c.Id, &c.Title, &c.Hex); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetColor(ctx context.Context, id string) (*domain.Color, error) {
	const op = "psql.GetColor"
	var c domain.Color
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, title, hex FROM colors WHERE id = $1`, id).Scan(&c.Id, &c.Title, &c.Hex); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &c, nil
}

func (d Driver) CreateColor(ctx context.Context, c *domain.Color) error {
	const op = "psql.CreateColor"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO colors (id, title, hex) VALUES ($1, $2, $3)`, c.Id, c.Title, c.Hex)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateColor(ctx context.Context, c *domain.Color) error {
	const op = "psql.UpdateColor"
	return d.execUpdate(op, ctx, `UPDATE colors SET title = $2, hex = $3 WHERE id = $1`, c.Id, c.Title, c.Hex)
}

func (d Driver) DeleteColor(ctx context.Context, id string) error {
	const op = "psql.DeleteColor"
	return d.execUpdate(op, ctx, `DELETE FROM colors WHERE id = $1`, id)
}

func (d Driver) GetAllSlides(ctx context.Context) ([]*domain.Slide, error) {
	const op = "psql.GetAllSlides"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, link, img, img762 FROM slides ORDER BY id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Slide, 0)
	for rows.Next() {
		var s domain.Slide
		if err := rows.Scan(&s.Id, &s.Link, &s.Img, &s.Img762); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &s)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetSlide(ctx context.Context, id string) (*domain.Slide, error) {
	const op = "psql.GetSlide"
	var s domain.Slide
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, link, img, img762 FROM slides WHERE id = $1`, id).Scan(&s.Id, &s.Link, &s.Img, &s.Img762); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &s, nil
}

func (d Driver) CreateSlide(ctx context.Context, s *domain.Slide) error {
	const op = "psql.CreateSlide"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO slides (id, link, img, img762) VALUES ($1, $2, $3, $4)`, s.Id, s.Link, s.Img, s.Img762)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateSlide(ctx context.Context, s *domain.Slide) error {
	const op = "psql.UpdateSlide"
	return d.execUpdate(op, ctx, `UPDATE slides SET link = $2, img = $3, img762 = $4 WHERE id = $1`, s.Id, s.Link, s.Img, s.Img762)
}

func (d Driver) DeleteSlide(ctx context.Context, id string) error {
	const op = "psql.DeleteSlide"
	return d.execUpdate(op, ctx, `DELETE FROM slides WHERE id = $1`, id)
}

func (d Driver) GetAllArticles(ctx context.Context) ([]*domain.Article, error) {
	const op = "psql.GetAllArticles"
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, title, img, text, created_at FROM articles ORDER BY created_at DESC, id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Article, 0)
	for rows.Next() {
		a, err := scanArticleRows(rows)
		if err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, a)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func (d Driver) GetArticle(ctx context.Context, id string) (*domain.Article, error) {
	const op = "psql.GetArticle"
	a, err := scanArticleRow(d.Driver.QueryRowContext(ctx, `SELECT id, title, img, text, created_at FROM articles WHERE id = $1`, id))
	if err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return a, nil
}

func (d Driver) CreateArticle(ctx context.Context, a *domain.Article) error {
	const op = "psql.CreateArticle"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO articles (id, title, img, text, created_at) VALUES ($1, $2, $3, $4, $5)`, a.Id, a.Title, a.Img, a.Text, time.Unix(a.CreationTime, 0).UTC())
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateArticle(ctx context.Context, a *domain.Article) error {
	const op = "psql.UpdateArticle"
	return d.execUpdate(op, ctx, `UPDATE articles SET title = $2, img = $3, text = $4 WHERE id = $1`, a.Id, a.Title, a.Img, a.Text)
}

func (d Driver) DeleteArticle(ctx context.Context, id string) error {
	const op = "psql.DeleteArticle"
	return d.execUpdate(op, ctx, `DELETE FROM articles WHERE id = $1`, id)
}

func scanArticleRows(rows *sql.Rows) (*domain.Article, error) {
	var a domain.Article
	var createdAt time.Time
	if err := rows.Scan(&a.Id, &a.Title, &a.Img, &a.Text, &createdAt); err != nil {
		return nil, err
	}
	a.CreationTime = createdAt.Unix()
	return &a, nil
}

func scanArticleRow(row *sql.Row) (*domain.Article, error) {
	var a domain.Article
	var createdAt time.Time
	if err := row.Scan(&a.Id, &a.Title, &a.Img, &a.Text, &createdAt); err != nil {
		return nil, err
	}
	a.CreationTime = createdAt.Unix()
	return &a, nil
}

func (d Driver) execUpdate(op string, ctx context.Context, query string, args ...any) error {
	res, err := d.Driver.ExecContext(ctx, query, args...)
	if err != nil {
		return mapCatalogExecError(op, err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return format.Error(op, err)
	}
	if rowsAffected == 0 {
		return format.Error(op, domain.ErrNotFound)
	}
	return nil
}

func mapCatalogRowError(op string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return format.Error(op, domain.ErrNotFound)
	}
	return format.Error(op, err)
}

func mapCatalogExecError(op string, err error) error {
	if err == nil {
		return nil
	}
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23505":
			return format.Error(op, domain.ErrAlreadyExists)
		case "23503":
			return format.Error(op, domain.ErrForeignKey)
		case "23514":
			return format.Error(op, domain.ErrInvalidType)
		}
	}
	return format.Error(op, err)
}
