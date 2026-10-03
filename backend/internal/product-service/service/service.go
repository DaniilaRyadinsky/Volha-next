package service

import "github.com/autumnterror/volha-backend/internal/product-service/repository"

type ProductsService struct {
	tx    TxRunner
	repos repository.Provider
}

func NewProductsService(
	tx TxRunner,
	repos repository.Provider,
) *ProductsService {
	return &ProductsService{
		tx:    tx,
		repos: repos,
	}
}
