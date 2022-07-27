package product

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/gocraft/dbr/v2"
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
func NewService(env *lib.Env) (*Service, error) {
	return &Service{
		env,
		&repo{
			env.DB,
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
func (s *Service) GetProduct(ctx context.Context, id int64) (*lib.Product, error) {
	return s.repo.GetProduct(ctx, id)
}

type repoi interface {
	GetProduct(ctx context.Context, id int64) (*lib.Product, error)
}

type repo struct {
	*dbr.Connection
}

func (r *repo) GetProduct(ctx context.Context, id int64) (*lib.Product, error) {
	sess := r.NewSession(nil)

	result := &lib.Product{}

	err := sess.Select("*").From(tables.products).Where("id = ?", id).LoadOneContext(ctx, result)
	if err != nil {
		return nil, err
	}

	if result.ID != id {
		return nil, &errors.ErrNoProductFound{
			ID: id,
		}
	}

	return result, nil
}
