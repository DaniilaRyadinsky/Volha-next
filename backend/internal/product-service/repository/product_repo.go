package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/lib/pq"
)

type ProductRepo interface {
	GetAllProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int, error)
	GetProduct(ctx context.Context, id string) (*domain.Product, error)
	CreateProduct(ctx context.Context, p *domain.ProductId) error
	UpdateProduct(ctx context.Context, p *domain.ProductId) error
	DeleteProduct(ctx context.Context, id string) error
	FilterProducts(ctx context.Context, filter *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error)
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (d Driver) GetAllProducts(ctx context.Context, offset, limit int) ([]*domain.Product, int, error) {
	const op = "PostgresDb.GetAllProducts"

	total, err := d.countProducts(ctx, op, nil, nil)
	if err != nil {
		return nil, 0, err
	}

	query := productListSelectQuery + " ORDER BY p.title ASC, p.id ASC LIMIT $1 OFFSET $2"
	rows, err := d.Driver.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, format.Error(op, err)
	}
	defer rows.Close()

	products, err := scanProductList(op, rows)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (d Driver) GetProduct(ctx context.Context, id string) (*domain.Product, error) {
	const op = "PostgresDb.GetProduct"

	query := productSelectQuery + " WHERE p.id = $1"
	p, err := scanProduct(d.Driver.QueryRowContext(ctx, query, id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, format.Error(op, err)
	}

	return p, nil
}

func (d Driver) CreateProduct(ctx context.Context, p *domain.ProductId) error {
	const op = "PostgresDb.CreateProduct"

	_, err := d.Driver.ExecContext(
		ctx,
		createProductQuery,
		p.Id, p.Title, p.Article, p.Brand, p.Category, p.Country, p.Width, p.Height, p.Depth,
		pq.Array(p.Photos), p.Price, p.Description, p.IsFeatured, p.PriceSale, p.Slug,
		pq.Array(p.Materials), pq.Array(p.Colors), pq.Array(p.Seems),
	)
	if err != nil {
		return mapPostgresError(op, err)
	}

	return nil
}

func (d Driver) UpdateProduct(ctx context.Context, p *domain.ProductId) error {
	const op = "PostgresDb.UpdateProduct"

	var updated int
	err := d.Driver.QueryRowContext(ctx, updateProductQuery,
		p.Id, p.Title, p.Article, p.Brand, p.Category, p.Country, p.Width, p.Height, p.Depth,
		pq.Array(p.Photos), p.Price, p.Description, p.IsFeatured, p.PriceSale, p.Slug,
		pq.Array(p.Materials), pq.Array(p.Colors), pq.Array(p.Seems),
	).Scan(&updated)
	if err != nil {
		return mapPostgresError(op, err)
	}
	if updated == 0 {
		return format.Error(op, domain.ErrNotFound)
	}

	return nil
}

func (d Driver) DeleteProduct(ctx context.Context, id string) error {
	const op = "PostgresDb.DeleteProduct"

	res, err := d.Driver.ExecContext(ctx, `DELETE FROM products WHERE id = $1`, id)
	if err != nil {
		return format.Error(op, err)
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

func (d Driver) FilterProducts(ctx context.Context, filter *domain.ProductFilter) ([]*domain.ProductCatalogItem, int, error) {
	const op = "PostgresDb.FilterProducts"

	conditions, args := buildProductFilter(filter)
	total, err := d.countProducts(ctx, op, conditions, args)
	if err != nil {
		return nil, 0, err
	}

	query := productCatalogSelectQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += productOrderBy(filter.SortBy, filter.SortOrder)
	limit := filter.Limit
	offset := filter.Offset
	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	rows, err := d.Driver.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, format.Error(op, err)
	}
	defer rows.Close()

	products, err := scanProductCatalogList(op, rows)
	if err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

func (d Driver) countProducts(ctx context.Context, op string, conditions []string, args []any) (int, error) {
	query := productCountQuery
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	var total int
	if err := d.Driver.QueryRowContext(ctx, query, args...).Scan(&total); err != nil {
		return 0, format.Error(op, err)
	}
	return total, nil
}

func buildProductFilter(filter *domain.ProductFilter) ([]string, []any) {
	var conditions []string
	var args []any
	argPos := 1

	buildFilter := func(field string, values []string) {
		if len(values) == 0 {
			return
		}
		placeholders := make([]string, len(values))
		for i, v := range values {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, v)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ",")))
	}

	buildFilter("p.brand_id", filter.Brand)
	buildFilter("p.category_id", filter.Category)
	buildFilter("p.country_id", filter.Country)

	numericFilters := []struct {
		field string
		value int
		op    string
	}{
		{"p.width", int(filter.MinWidth), ">="},
		{"p.width", int(filter.MaxWidth), "<="},
		{"p.height", int(filter.MinHeight), ">="},
		{"p.height", int(filter.MaxHeight), "<="},
		{"p.depth", int(filter.MinDepth), ">="},
		{"p.depth", int(filter.MaxDepth), "<="},
		{"p.price", int(filter.MinPrice), ">="},
		{"p.price", int(filter.MaxPrice), "<="},
	}

	for _, f := range numericFilters {
		if f.value > 0 {
			conditions = append(conditions, fmt.Sprintf("%s %s $%d", f.field, f.op, argPos))
			args = append(args, f.value)
			argPos++
		}
	}

	if filter.Title != "" {
		conditions = append(conditions, fmt.Sprintf("p.title ILIKE $%d", argPos))
		args = append(args, "%"+filter.Title+"%")
		argPos++
	}

	if len(filter.Materials) > 0 {
		placeholders := make([]string, len(filter.Materials))
		for i, v := range filter.Materials {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, v)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf(`
EXISTS (
	SELECT 1 FROM product_materials pm
	WHERE pm.product_id = p.id AND pm.material_id IN (%s)
)`, strings.Join(placeholders, ",")))
	}

	if len(filter.Colors) > 0 {
		placeholders := make([]string, len(filter.Colors))
		for i, v := range filter.Colors {
			placeholders[i] = fmt.Sprintf("$%d", argPos)
			args = append(args, v)
			argPos++
		}
		conditions = append(conditions, fmt.Sprintf(`
EXISTS (
	SELECT 1 FROM product_colors pc
	WHERE pc.product_id = p.id AND pc.color_id IN (%s)
)`, strings.Join(placeholders, ",")))
	}

	return conditions, args
}

