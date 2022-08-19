package lib

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

func convertOrdersToProto(orders []*Order) (result []*proto.Order, err error) {
	result = make([]*proto.Order, len(orders))

	for i, order := range orders {
		o, err := convertOrderToProto(order)
		if err != nil {
			return nil, err
		}
		result[i] = o
	}

	return
}

func convertInquiryToProto(info *Inquiry) *proto.Inquiry {
	if info == nil {
		return nil
	}
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

func convertOrderToProto(order *Order) (result *proto.Order, err error) {
	result = new(proto.Order)

	result.PaymentMethod = convertPaymentMethodToProto(order.PaymentMethod)
	result.Status = convertOrderStatusToProto(order.Status)
	result.Inquiry = convertInquiryToProto(order.Inquiry)
	result.InquiryId = order.InquiryID.String()
	result.Id = order.ID.String()
	result.ExtId = order.ExtID
	result.Total = order.Total

	due, err := ptypes.TimestampProto(order.Due)

	if err != nil {
		return nil, err
	}

	result.Due = due

	return
}

func convertOrder(order *proto.Order) (result *Order, err error) {

	result = new(Order)

	due, err := ptypes.Timestamp(order.Due)
	if err != nil {
		return nil, err
	}

	result.Due = due

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

	if order.Inquiry != nil {
		result.Inquiry = convertInquiry(order.Inquiry)
	}

	return

}

func convertCart(cart []*proto.CartContents) (result []*Cart, err error) {
	result = make([]*Cart, len(cart))

	for i, pcontent := range cart {
		content := new(Cart)

		order, err := uuid.Parse(pcontent.OrderId)
		if err != nil {
			return nil, err
		}

		content.OrderID = order

		product, err := uuid.Parse(pcontent.ProductId)
		if err != nil {
			return nil, err
		}

		content.ProductID = product

		id, err := uuid.Parse(pcontent.Id)
		if err != nil {
			return nil, err
		}

		content.ID = id

		content.Quantity = pcontent.Quantity

		result[i] = content
	}

	return
}

func convertCartToProto(cart []*Cart) (result []*proto.CartContents) {
	result = make([]*proto.CartContents, len(cart))

	for i, content := range cart {
		pcontent := new(proto.CartContents)
		pcontent.ProductId = content.ProductID.String()
		pcontent.OrderId = content.OrderID.String()
		pcontent.Quantity = content.Quantity
		pcontent.Id = content.ID.String()

		result[i] = pcontent
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
