package repository

import (
	"context"
	"database/sql"
	"errors"

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
	rows, err := d.Driver.QueryContext(ctx, `SELECT c.id, c.slug, c.title, c.description, COUNT(p.id)::int, c.img FROM categories c LEFT JOIN products p ON p.category_id = c.id GROUP BY c.id ORDER BY c.title ASC, c.id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Category, 0)
	for rows.Next() {
		var c domain.Category
		if err := rows.Scan(&c.Id, &c.Slug, &c.Title, &c.Description, &c.Count, &c.Img); err != nil {
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
	if err := d.Driver.QueryRowContext(ctx, `SELECT c.id, c.slug, c.title, c.description, COUNT(p.id)::int, c.img FROM categories c LEFT JOIN products p ON p.category_id = c.id WHERE c.id = $1 GROUP BY c.id`, id).Scan(&c.Id, &c.Slug, &c.Title, &c.Description, &c.Count, &c.Img); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &c, nil
}

func (d Driver) CreateCategory(ctx context.Context, c *domain.Category) error {
	const op = "psql.CreateCategory"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO categories (id, slug, title, img, description) VALUES ($1, $2, $3, $4, $5)`, c.Id, c.Slug, c.Title, c.Img, c.Description)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateCategory(ctx context.Context, c *domain.Category) error {
	const op = "psql.UpdateCategory"
	return d.execUpdate(op, ctx, `UPDATE categories SET slug = $2, title = $3, img = $4, description = $5 WHERE id = $1`, c.Id, c.Slug, c.Title, c.Img, c.Description)
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
	rows, err := d.Driver.QueryContext(ctx, `SELECT id, link, img, img762, title, description, badge FROM slides ORDER BY id ASC`)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	items := make([]*domain.Slide, 0)
	for rows.Next() {
		var s domain.Slide
		if err := rows.Scan(&s.Id, &s.Link, &s.Img, &s.Img762, &s.Title, &s.Description, &s.Badge); err != nil {
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
	if err := d.Driver.QueryRowContext(ctx, `SELECT id, link, img, img762, title, description, badge FROM slides WHERE id = $1`, id).Scan(&s.Id, &s.Link, &s.Img, &s.Img762, &s.Title, &s.Description, &s.Badge); err != nil {
		return nil, mapCatalogRowError(op, err)
	}
	return &s, nil
}

func (d Driver) CreateSlide(ctx context.Context, s *domain.Slide) error {
	const op = "psql.CreateSlide"
	_, err := d.Driver.ExecContext(ctx, `INSERT INTO slides (id, link, img, img762, title, description, badge) VALUES ($1, $2, $3, $4, $5, $6, $7)`, s.Id, s.Link, s.Img, s.Img762, s.Title, s.Description, s.Badge)
	return mapCatalogExecError(op, err)
}

func (d Driver) UpdateSlide(ctx context.Context, s *domain.Slide) error {
	const op = "psql.UpdateSlide"
	return d.execUpdate(op, ctx, `UPDATE slides SET link = $2, img = $3, img762 = $4, title = $5, description = $6, badge = $7 WHERE id = $1`, s.Id, s.Link, s.Img, s.Img762, s.Title, s.Description, s.Badge)
}

func (d Driver) DeleteSlide(ctx context.Context, id string) error {
	const op = "psql.DeleteSlide"
	return d.execUpdate(op, ctx, `DELETE FROM slides WHERE id = $1`, id)
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
