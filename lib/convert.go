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
		Attachments: convertAttachmentsToProto(info.Attachments),
		Body:        info.Description,
		Id:          info.ID.String(),
		FirstName:   info.FirstName,
		LastName:    info.LastName,
		PhoneNumber: info.Number,
		Email:       info.Email,
	}
}

func convertInquiry(info *proto.Inquiry) (inquiry *Inquiry) {
	if info == nil {
		return nil
	}

	inquiry = new(Inquiry)

	inquiry.Attachments = convertAttachments(info.Attachments)
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

func convertAttachmentsToProto(attachments []*Attachment) (result []*proto.Attachment) {
	result = make([]*proto.Attachment, len(attachments))

	for i, attachment := range attachments {
		result[i] = convertAttachmentToProto(attachment)
	}

	return
}

func convertAttachments(attachments []*proto.Attachment) (result []*Attachment) {
	result = make([]*Attachment, len(attachments))

	for i, attachment := range attachments {
		result[i] = convertAttachment(attachment)
	}

	return
}

func convertAttachmentToProto(attachment *Attachment) (result *proto.Attachment) {
	if attachment == nil {
		return nil
	}

	result = new(proto.Attachment)

	result.Type = convertAttachmentTypeToProto(attachment.Type)
	result.Url = attachment.URL

	return
}

func convertAttachment(attachment *proto.Attachment) (result *Attachment) {
	if attachment == nil {
		return nil
	}

	result = new(Attachment)

	result.Type = convertAttachmentType(attachment.Type)
	result.URL = attachment.Url

	return
}

func convertAttachmentType(atype proto.AttachmentType) (result AttachmentType) {
	result = AttachmentTypeNotImplemented

	switch atype {
	case proto.AttachmentType_AttachmentTypeFile:
		result = AttachmentTypeFile
	case proto.AttachmentType_AttachmentTypeImage:
		result = AttachmentTypeImage
	}

	return
}

func convertAttachmentTypeToProto(atype AttachmentType) (result proto.AttachmentType) {
	result = proto.AttachmentType_AttachmentTypeNotImplemented

	switch atype {
	case AttachmentTypeFile:
		result = proto.AttachmentType_AttachmentTypeFile
	case AttachmentTypeImage:
		result = proto.AttachmentType_AttachmentTypeImage
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

	result.Cart = convertCartToProto(order.Cart)

	return
}

func convertOrder(order *proto.Order) (result *Order, err error) {

	result = new(Order)

	result.PaymentMethod = convertPaymentMethod(order.PaymentMethod)

	if result.Due, err = ptypes.Timestamp(order.Due); err != nil {
		return nil, err
	}

	if uuid, err := uuid.Parse(order.InquiryId); err == nil {
		result.InquiryID = uuid
	}

	if uuid, err := uuid.Parse(order.Id); err == nil {
		result.ID = uuid
	}

	result.Status = convertOrderStatus(order.Status)
	result.ExtID = order.ExtId
	result.Total = order.Total

	if order.Inquiry != nil {
		result.Inquiry = convertInquiry(order.Inquiry)
	}

	result.Cart = convertCart(order.Cart)

	return

}

func convertCart(cart []*proto.CartContents) (result []*Cart) {
	if cart == nil {
		return nil
	}
	result = make([]*Cart, len(cart))

	for i, pcontent := range cart {

		content := new(Cart)

		if product, err := uuid.Parse(pcontent.ProductId); err == nil {
			content.ProductID = product
		}

		if order, err := uuid.Parse(pcontent.OrderId); err == nil {
			content.OrderID = order
		}

		if id, err := uuid.Parse(pcontent.Id); err == nil {
			content.ID = id
		}

		content.Quantity = pcontent.Quantity

		result[i] = content

	}

	return
}

func convertCartToProto(cart []*Cart) (result []*proto.CartContents) {
	if cart == nil {
		return nil
	}

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
func convertPaymentMethod(method proto.PaymentMethod) (result PaymentMethod) {

	switch method {
	case proto.PaymentMethod_PaymentMethodPaypal:
		result = PaymentMethodPaypal
	default:
		result = PaymentMethodNotImplemented
	}

	return
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

func convertInquiresToProto(inquires []*Inquiry) []*proto.Inquiry {
	result := make([]*proto.Inquiry, len(inquires))

	for i, inquiry := range inquires {
		result[i] = convertInquiryToProto(inquiry)
	}

	return result
}

func convertProductsToProto(products []*Product) (result []*proto.Product) {
	result = make([]*proto.Product, len(products))

	for i, product := range products {
		result[i] = convertProductToProto(product)
	}

	return
}

func convertProductToProto(product *Product) (result *proto.Product) {
	result = new(proto.Product)

	result.Inventory = int64(product.Inventory)
	result.Description = product.Description
	result.Id = product.ID.String()
	result.Cost = product.Cost
	result.Name = product.Name

	return
}

func convertProductsFromProto(products []*proto.Product) (result []*Product) {
	result = make([]*Product, len(products))

	for i, product := range products {
		result[i] = convertProductFromProto(product)
	}

	return
}

func convertProductFromProto(product *proto.Product) (result *Product) {
	result = new(Product)

	if id, err := uuid.Parse(product.Id); err == nil {
		result.ID = id
	}

	result.Inventory = int(product.Inventory)
	result.Description = product.Description
	result.Cost = product.Cost
	result.Name = product.Name

	return
}
