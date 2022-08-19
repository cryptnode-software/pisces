package product_test

import (
	"context"
	"errors"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/product"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/stretchr/testify/assert"
)

var (
	env = utility.NewEnv(utility.NewLogger())

	service, err = product.NewService(env)

	ctx = context.Background()

	products = []*lib.Product{
		{
			Description: "Test One Product Description",
			Name:        "Test One Product",
			Cost:        0.00,
		},
		{
			Description: "Test Two Product Description",
			Name:        "Test Two Product",
			Cost:        0.00,
		},
		{
			Description: "Test Three Product Description",
			Name:        "Test Three Product",
			Cost:        0.00,
		},
	}
)

func TestGetProduct(t *testing.T) {

	if err != nil {
		t.Error(err)
		return
	}

	products := products

	if err := seed(products); err != nil {
		t.Error(err)
		return
	}

	for _, p := range products {
		product, err := service.GetProduct(ctx, p.ID, nil)

		if err != nil {
			t.Error(err)
			return
		}

		p.Model = product.Model

		assert.Equal(t, p, product)
	}

	if err = deseed(products); err != nil {
		t.Error(err)
		return
	}

}

func TestGetProductFailure(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	products := products

	for _, p := range products {

		product, err := service.GetProduct(ctx, p.ID, nil)

		if err != nil {
			t.Error(err)
			return
		}

		if product == nil {
			t.Error("product was successfully found when it should have failed")
			return
		}
	}

}

func TestSoftDeleteProduct(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	products := products

	if err := seed(products); err != nil {
		t.Error(err)
		return
	}

	for _, p := range products {

		err := service.DeleteProduct(ctx, p, &lib.DeleteConditions{
			HardDelete: false,
		})

		if err != nil {
			t.Error(err)
			return
		}

		product, err := service.GetProduct(ctx, p.ID, &lib.GetProductCondtions{
			Archived: true,
		})

		if product.ID != p.ID {
			t.Error(errors.New("product returned with a different id when trying to fetching an archived product"))
			return
		}

		p.Model = product.Model

		assert.Equal(t, p, product)
	}
}

func TestSaveProduct(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		product *lib.Product
	}{
		{
			product: &lib.Product{
				Cost:        40,
				Description: "A dozen cookies",
				Name:        "A dozen cookies",
			},
		},
	}

	for _, table := range tables {
		product, err := service.SaveProduct(ctx, table.product)
		if err != nil {
			t.Error(err)
			return
		}

		table.product = product

		assert.Equal(t, table.product, product)

	}

}

func seed(products []*lib.Product) error {
	for i, p := range products {
		product, err := service.SaveProduct(ctx, p)
		if err != nil {
			return err
		}

		products[i] = product
	}
	return nil
}

func deseed(products []*lib.Product) error {
	for _, p := range products {
		err := service.DeleteProduct(ctx, p, &lib.DeleteConditions{
			HardDelete: true,
		})
		if err != nil {
			return err
		}
	}
	return nil
}
