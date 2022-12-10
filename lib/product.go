package lib

import (
	"context"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/google/uuid"
)

type GetProductsOption struct {
	Sort *SortBy
}

type SortBy struct {
	Direction SortDirection
	Field     string
}

type SortDirection string

const (
	Descending SortDirection = "DESC"
	Ascending  SortDirection = "ASC"
)

type WithGetProductsOptions func(o *GetProductsOption) error

func WithProductSort(field string, direction SortDirection) WithGetProductsOptions {
	return func(o *GetProductsOption) error {
		o.Sort = &SortBy{
			Direction: direction,
			Field:     field,
		}
		return nil
	}
}

// ProductService ...
type ProductService interface {
	GetProduct(ctx context.Context, id uuid.UUID, conditions *GetProductCondtions) (*Product, error)
	DeleteProduct(ctx context.Context, product *Product, conditions *DeleteConditions) error
	GetProducts(ctx context.Context, opts ...WithGetProductsOptions) ([]*Product, error)
	SaveProduct(ctx context.Context, product *Product) (*Product, error)
}

// Product ...
type Product struct {
	Cost        float32
	Description string
	Name        string
	Inventory   int
	commons.Model
}

type DeleteConditions struct {
	HardDelete bool
}

type GetProductCondtions struct {
	Archived bool
}
