package cart_test

import (
	"github.com/cryptnode-software/pisces/lib/cart"
	"github.com/cryptnode-software/pisces/lib/utility"
)

var (
	service, err = cart.NewService(
		utility.NewEnv(
			utility.NewLogger(),
		),
	)
)

// func TestGetCart(t *testing.T) {
// 	tables := []struct {
// 		order *lib.Order
// 		cart  *lib.Cart
// 		fail  bool
// 	}{
// 		{
// 			fail:  true,
// 			order: &lib.Order{

// 			},
// 			cart: &lib.Cart{
// 				Product: &lib.Product{},
// 			},
// 		},
// 	}
// }
