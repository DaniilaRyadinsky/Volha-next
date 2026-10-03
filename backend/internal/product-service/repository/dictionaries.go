package repository

import (
	"context"
	"database/sql"

	"strconv"
	"strings"

	"github.com/autumnterror/volha-backend/internal/platform/format"
	"github.com/autumnterror/volha-backend/internal/product-service/domain"
)

type DictRepo interface {
	GetDictionaries(ctx context.Context) (*domain.Dictionaries, error)
	GetDictionariesByCategoryID(ctx context.Context, idCat string) (*domain.Dictionaries, error)
}

// GetDictionaries извлекает все справочники для общего фильтра.
func (d Driver) GetDictionaries(ctx context.Context) (*domain.Dictionaries, error) {
	const op = "PostgresDb.GetDictionaries"

	rows, err := d.Driver.QueryContext(ctx, getDicQuery)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	return parseDictionaryRows(rows)
}

// GetDictionariesByCategoryID извлекает справочники, относящиеся к определенной категории товаров.
func (d Driver) GetDictionariesByCategoryID(ctx context.Context, idCat string) (*domain.Dictionaries, error) {
	const op = "PostgresDb.GetDictionariesByCategoryID"

	rows, err := d.Driver.QueryContext(ctx, getDicByCatQuery, idCat)
	if err != nil {
		return nil, format.Error(op, err)
	}
	defer rows.Close()

	return parseDictionaryRows(rows)
}

// parseDictionaryRows обрабатывает строки из SQL-запроса и преобразует их в структуру domain.Dictionaries.
func parseDictionaryRows(rows *sql.Rows) (*domain.Dictionaries, error) {
	const op = "parseDictionaryRows"

	result := &domain.Dictionaries{
		Brands:     []*domain.Brand{},
		Categories: []*domain.Category{},
		Countries:  []*domain.Country{},
		Materials:  []*domain.Material{},
		Colors:     []*domain.Color{},
		MinPrice:   0,
		MaxPrice:   0,
		MinWidth:   0,
		MaxWidth:   0,
		MinHeight:  0,
		MaxHeight:  0,
		MinDepth:   0,
		MaxDepth:   0,
	}

	for rows.Next() {
		var typ, id, field1, field2, field3 string
		if err := rows.Scan(&typ, &id, &field1, &field2, &field3); err != nil {
			return nil, format.Error(op, err)
		}

		switch typ {
		case "brand":
			result.Brands = append(result.Brands, &domain.Brand{Id: id, Title: field1})
		case "category":
			result.Categories = append(result.Categories, &domain.Category{Id: id, Title: field1, Uri: field2, Img: field3})
		case "country":
			result.Countries = append(result.Countries, &domain.Country{Id: id, Title: field1, Friendly: field2})
		case "material":
			result.Materials = append(result.Materials, &domain.Material{Id: id, Title: field1})
		case "color":
			result.Colors = append(result.Colors, &domain.Color{Id: id, Title: field1, Hex: field2})
		case "stats":
			minPrice, err := strconv.Atoi(field1)
			if err != nil {
				return nil, format.Error(op, err)
			}
			result.MinPrice = int32(minPrice)

			maxPrice, err := strconv.Atoi(field2)
			if err != nil {
				return nil, format.Error(op, err)
			}
			result.MaxPrice = int32(maxPrice)

			parts := strings.Split(field3, ",")
			if len(parts) != 6 {
				return nil, format.Error(op, strconv.ErrSyntax)
			}
			dimensions := make([]int, len(parts))
			for i, part := range parts {
				dimensions[i], err = strconv.Atoi(part)
				if err != nil {
					return nil, format.Error(op, err)
				}
			}
			result.MinWidth = int32(dimensions[0])
			result.MaxWidth = int32(dimensions[1])
			result.MinHeight = int32(dimensions[2])
			result.MaxHeight = int32(dimensions[3])
			result.MinDepth = int32(dimensions[4])
			result.MaxDepth = int32(dimensions[5])
		}
	}

	if err := rows.Err(); err != nil {
		return nil, format.Error(op, err)
	}

	return result, nil
}
