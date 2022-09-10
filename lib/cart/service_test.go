package cart_test

import (
	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/cart"
)

var (
	service, err = cart.NewService(
		lib.NewEnv(
			lib.NewLogger(lib.EnvDev),
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
