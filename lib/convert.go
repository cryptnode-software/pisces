package lib

import (
	"github.com/aws/smithy-go/time"
	"github.com/google/uuid"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

func convertOrdersToProto(orders []*Order) []*proto.Order {
	result := make([]*proto.Order, len(orders))
	for i, order := range orders {
		result[i] = convertOrderToProto(order)
	}
	return result
}

func convertInquiryToProto(info *Inquiry) *proto.Inquiry {
	return &proto.Inquiry{
		Body:        info.Description,
		Id:          info.ID.String(),
		FirstName:   info.FirstName,
		LastName:    info.LastName,
		PhoneNumber: info.Number,
		Email:       info.Email,
	}
}

func convertInquiry(info *proto.Inquiry) (inquiry *Inquiry) {
	inquiry = new(Inquiry)

	inquiry.FirstName = info.FirstName
	inquiry.Number = info.PhoneNumber
	inquiry.LastName = info.LastName
	inquiry.Description = info.Body
	inquiry.Email = info.Email

	if uuid, err := uuid.Parse(info.Id); err == nil {
		inquiry.ID = uuid
	}

	return
}

func convertOrderToProto(order *Order) *proto.Order {

	return &proto.Order{
		PaymentMethod: convertPaymentMethodToProto(order.PaymentMethod),
		Status:        convertOrderStatusToProto(order.Status),
		InquiryId:     order.InquiryID.String(),
		Due:           order.Due.String(),
		Id:            order.ID.String(),
		ExtId:         order.ExtID,
		Total:         order.Total,
	}
}

func convertOrder(order *proto.Order) (result *Order) {

	result = new(Order)

	if uuid, err := uuid.Parse(order.InquiryId); err == nil {
		result.InquiryID = uuid
	}

	if uuid, err := uuid.Parse(order.Id); err == nil {
		result.ID = uuid
	}

	result.PaymentMethod = convertPaymentMethod(order.PaymentMethod)
	result.Status = convertOrderStatus(order.Status)
	result.ExtID = order.ExtId
	result.Total = order.Total
	if due, err := time.ParseDateTime(order.Due); err == nil {
		result.Due = due
	}

	return

}

func convertCart(cart *proto.Cart) (result *Cart) {
	result = new(Cart)

	result.Contents = make([]*CartContents, len(cart.Contents))
	result.OrderID = OrderID(cart.OrderId)

	for i, product := range cart.Contents {
		result.Contents[i] = &CartContents{
			ProductID: product.ProductId,
			Quantity:  product.Quantity,
			ID:        product.Id,
		}
	}

	return
}

func convertCartToProto(cart *Cart) (result *proto.Cart) {
	result = new(proto.Cart)

	result.Contents = make([]*proto.CartContents, len(cart.Contents))
	result.OrderId = string(cart.OrderID)

	for i, product := range cart.Contents {

		result.Contents[i] = &proto.CartContents{
			ProductId: product.ProductID,
			Quantity:  product.Quantity,
			Id:        product.ID,
		}

	}

	return
}

func convertPaymentMethodToProto(method PaymentMethod) (result proto.PaymentMethod) {
	switch method {
	case PaymentMethodPaypal:
		result = proto.PaymentMethod_PaymentMethodPaypal
	default:
		result = proto.PaymentMethod_PaymentMethodNotImplemented
	}

	return
}
func convertPaymentMethod(method proto.PaymentMethod) PaymentMethod {

	result := PaymentMethodNotImplemented

	switch method {
	case proto.PaymentMethod_PaymentMethodPaypal:
		result = PaymentMethodPaypal
	}

	return result
}

func convertOrderStatus(status proto.OrderStatus) (result OrderStatus) {

	switch status {
	case proto.OrderStatus_AdminPending:
		result = OrderStatusAdminPending
	case proto.OrderStatus_UserPending:
		result = OrderStatusUserPending
	case proto.OrderStatus_Accepted:
		result = OrderStatusAccepted
	default:
		result = OrderStatusNotImplemented
	}

	return
}

func convertOrderStatusToProto(status OrderStatus) (result proto.OrderStatus) {

	switch status {
	case OrderStatusAccepted:
		result = proto.OrderStatus_Accepted
	case OrderStatusUserPending:
		result = proto.OrderStatus_UserPending
	case OrderStatusAdminPending:
		result = proto.OrderStatus_AdminPending
	default:
		result = proto.OrderStatus_NotImplemented
	}

	return
}

func convertUserFromProto(user *proto.User) *User {
	return &User{
		Username: user.Username,
		Email:    user.Email,
		Admin:    user.Admin,
	}
}

func convertInquiresToProto(inquires []*Inquiry) []*proto.Inquiry {
	result := make([]*proto.Inquiry, len(inquires))

	for i, inquiry := range inquires {
		result[i] = convertInquiryToProto(inquiry)
	}

	return result
}
