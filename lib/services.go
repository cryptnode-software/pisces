package lib

import (
	"context"

	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/gocraft/dbr/v2"
)

//Services ...
type Services struct {
	ProductService ProductService
	UploadService  UploadService
	PaypalService  PaypalService
	OrderService   OrderService
	AuthService    AuthService
	CartService    CartService
}

//We can put methods that interact with multiple services here as it is an intermediary.
//For example we don't necessarily want to expose GetTotal to the outside of our app,
//therefor we can add it here since we don't expose our internal services to the out side

//GetTotal get total
func (services *Services) GetTotal(ctx context.Context, order *Order) (total float32, err error) {
	if order == nil {
		return 0, errors.ErrCartOrderNotProvided
	}

	cart, err := services.CartService.GetCart(ctx, order)
	if err != nil && err != dbr.ErrNotFound {
		return
	}

	for _, content := range cart.Contents {
		product, err := services.ProductService.GetProduct(ctx, content.ProductID)
		if err != nil {
			return 0, err
		}

		total += (product.Cost * float32(content.Quantity))
	}

	return
}
