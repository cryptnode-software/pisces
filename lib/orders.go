package lib

import "context"

//OrderService represents the OrderService interface
type OrderService interface {
	GetInquires(ctx context.Context, conditions *GetInquiryConditions) ([]*Inquiry, error)
	GetOrders(context.Context, *OrderConditions) ([]*Order, error)
	GetInquiry(ctx context.Context, id int64) (*Inquiry, error)
	SaveInquiry(context.Context, *Inquiry) (*Inquiry, error)
	GetOrder(ctx context.Context, id int64) (*Order, error)
	ArchiveOrder(context.Context, *Order) (*Order, error)
	SaveOrder(context.Context, *Order) (*Order, error)
}

//OrderConditions defines the different conditions that
//we can filter and sort orders by.
type OrderConditions struct {
	Status OrderStatus
	SortBy OrdersSortBy
}

//OrdersSortBy represents the primitive type for all the sorting capabilities
type OrdersSortBy string

const (
	//OrdersSortByDateAscending ...
	OrdersSortByDateAscending OrdersSortBy = "DATE_ASCENDING"

	//OrdersSortByDateDescending ...
	OrdersSortByDateDescending OrdersSortBy = "DATE_DESCENDING"

	//OrdersSortByDueDescending ...
	OrdersSortByDueDescending OrdersSortBy = "DUE_DESCENDING"

	//OrderSortByDueAscending ...
	OrderSortByDueAscending OrdersSortBy = "DUE_ASCENDING"
)

//Order the general structure of an order
type Order struct {
	PaymentMethod PaymentMethod `json:"payment_method"`
	InquiryID     int64         `json:"inquiry_id"`
	Status        OrderStatus   `json:"status"`
	ExtID         string        `json:"ext_id"`
	Total         float32       `json:"total"`
	Due           string        `json:"due"`
	ID            int64         `json:"id"`
}

//PaymentMethod the primitive data type for all of
//our payment methods.
type PaymentMethod string

const (
	//PaymentMethodNotImplemented this is our default payment method if
	//there are no other matching ones
	PaymentMethodNotImplemented PaymentMethod = "NOT_IMPLEMENTED"
	//PaymentMethodPaypal is the payment method that indicates that
	//the user is using paypal to checkout
	PaymentMethodPaypal PaymentMethod = "PAYPAL"
)

//Inquiry the structure contact info of a customer
type Inquiry struct {
	Description string   `json:"description"`
	Attachments []string `json:"attachments"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Number      string   `json:"number"`
	Email       string   `json:"email"`
	ID          int64    `json:"id"`
}

//OrderStatus this is the primitive datatype for OrderStatus' for
//how we handle it through the rest of the application
type OrderStatus string

const (
	//OrderStatusNotImplemented is our default when there is a mapping
	//or something else has gone wrong in the system.
	OrderStatusNotImplemented OrderStatus = "NOT_IMPLEMENTED"
	//OrderStatusAdminPending represents when the order is pending
	//on an admin to accept it. Once the/an admin has accepted the
	//order the order should go into a state of 'ACCEPTED'.
	OrderStatusAdminPending OrderStatus = "ADMIN_PENDING"
	//OrderStatusUserPending represents when the order is pending
	//on an user to accept/finalize their order. This is typically
	//the first step that the order should be in when initialized.
	OrderStatusUserPending OrderStatus = "USER_PENDING"
	//OrderStatusAccepted represents when both parties have consented
	//to the order and will fulfill it on either end. The consumer
	//using their goods to pay for a product that the business is
	//selling. This is typically the final step in the ordering
	//process
	OrderStatusAccepted OrderStatus = "ACCEPTED"
)

//GetInquiryConditions represents the different conditions that we
//can define when using the
type GetInquiryConditions struct {
	//InquiryID returns a single inquiry that matches the id
	InquiryID int64
	//WithoutOrder will return inquires that have been created
	//without any order associated with it.
	WithoutOrder bool
}
