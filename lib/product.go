package lib

import (
	"context"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/google/uuid"
)

type GetProductsOption struct {
	ID       *uuid.UUID
	Sort     *SortBy
	Name     *string
	Archived bool
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
		if field == "" || direction == "" {
			return nil
		}
		o.Sort = &SortBy{
			Direction: direction,
			Field:     field,
		}
		return nil
	}
}

func WithProductName(name string) WithGetProductsOptions {
	return func(o *GetProductsOption) error {
		if name == "" {
			return nil
		}
		o.Name = &name
		return nil
	}
}

func WithProductID(id uuid.UUID) WithGetProductsOptions {
	return func(o *GetProductsOption) error {
		if id == uuid.Nil {
			return nil
		}
		o.ID = &id
		return nil
	}
}

func WithProductArchive() WithGetProductsOptions {
	return func(o *GetProductsOption) error {
		o.Archived = true
		return nil
	}
}

// ProductService ...
type ProductService interface {
	GetProduct(ctx context.Context, opts ...WithGetProductsOptions) (*Product, error)
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
