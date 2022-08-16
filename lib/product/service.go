package product

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/google/uuid"
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

func (s *Service) GetProduct(ctx context.Context, id uuid.UUID, conditions *lib.GetProductCondtions) (*lib.Product, error) {
	if conditions != nil && conditions.Archived {
		return s.repo.GetArchivedProduct(ctx, id)
	}
	return s.repo.GetProduct(ctx, id)
}

func (s *Service) SaveProduct(ctx context.Context, product *lib.Product) (result *lib.Product, err error) {
	return s.repo.SaveProduct(ctx, product)
}

func (s *Service) DeleteProduct(ctx context.Context, product *lib.Product, conditions *lib.ProductDeleteConditions) error {

	if conditions != nil && conditions.HardDelete {
		return s.repo.HardDelete(ctx, product)
	}

	return s.repo.SoftDelete(ctx, product)
}

type repoi interface {
	SaveProduct(ctx context.Context, product *lib.Product) (*lib.Product, error)
	GetArchivedProduct(ctx context.Context, id uuid.UUID) (*lib.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*lib.Product, error)
	HardDelete(ctx context.Context, product *lib.Product) error
	SoftDelete(ctx context.Context, product *lib.Product) error
}

type repo struct {
	*gorm.DB
}

func (r *repo) GetProduct(ctx context.Context, id uuid.UUID) (product *lib.Product, err error) {
	product = new(lib.Product)
	r.DB.Model(new(lib.Product)).First(product, "id = ?", id)
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

func (r *repo) GetArchivedProduct(ctx context.Context, id uuid.UUID) (product *lib.Product, err error) {
	product = new(lib.Product)
	err = r.DB.Unscoped().Where("id = ?", id).Find(product).Error
	return product, err
}
