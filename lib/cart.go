package lib

import (
	"context"

	"github.com/google/uuid"
)

//CartService represents the structure that the cart service should be
//when we implement it in its functional form. i.e. lib/cart/service.go
type CartService interface {
	SaveProduct(ctx context.Context, order *Order, product *Product, action CartAction, quantity int) error
	SaveCart(ctx context.Context, cart Cart) (Cart, error)
	GetCart(context.Context, *Order) (Cart, error)
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

//Cart represents how we structure a cart in the database
//used to pair an assortment of products to a single order.
//currently the OrderID is required, but in the future you
//should be able to make an Cart w/o pairing it to an order.
type Cart []*CartContent

//CartContents is the internal structure of a cart. This
//includes the product which it is associated with the
//quantity of said product and the ID that it holds in
//the database. When updating a cart we currently erase
//the previous one and create a new one. In the future
//it will update the previous cart instead of erasing it.
type CartContent struct {
	ProductID uuid.UUID
	Product   *Product `gorm:"references:ID;"`
	OrderID   uuid.UUID
	Quantity  int64
	Model
}
