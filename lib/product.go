package lib

import (
	"context"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/google/uuid"
)

// ProductService ...
type ProductService interface {
	DeleteProduct(ctx context.Context, product *Product, conditions *DeleteConditions) error
	GetProduct(ctx context.Context, id uuid.UUID, conditions *GetProductCondtions) (*Product, error)
	SaveProduct(ctx context.Context, product *Product) (*Product, error)
}

// Product ...
type Product struct {
	Cost        float32
	Description string
	Name        string
	commons.Model
}

type DeleteConditions struct {
	HardDelete bool
}

type GetProductCondtions struct {
	Archived bool
}
