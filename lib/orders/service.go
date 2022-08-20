package orders

import (
	"context"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
			env.GormDB,
		},
	}, nil
}

//GetOrders returns orders sorted and filtered by the conditions provided
func (s *Service) GetOrders(ctx context.Context, conditions *lib.OrderConditions) ([]*lib.Order, error) {
	return s.repo.GetOrders(ctx, conditions)
}

//GetOrder returns a specific order and any information that is associated with
//it
func (s *Service) GetOrder(ctx context.Context, id uuid.UUID) (*lib.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

//GetInquiry returns a specific inquiry based on the id provided, if there is
//no inquiry found an exception will be raised.
func (s *Service) GetInquiry(ctx context.Context, id uuid.UUID) (*lib.Inquiry, error) {
	return s.repo.GetInquiry(ctx, id)
}

//SaveOrder will either create a new order (if required, i.e. the ID property is empty). Or
//will update a preexisting one.
func (s *Service) SaveOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {

	//inquiry should be required to create/update a order
	if order.Inquiry == nil && order.InquiryID == uuid.Nil {
		return nil, &errors.ErrNoOrderInquiryProvided{
			OrderID: order.ID.String(),
		}
	}

	//create new order
	if order.ID == uuid.Nil {
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
	if inquiry.ID == uuid.Nil {
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

func (s *Service) DeleteOrder(ctx context.Context, order *lib.Order, conditions *lib.DeleteConditions) error {
	if conditions != nil && conditions.HardDelete {
		return s.repo.HardDeleteOrder(ctx, order)
	}

	return s.repo.SoftDeleteOrder(ctx, order)
}

func (s *Service) DeleteInquiry(ctx context.Context, inquiry *lib.Inquiry, conditions *lib.DeleteConditions) error {
	if conditions != nil && conditions.HardDelete {
		return s.repo.HardDeleteInquiry(ctx, inquiry)
	}

	return s.repo.SoftDeleteInquiry(ctx, inquiry)
}

type repoi interface {
	GetInquires(ctx context.Context, conditions *lib.GetInquiryConditions) ([]*lib.Inquiry, error)
	UpdateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error)
	CreateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error)
	GetOrders(context.Context, *lib.OrderConditions) ([]*lib.Order, error)
	CreateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error)
	UpdateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error)
	GetInquiry(ctx context.Context, id uuid.UUID) (*lib.Inquiry, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*lib.Order, error)
	HardDeleteOrder(ctx context.Context, order *lib.Order) error
	SoftDeleteOrder(ctx context.Context, order *lib.Order) error
	HardDeleteInquiry(ctx context.Context, inquiry *lib.Inquiry) error
	SoftDeleteInquiry(ctx context.Context, inquiry *lib.Inquiry) error
}

type repo struct {
	*gorm.DB
}

func (r *repo) GetInquiry(ctx context.Context, id uuid.UUID) (inquiry *lib.Inquiry, err error) {
	inquiry = new(lib.Inquiry)
	r.DB.Model(new(lib.Inquiry)).First(inquiry, "id = ?", id)
	return
}

func (r *repo) CreateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {
	err := r.DB.Save(order).Error

	return order, err
}

func (r *repo) UpdateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {
	err := r.DB.Model(new(lib.Order)).
		Where("id = ?", order.ID).
		Update("payment_method", order.PaymentMethod).
		Update("status", order.Status).
		Update("ext_id", order.ExtID).
		Update("due", order.Due).
		Error

	return order, err
}

func (r *repo) GetOrders(ctx context.Context, conditions *lib.OrderConditions) ([]*lib.Order, error) {

	var result []*lib.Order

	tx := r.DB

	if conditions != nil {
		if conditions.Status != lib.OrderStatusNotImplemented {
			if err := tx.Model(new(lib.Order)).
				Preload("Inquiry").
				Where("status = ?", conditions.Status).
				Find(&result).Error; err != nil {
				return nil, err
			}
		}
	}

	if conditions == nil {
		if err := r.DB.Preload("Inquiry").Find(&result).Error; err != nil {
			return nil, err
		}
	}

	if tx.Error != nil {
		return nil, tx.Error
	}

	return result, nil
}

func (r *repo) GetOrder(ctx context.Context, id uuid.UUID) (order *lib.Order, err error) {
	order = new(lib.Order)
	r.DB.Preload("Inquiry").Model(new(lib.Order)).First(order, "id = ?", id)

	return
}

func (r *repo) GetInquires(ctx context.Context, conditions *lib.GetInquiryConditions) ([]*lib.Inquiry, error) {

	var result []*lib.Inquiry

	tx := r.DB.Model(new(lib.Inquiry))

	if conditions != nil {
		if conditions.WithoutOrder {
			tx = tx.Where("order_id IS NULL")
		}

		if conditions.InquiryID != 0 {
			tx = tx.Where("id = ?", conditions.InquiryID)
		}
	}
	tx.Find(&result)

	return result, nil
}

func (r *repo) CreateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error) {

	err := r.DB.Save(inquiry).Error

	return inquiry, err
}

func (r *repo) UpdateInquiry(ctx context.Context, inquiry *lib.Inquiry) (*lib.Inquiry, error) {

	err := r.DB.Save(inquiry).Error

	return inquiry, err
}

func (r *repo) HardDeleteOrder(ctx context.Context, order *lib.Order) error {
	return r.DB.Unscoped().Delete(order).Error
}

func (r *repo) SoftDeleteOrder(ctx context.Context, order *lib.Order) error {
	return r.DB.Delete(order).Error
}

func (r *repo) HardDeleteInquiry(ctx context.Context, inquiry *lib.Inquiry) error {
	return r.DB.Unscoped().Delete(inquiry).Error
}

func (r *repo) SoftDeleteInquiry(ctx context.Context, inquiry *lib.Inquiry) error {
	return r.DB.Delete(inquiry).Error
}
