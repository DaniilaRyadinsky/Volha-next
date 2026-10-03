package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/autumnterror/volha-backend/internal/product-service/domain"
	"github.com/autumnterror/volha-backend/internal/product-service/mocks"
	"github.com/autumnterror/volha-backend/internal/product-service/service"
	"github.com/rs/xid"
)

func TestGetAllProductsUsesOffsetLimit(t *testing.T) {
	t.Parallel()

	repo := mocks.ProductRepo{
		GetAllProductsFn: func(_ context.Context, offset, limit int) ([]*domain.Product, int, error) {
			if offset != 20 || limit != 10 {
				t.Fatalf("unexpected pagination: offset=%d limit=%d", offset, limit)
			}
			return []*domain.Product{{Id: "product"}}, 1, nil
		},
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{ProductRepo: repo})

	items, total, err := svc.GetAllProducts(context.Background(), 20, 10)
	if err != nil {
		t.Fatalf("GetAllProducts() error = %v", err)
	}
	if total != 1 || len(items) != 1 {
		t.Fatalf("unexpected result: total=%d len=%d", total, len(items))
	}
}

func TestGetAllProductsRejectsBadPagination(t *testing.T) {
	t.Parallel()

	repo := mocks.ProductRepo{
		GetAllProductsFn: func(context.Context, int, int) ([]*domain.Product, int, error) {
			t.Fatal("repository must not be called")
			return nil, 0, nil
		},
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{ProductRepo: repo})

	_, _, err := svc.GetAllProducts(context.Background(), 0, 0)
	if err == nil {
		t.Fatal("expected validation error")
	}
	if !errors.Is(err, service.ErrBadServiceCheck) {
		t.Fatalf("expected ErrBadServiceCheck, got %v", err)
	}
}

func TestGetProductHasNoAnalyticsSideEffect(t *testing.T) {
	t.Parallel()

	id := xid.New().String()
	repo := mocks.ProductRepo{
		GetProductFn: func(_ context.Context, gotID string) (*domain.Product, error) {
			if gotID != id {
				t.Fatalf("unexpected id: %s", gotID)
			}
			return &domain.Product{Id: gotID}, nil
		},
	}
	svc := service.NewProductsService(mocks.TxRunner{}, mocks.Provider{ProductRepo: repo})

	product, err := svc.GetProduct(context.Background(), id)
	if err != nil {
		t.Fatalf("GetProduct() error = %v", err)
	}
	if product.Id != id {
		t.Fatalf("unexpected product id: %s", product.Id)
	}
}
