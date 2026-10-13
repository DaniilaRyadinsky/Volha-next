package domain

import "errors"

var (
	ErrInvalidInput    = errors.New("invalid input")
	ErrAlreadyExists   = errors.New("already exists")
	ErrUnauthenticated = errors.New("unauthenticated")
	ErrNotFound        = errors.New("not found")
	ErrCartLimit       = errors.New("cart quantity limit exceeded")
)

type User struct {
	ID    string
	Email string
	Login string
}

type CartItem struct {
	ProductID string
	Slug      string
	Title     string
	Quantity  int32
	Price     int32
	PriceSale int32
	Image     string
	LineTotal int64
}

type Cart struct {
	Items         []*CartItem
	TotalQuantity int32
	TotalPrice    int64
}
