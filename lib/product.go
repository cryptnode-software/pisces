package lib

import "context"

//ProductService ...
type ProductService interface {
	GetProduct(ctx context.Context, id int64) (*Product, error)
}

//Product ...
type Product struct {
	Description string  `json:"description"`
	Cost        float32 `json:"cost"`
	Name        string  `json:"name"`
	ID          int64   `json:"id"`
}