func productOrderBy(sortBy, sortOrder string) string {
	validSort := map[string]string{
		"price":  "p.price",
		"width":  "p.width",
		"height": "p.height",
		"depth":  "p.depth",
		"title":  "p.title",
	}
	column, ok := validSort[sortBy]
	if !ok {
		return " ORDER BY p.title ASC, p.id ASC"
	}

	order := "ASC"
	if strings.ToLower(sortOrder) == "desc" {
		order = "DESC"
	}
	return fmt.Sprintf(" ORDER BY %s %s, p.id ASC", column, order)
}

func scanProductList(op string, rows *sql.Rows) ([]*domain.Product, error) {
	products := make([]*domain.Product, 0)
	for rows.Next() {
		p, err := scanProduct(rows)
		if err != nil {
			return nil, format.Error(op, err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}

	return products, nil
}

func scanProductCatalogList(op string, rows *sql.Rows) ([]*domain.ProductCatalogItem, error) {
	items := make([]*domain.ProductCatalogItem, 0)
	for rows.Next() {
		var item domain.ProductCatalogItem
		item.CategoryDetails = &domain.Category{}
		if err := rows.Scan(
			&item.Id,
			&item.Slug,
			&item.Category,
			&item.Title,
			&item.Price,
			&item.Image,
			&item.PriceSale,
			&item.IsFeatured,
			&item.CategoryDetails.Id, &item.CategoryDetails.Slug, &item.CategoryDetails.Title,
			&item.CategoryDetails.Description, &item.CategoryDetails.Count, &item.CategoryDetails.Img,
			pq.Array(&item.Photos),
		); err != nil {
			return nil, format.Error(op, err)
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}
	return items, nil
}

func scanProduct(row rowScanner) (*domain.Product, error) {
	p, fields, err := scanProductFields(row)
	if err != nil {
		return nil, err
	}
	if err := decodeProductJSON(p, fields); err != nil {
		return nil, err
	}
	return p, nil
}

type productJSONFields struct {
	materials []byte
	colors    []byte
	similar   []byte
}

func scanProductFields(row rowScanner, extra ...any) (*domain.Product, productJSONFields, error) {
	p := domain.NewEmptyProduct()
	fields := productJSONFields{}
	dest := []any{
		&p.Id,
		&p.Title,
		&p.Article,
		&p.Slug,
		&p.Width,
		&p.Height,
		&p.Depth,
		pq.Array(&p.Photos),
		&p.Price,
		&p.Description,
		&p.IsFeatured,
		&p.PriceSale,
		&p.Brand.Id,
		&p.Brand.Title,
		&p.Category.Id,
		&p.Category.Slug,
		&p.Category.Title,
		&p.Category.Description,
		&p.Category.Count,
		&p.Category.Img,
		&p.Country.Id,
		&p.Country.Title,
		&p.Country.Friendly,
		&fields.materials,
		&fields.colors,
		&fields.similar,
	}
	dest = append(dest, extra...)

	if err := row.Scan(dest...); err != nil {
		return nil, fields, err
	}
	return &p, fields, nil
}

func decodeProductJSON(p *domain.Product, fields productJSONFields) error {
	if err := json.Unmarshal(fields.materials, &p.Materials); err != nil {
		return err
	}
	if err := json.Unmarshal(fields.colors, &p.Colors); err != nil {
		return err
	}
	// The SQL projection has its own schema; public JSON tags must not change
	// how stored article/photos/discount fields are decoded.
	var similar []struct {
		ID          string   `json:"id"`
		Title       string   `json:"title"`
		Article     string   `json:"article"`
		Slug        string   `json:"slug"`
		Width       int32    `json:"width"`
		Height      int32    `json:"height"`
		Depth       int32    `json:"depth"`
		Photos      []string `json:"photos"`
		Price       int32    `json:"price"`
		Description string   `json:"description"`
		IsFeatured  bool     `json:"isFeatured"`
		PriceSale   int32    `json:"pricesale"`
	}
	if err := json.Unmarshal(fields.similar, &similar); err != nil {
		return err
	}
	p.Seems = make([]*domain.Product, 0, len(similar))
	for _, item := range similar {
		p.Seems = append(p.Seems, &domain.Product{
			Id: item.ID, Title: item.Title, Article: item.Article, Slug: item.Slug,
			Width: item.Width, Height: item.Height, Depth: item.Depth,
			Photos: item.Photos, Price: item.Price, Description: item.Description,
			IsFeatured: item.IsFeatured, PriceSale: item.PriceSale,
		})
	}
	return nil
}

func mapPostgresError(op string, err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code {
		case "23503":
			return format.Error(op, domain.ErrForeignKey)
		case "23505":
			return format.Error(op, domain.ErrAlreadyExists)
		case "23514":
			return format.Error(op, domain.ErrInvalidType)
		}
	}
	return format.Error(op, err)
}
