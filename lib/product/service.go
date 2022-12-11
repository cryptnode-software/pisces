package product

import (
	"context"
	"fmt"

	"github.com/cryptnode-software/pisces/lib"
	"gorm.io/gorm"
)

var (
	//tables the various tables that are used within our
	//product service.
	tables = struct {
		products string
	}{
		products: "products",
	}
)

//NewService creates a new paypal service that satisfies the PaypalService interface
func NewService(env *lib.Env) (lib.ProductService, error) {
	return &Service{
		env,
		&repo{
			env.GormDB,
		},
	}, nil
}

//Service the product service the acts a proxy between the
//product table/sql repo. Requires the general libraries env
//structure. The repo can be over written only from within
//product portion of the project pecisis
type Service struct {
	*lib.Env
	repo repoi
}

//GetProduct supplies us with an easy way to fetch a product by its provided id.
//if no product is found in the database with the provided id an error will be
//raised.
//
//id [github.com/google/uuid.UUID] required
//
//conditions [github.com/cryptnode-software/pisces/lib.GetProductCondtions] optional

func (s *Service) GetProduct(ctx context.Context, opts ...lib.WithGetProductsOptions) (*lib.Product, error) {
	return s.repo.GetProduct(ctx, opts...)
}

func (s *Service) GetProducts(ctx context.Context, opts ...lib.WithGetProductsOptions) ([]*lib.Product, error) {
	return s.repo.GetProducts(ctx, opts...)
}

func (s *Service) SaveProduct(ctx context.Context, product *lib.Product) (result *lib.Product, err error) {
	return s.repo.SaveProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, product *lib.Product, conditions *lib.DeleteConditions) error {

	if conditions != nil && conditions.HardDelete {
		return s.repo.HardDelete(ctx, product)
	}

	return s.repo.SoftDelete(ctx, product)
}

type repoi interface {
	GetProducts(ctx context.Context, opts ...lib.WithGetProductsOptions) (products []*lib.Product, err error)
	SaveProduct(ctx context.Context, product *lib.Product) (*lib.Product, error)
	GetProduct(ctx context.Context, opts ...lib.WithGetProductsOptions) (*lib.Product, error)
	HardDelete(ctx context.Context, product *lib.Product) error
	SoftDelete(ctx context.Context, product *lib.Product) error
}

type repo struct {
	*gorm.DB
}

func (r *repo) GetProduct(ctx context.Context, opts ...lib.WithGetProductsOptions) (product *lib.Product, err error) {
	product = new(lib.Product)

	options := new(lib.GetProductsOption)
	for _, opt := range opts {
		opt(options)
	}

	if options.ID != nil {
		tx := r.DB

		if options.Archived {
			tx = tx.Unscoped()
		}

		err = tx.First(product, "id = ?", options.ID).Error

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return
	}

	if options.Name != nil {
		tx := r.DB

		if options.Archived {
			tx = tx.Unscoped()
		}

		err = tx.First(product, "name = ?", options.Name).Error

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return
	}

	return
}

func (r *repo) GetProducts(ctx context.Context, opts ...lib.WithGetProductsOptions) (products []*lib.Product, err error) {
	options := new(lib.GetProductsOption)
	for _, opt := range opts {
		opt(options)
	}

	products = make([]*lib.Product, 0)

	if options.Sort != nil {
		err = r.DB.Order(fmt.Sprintf("%s %s", options.Sort.Field, options.Sort.Direction)).Find(&products).Error
		return
	}

	err = r.DB.Find(&products).Error

	return
}

func (r *repo) SaveProduct(ctx context.Context, product *lib.Product) (*lib.Product, error) {
	db := r.DB.Save(product)
	return product, db.Error
}

func (r *repo) HardDelete(ctx context.Context, product *lib.Product) error {
	return r.DB.Unscoped().Delete(product).Error
}

func (r *repo) SoftDelete(ctx context.Context, product *lib.Product) error {
	return r.DB.Delete(product).Error
}
