package cart

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/gocraft/dbr/v2"
)

var (
	tables = struct {
		cart string
	}{
		cart: "carts",
	}
)

//Service the cart service handles any interactions/validations that we
//might need for the cart table and logic that should be handled accordingly
type Service struct {
	*lib.Env
	repo repoi
}

//NewService simply creates a new CartService to handle any sort of validation
//that we might need to interact with the cart table.
func NewService(env *lib.Env) (lib.CartService, error) {
	return &Service{
		env,
		&repo{
			env.DB,
		},
	}, nil
}

//SaveProduct validates and delegates the tasks required to add/remove a product to/from an order
func (service *Service) SaveProduct(ctx context.Context, order *lib.Order, product *lib.Product, action lib.CartAction, quantity int) error {

	if order == nil {
		return errors.ErrCartOrderNotProvided
	}

	if product == nil {
		return errors.ErrProductNotProvided
	}

	var err error
	switch action {
	case lib.RemoveProduct:
		err = service.repo.RemoveProduct(ctx, order, product)
	case lib.AddProduct:
		err = service.repo.AddProduct(ctx, order, product, quantity)
	default:
		return &errors.ErrCartActionNotRecognized{
			Action: string(action),
		}
	}

	return err
}

//SaveCart ...
func (service *Service) SaveCart(ctx context.Context, cart *lib.Cart) (*lib.Cart, error) {
	return service.repo.SaveCart(ctx, cart)
}

//GetCart simply accepts an order and returns (if any) products that are associated with it
//if non are found then it will return a nil value
func (service *Service) GetCart(ctx context.Context, order *lib.Order) (*lib.Cart, error) {
	if order == nil {
		return nil, errors.ErrCartOrderNotProvided
	}

	return service.repo.GetCart(ctx, order)
}

type repoi interface {
	AddProduct(ctx context.Context, order *lib.Order, product *lib.Product, quantity int) error
	SaveCart(ctx context.Context, cart *lib.Cart) (*lib.Cart, error)
	RemoveProduct(ctx context.Context, order *lib.Order, product *lib.Product) error
	GetCart(context.Context, *lib.Order) (*lib.Cart, error)
}

type repo struct {
	*dbr.Connection
}

//RemoveProduct gives us a simple way of directly writing to the cart table w/o any validation
//other than the ones that are within the dbr module itself. If you want any sort of validation
//you should do it in the service itself
func (repo *repo) RemoveProduct(ctx context.Context, order *lib.Order, product *lib.Product) error {
	sess := repo.NewSession(nil)

	_, err := sess.DeleteFrom(tables.cart).Where("order_id = ?", order.ID).Where("product_id = ?", product.ID).ExecContext(ctx)

	return err
}

//AddProduct: tldr; adds a product to the provided order directly into the cart table.
//gives us a simple way of directly writing to the cart table w/o any validation
//other than the ones that are within the dbr module itself. If you want any sort of validation
//you should do it in the service itself
func (repo *repo) AddProduct(ctx context.Context, order *lib.Order, product *lib.Product, quantity int) error {

	sess := repo.NewSession(nil)

	if rows, err := sess.Select("*").From(tables.cart).Where("order_id = ?", order.ID).Where("product_id = ?", product.ID).ReturnInt64(); rows >= 0 || err == nil {

		_, err := sess.Update(tables.cart).
			Where("order_id = ?", order.ID).
			Where("product_id = ?", product.ID).
			Set("quantity", quantity).
			ExecContext(ctx)

		return err
	}
	_, err := sess.InsertInto(tables.cart).
		Pair("order_id", order.ID).
		Pair("product_id", product.ID).
		Pair("quantity", quantity).
		ExecContext(ctx)

	return err

}

//GetCart accepts an entire order and returns any products and the quantity that have been
//added to the order.
func (repo *repo) GetCart(ctx context.Context, order *lib.Order) (cart *lib.Cart, err error) {
	sess := repo.NewSession(nil)
	cart = new(lib.Cart)

	_, err = sess.Select("*").From(tables.cart).Where("order_id = ?", order.ID).LoadContext(ctx, &cart.Contents)

	cart.OrderID = order.ID

	return
}

func (repo *repo) SaveCart(ctx context.Context, cart *lib.Cart) (*lib.Cart, error) {
	sess := repo.NewSession(nil)

	sess.DeleteFrom(tables.cart).Where("order_id = ?", cart.OrderID).ExecContext(ctx)

	for _, content := range cart.Contents {
		result, err := sess.InsertInto(tables.cart).
			Pair("product_id", content.ProductID).
			Pair("quantity", content.Quantity).
			Pair("order_id", cart.OrderID).
			ExecContext(ctx)

		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()

		if err != nil {
			return nil, err
		}

		content.ID = id
	}

	return cart, nil
}
