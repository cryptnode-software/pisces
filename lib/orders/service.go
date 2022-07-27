package orders

import (
	"context"
	"strings"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/gocraft/dbr/v2"
)

var (
	tables = struct {
		orders      string
		inquires    string
		attachments string
	}{
		attachments: "inquiry_attachments",
		inquires:    "inquires",
		orders:      "orders",
	}
)

//Service the order service, handles everything related to an order
type Service struct {
	*lib.Env
	repo repoi
}

//NewService returns a new `Orders` service to handle every
//thing related to an order
func NewService(env *lib.Env) (lib.OrderService, error) {
	return &Service{
		env,
		&repo{
			env.DB,
		},
	}, nil
}

//GetOrders returns orders sorted and filtered by the conditions provided
func (s *Service) GetOrders(ctx context.Context, conditions *lib.OrderConditions) ([]*lib.Order, error) {
	return s.repo.GetOrders(ctx, conditions)
}

//GetOrder returns a specific order and any information that is associated with
//it
func (s *Service) GetOrder(ctx context.Context, id int64) (*lib.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

//GetInquiry returns a specific inquiry based on the id provided, if there is
//no inquiry found an exception will be raised.
func (s *Service) GetInquiry(ctx context.Context, id int64) (*lib.Inquiry, error) {
	return s.repo.GetInquiry(ctx, id)
}

//SaveOrder will either create a new order (if required, i.e. the ID property is empty). Or
//will update a preexisting one.
func (s *Service) SaveOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {

	//inquiry should be required to create/update a order
	if order.InquiryID == 0 {
		return nil, &errors.ErrNoOrderInquiryProvided{
			OrderID: order.ID,
		}
	}

	//create new order
	if order.ID == 0 {
		return s.repo.CreateOrder(ctx, order)
	}

	//otherwise update a preexisting order
	return s.repo.UpdateOrder(ctx, order)
}

//SaveInquiry will either create a new inquiry or update a pre-existing one. The optional
//parameter that defines whether it will create/update one is dependant on the ID. If the
//ID <= zero (basically zero but in the off chance someone tries to update an id that is
//less then 0, we can catch it) create a new inquiry, otherwise update with the id that
//is provided.
func (s *Service) SaveInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error) {
	if inquiry.ID <= 0 {
		return s.repo.CreateInquiry(ctx, inquiry)
	}

	return s.repo.UpdateInquiry(ctx, inquiry)
}

//GetInquires returns all of the inquires that match the *optional conditions from
//the second parameter. If there are no conditions provided and/or met then all
//inquires will be returned.
func (s *Service) GetInquires(ctx context.Context, conditions *lib.GetInquiryConditions) ([]*lib.Inquiry, error) {
	return s.repo.GetInquires(ctx, conditions)
}

//ArchiveOrder archives a provided order
func (s *Service) ArchiveOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {
	return order, nil
}

type repoi interface {
	GetInquires(ctx context.Context, conditions *lib.GetInquiryConditions) ([]*lib.Inquiry, error)
	UpdateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error)
	CreateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error)
	CreateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error)
	UpdateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error)
	GetInquiry(ctx context.Context, id int64) (*lib.Inquiry, error)
	GetOrder(ctx context.Context, id int64) (*lib.Order, error)
	GetOrders(context.Context, *lib.OrderConditions) ([]*lib.Order, error)
}

type repo struct {
	*dbr.Connection
}

func (r *repo) GetInquiry(ctx context.Context, id int64) (inquiry *lib.Inquiry, err error) {
	sess := r.NewSession(nil)

	err = sess.Select("*").From(tables.inquires).Where("id = ?", id).LoadOneContext(ctx, &inquiry)

	return
}

func (r *repo) CreateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {
	sess := r.NewSession(nil)

	result, err := sess.InsertInto(tables.orders).
		Pair("payment_method", order.PaymentMethod).
		Pair("inquiry_id", order.InquiryID).
		Pair("status", order.Status).
		Pair("ext_id", order.ExtID).
		Pair("due", order.Due).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	order.ID = id

	_, err = sess.Update(tables.inquires).
		Where("id = ?", order.InquiryID).
		Set("order_id", order.ID).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *repo) UpdateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {
	sess := r.NewSession(nil)

	_, err := sess.Update(tables.orders).
		Where("id = ?", order.ID).
		Set("payment_method", order.PaymentMethod).
		Set("status", order.Status).
		Set("ext_id", order.ExtID).
		Set("due", order.Due).
		Exec()

	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *repo) GetOrders(ctx context.Context, conditions *lib.OrderConditions) ([]*lib.Order, error) {
	sess := r.NewSession(nil)

	var result []*lib.Order

	stmt := sess.Select("*").
		From(tables.orders)

	if conditions != nil {
		if conditions.Status != lib.OrderStatusNotImplemented {
			stmt.Where("status = ?", conditions.Status)
		}
	}

	_, err := stmt.LoadContext(ctx, &result)

	if err != nil {
		return nil, err
	}

	for i, order := range result {
		result[i].Due = strings.Replace(order.Due, "T", " ", 1)
		result[i].Due = strings.Replace(order.Due, "Z", " ", 1)
	}

	return result, nil
}

func (r *repo) GetOrder(ctx context.Context, id int64) (order *lib.Order, err error) {
	sess := r.NewSession(nil)

	err = sess.Select("*").From(tables.orders).Where("id = ?", id).LoadOne(&order)

	order.Due = strings.Replace(order.Due, "T", " ", 1)
	order.Due = strings.Replace(order.Due, "Z", " ", 1)

	return
}

func (r *repo) GetInquires(ctx context.Context, conditions *lib.GetInquiryConditions) ([]*lib.Inquiry, error) {
	sess := r.NewSession(nil)

	var result []*lib.Inquiry

	stmt := sess.Select("*").From(tables.inquires)

	if conditions != nil {
		if conditions.WithoutOrder {
			stmt = stmt.Where("order_id IS NULL")
		}

		if conditions.InquiryID != 0 {
			stmt = stmt.Where("id = ?", conditions.InquiryID)
		}
	}

	stmt.LoadContext(ctx, &result)

	return result, nil
}

func (r *repo) CreateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error) {
	sess := r.NewSession(nil)

	result, err := sess.InsertInto(tables.inquires).
		Pair("description", inquiry.Description).
		Pair("first_name", inquiry.FirstName).
		Pair("last_name", inquiry.LastName).
		Pair("number", inquiry.Number).
		Pair("email", inquiry.Email).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	inquiry.ID = id

	return inquiry, nil
}

func (r *repo) UpdateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error) {
	sess := r.NewSession(nil)

	_, err := sess.Update(tables.inquires).
		Set("description", inquiry.Description).
		Set("first_name", inquiry.FirstName).
		Set("last_name", inquiry.LastName).
		Set("number", inquiry.Number).
		Set("email", inquiry.Email).
		Where("id = ?", inquiry.ID).
		ExecContext(ctx)

	if err != nil {
		return nil, err
	}

	return inquiry, nil
}
