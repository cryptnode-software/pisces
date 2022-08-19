package cart

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"gorm.io/gorm"
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
			env.GormDB,
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
func (service *Service) SaveCart(ctx context.Context, cart []*lib.Cart) ([]*lib.Cart, error) {
	return service.repo.SaveCart(ctx, cart)
}

//GetCart simply accepts an order and returns (if any) products that are associated with it
//if non are found then it will return a nil value
func (service *Service) GetCart(ctx context.Context, order *lib.Order) ([]*lib.Cart, error) {
	if order == nil {
		return nil, errors.ErrCartOrderNotProvided
	}

	return service.repo.GetCart(ctx, order)
}

type repoi interface {
	AddProduct(ctx context.Context, order *lib.Order, product *lib.Product, quantity int) error
	SaveCart(ctx context.Context, cart []*lib.Cart) ([]*lib.Cart, error)
	RemoveProduct(ctx context.Context, order *lib.Order, product *lib.Product) error
	GetCart(context.Context, *lib.Order) ([]*lib.Cart, error)
}

type repo struct {
	*gorm.DB
}

//RemoveProduct gives us a simple way of directly writing to the cart table w/o any validation
//other than the ones that are within the grom module itself. If you want any sort of validation
//you should do it in the service itself
func (repo *repo) RemoveProduct(ctx context.Context, order *lib.Order, product *lib.Product) error {
	return repo.DB.Delete(new(lib.Cart), "order_id = ?", order.ID, "product_id = ?", product.ID).Error
}

//AddProduct: tldr; adds a product to the provided order directly into the cart table.
//gives us a simple way of directly writing to the cart table w/o any validation
//other than the ones that are within the gorm module itself. If you want any sort of validation
//you should do it in the service itself
func (repo *repo) AddProduct(ctx context.Context, order *lib.Order, product *lib.Product, quantity int) error {
	return repo.DB.Save(new(lib.Cart)).Error
}

//GetCart accepts an entire order and returns any products and the quantity that have been
//added to the order.
func (repo *repo) GetCart(ctx context.Context, order *lib.Order) (cart []*lib.Cart, err error) {

	repo.DB.Model(new(lib.Cart)).Find(cart, "order_id = ?", order.ID)

	return
}

func (repo *repo) SaveCart(ctx context.Context, cart []*lib.Cart) (result []*lib.Cart, err error) {

	err = repo.DB.Save(cart).Error

	return cart, err
}
