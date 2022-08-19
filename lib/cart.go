package lib

import (
	"context"

	"github.com/google/uuid"
)

//CartService represents the structure that the cart service should be
//when we implement it in its functional form. i.e. lib/cart/service.go
type CartService interface {
	SaveProduct(ctx context.Context, order *Order, product *Product, action CartAction, quantity int) error
	SaveCart(ctx context.Context, cart []*Cart) ([]*Cart, error)
	GetCart(context.Context, *Order) ([]*Cart, error)
}

//CartAction represents the primitive type for all of the CartActions.
//This is used for add or removing a product from the provided order
type CartAction string

const (
	//AddProduct dispatches an action that lets our cart service know that it
	//needs to add a product from its collection
	AddProduct CartAction = "ADD"
	//RemoveProduct dispates an action that lets our cart service know that it
	//needs to remove a product from its collection
	RemoveProduct CartAction = "REMOVE"
)

type Cart struct {
	ProductID uuid.UUID
	Product   *Product `gorm:"references:ID;"`
	OrderID   uuid.UUID
	Quantity  int64
	Model
}
