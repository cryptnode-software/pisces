package lib

import (
	"context"

	"github.com/google/uuid"
)

//ProductService ...
type ProductService interface {
	DeleteProduct(ctx context.Context, product *Product, conditions *ProductDeleteConditions) error
	GetProduct(ctx context.Context, id uuid.UUID, conditions *GetProductCondtions) (*Product, error)
	SaveProduct(ctx context.Context, product *Product) (*Product, error)
}

//Product ...
type Product struct {
	Cost        float32
	Description string
	Name        string
	Model
}

type ProductDeleteConditions struct {
	HardDelete bool
}

type GetProductCondtions struct {
	Archived bool
}
