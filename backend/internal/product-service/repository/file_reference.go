package repository

import (
	"context"

	"github.com/autumnterror/volha-backend/internal/platform/format"
)

type FileReferenceRepo interface {
	IsFileReferenced(ctx context.Context, filename string) (bool, error)
}

func (d Driver) IsFileReferenced(ctx context.Context, filename string) (bool, error) {
	const op = "psql.IsFileReferenced"
	query := `
SELECT EXISTS (
	SELECT 1 FROM categories WHERE img = $1
	UNION ALL
	SELECT 1 FROM products WHERE $1 = ANY(photos)
	UNION ALL
	SELECT 1 FROM product_color_photos WHERE $1 = ANY(photos)
	UNION ALL
	SELECT 1 FROM slides WHERE img = $1 OR img762 = $1
)`

	var referenced bool
	if err := d.Driver.QueryRowContext(ctx, query, filename).Scan(&referenced); err != nil {
		return false, format.Error(op, err)
	}

	return referenced, nil
}
