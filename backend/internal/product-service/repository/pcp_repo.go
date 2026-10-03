package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/lib/pq"
)

type ProductColorPhotosRepo interface {
	GetPhotosByProductColor(ctx context.Context, productID string, colorID string) ([]string, error)
	GetAllProductColorPhotos(ctx context.Context) ([]domain.ProductColorPhotos, error)
	GetProductColorPhotosByProduct(ctx context.Context, productID string) ([]domain.ProductColorPhotos, error)
	CreateProductColorPhotos(ctx context.Context, pcp *domain.ProductColorPhotos) error
	UpdateProductColorPhotos(ctx context.Context, pcp *domain.ProductColorPhotos) error
	DeleteProductColorPhotos(ctx context.Context, productID string, colorID string) error
}

func (d Driver) GetPhotosByProductColor(
	ctx context.Context,
	productID string,
	colorID string,
) ([]string, error) {
	op := "psql.Get.ProductColorPhotos"
	query := `SELECT photos FROM product_color_photos WHERE product_id = $1 AND color_id = $2`

	row := d.Driver.QueryRowContext(ctx, query, productID, colorID)

	var res []string
	if err := row.Scan(pq.Array(&res)); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, format.Error(op, domain.ErrNotFound)
		}
		return nil, format.Error(op, err)
	}

	return res, nil
}

func (d Driver) GetAllProductColorPhotos(
	ctx context.Context,
) ([]domain.ProductColorPhotos, error) {
	op := "psql.GetAll.ProductColorPhotos"
	query := `SELECT product_id, color_id, photos FROM product_color_photos ORDER BY product_id ASC, color_id ASC`

	rows, err := d.Driver.QueryContext(ctx, query)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	var resLs []domain.ProductColorPhotos
	for rows.Next() {
		var res domain.ProductColorPhotos
		if err := rows.Scan(&res.ProductId, &res.ColorId, pq.Array(&res.Photos)); err != nil {
			return nil, format.Error(op, err)
		}
		resLs = append(resLs, res)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}

	return resLs, nil
}

func (d Driver) GetProductColorPhotosByProduct(
	ctx context.Context,
	productID string,
) ([]domain.ProductColorPhotos, error) {
	op := "psql.GetByProduct.ProductColorPhotos"
	query := `SELECT product_id, color_id, photos FROM product_color_photos WHERE product_id = $1 ORDER BY color_id ASC`

	rows, err := d.Driver.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	var resLs []domain.ProductColorPhotos
	for rows.Next() {
		var res domain.ProductColorPhotos
		if err := rows.Scan(&res.ProductId, &res.ColorId, pq.Array(&res.Photos)); err != nil {
			return nil, format.Error(op, err)
		}
		resLs = append(resLs, res)
	}
	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}

	return resLs, nil
}

func (d Driver) CreateProductColorPhotos(
	ctx context.Context,
	pcp *domain.ProductColorPhotos,
) error {
	op := "psql.Create.ProductColorPhotos"
	query := `INSERT INTO product_color_photos (product_id, color_id, photos) VALUES ($1, $2, $3)`

	_, err := d.Driver.ExecContext(ctx, query, pcp.ProductId, pcp.ColorId, pq.Array(pcp.Photos))
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505": // unique_violation
				return format.Error(op, domain.ErrAlreadyExists)
			case "23503": // foreign_key_violation
				return format.Error(op, domain.ErrForeignKey)
			}
		}
		return format.Error(op, err)
	}
	return nil
}

func (d Driver) UpdateProductColorPhotos(
	ctx context.Context,
	pcp *domain.ProductColorPhotos,
) error {
	op := "psql.Update.ProductColorPhotos"
	query := `UPDATE product_color_photos SET photos = $3 WHERE product_id = $1 AND color_id = $2`

	res, err := d.Driver.ExecContext(ctx, query, pcp.ProductId, pcp.ColorId, pq.Array(pcp.Photos))
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

func (d Driver) DeleteProductColorPhotos(
	ctx context.Context,
	productID string,
	colorID string,
) error {
	op := "psql.deleteProductColorPhotos"
	query := `DELETE FROM product_color_photos WHERE product_id = $1 AND color_id = $2`
	res, err := d.Driver.ExecContext(ctx, query, productID, colorID)
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
