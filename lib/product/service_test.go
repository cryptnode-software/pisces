package product

import (
	"context"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/stretchr/testify/assert"
)

var (
	env = utility.NewEnv(utility.NewLogger())

	service, err = NewService(env)

	ctx = context.Background()
)

func TestGetProduct(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		id       int64
		expected *lib.Product
	}{
		{
			id: 1,
			expected: &lib.Product{
				Description: "Test One Product Description",
				Name:        "Test One Product",
				Cost:        0.00,
				ID:          1,
			},
		},
		{
			id: 2,
			expected: &lib.Product{
				Description: "Test Two Product Description",
				Name:        "Test Two Product",
				Cost:        0.00,
				ID:          2,
			},
		},
		{
			id: 3,
			expected: &lib.Product{
				Description: "Test Three Product Description",
				Name:        "Test Three Product",
				Cost:        0.00,
				ID:          3,
			},
		},
	}

	for _, table := range tables {
		product, err := service.GetProduct(ctx, table.id)

		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, table.expected, product)
	}
}
