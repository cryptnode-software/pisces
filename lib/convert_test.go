package lib

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
)

var (
	protodue, err = ptypes.TimestampProto(due)
	id            = uuid.New()
	due           = time.Now()
)

func TestConvertOrdersToProto(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		expected []*proto.Order
		orders   []*Order
	}{
		{
			expected: []*proto.Order{
				{
					PaymentMethod: proto.PaymentMethod_PaymentMethodPaypal,
					Status:        proto.OrderStatus_UserPending,
					Inquiry: &proto.Inquiry{
						Body:        "Quis incididunt aliqua ex duis proident sunt sit.",
						Email:       "test.user@test.io",
						PhoneNumber: "000-000-0000",
						Attachments: []*proto.Attachment{},
						Id:          id.String(),
						LastName:    "user",
						FirstName:   "test",
					},
					InquiryId: uuid.Nil.String(),
					Id:        id.String(),
					Due:       protodue,
					Total:     40,
				},
			},
			orders: []*Order{
				{
					PaymentMethod: PaymentMethodPaypal,
					Inquiry: &Inquiry{
						Description: "Quis incididunt aliqua ex duis proident sunt sit.",
						Email:       "test.user@test.io",
						Number:      "000-000-0000",
						LastName:    "user",
						FirstName:   "test",
						Attachments: []*Attachment{},
						Model: Model{
							ID: id,
						},
					},
					Status: OrderStatusUserPending,
					Model: Model{
						ID: id,
					},
					Due:   due,
					Total: 40.00,
				},
			},
		},
	}

	for _, table := range tables {
		orders, err := convertOrdersToProto(table.orders)

		if err != nil {
			t.Error(err)
			continue
		}

		assert.Equal(t, table.expected, orders)
	}
}

func TestConvertOrderToProto(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		order    *Order
		expected *proto.Order
	}{
		{
			order: &Order{
				PaymentMethod: PaymentMethodPaypal,
				Inquiry: &Inquiry{
					Description: "Quis incididunt aliqua ex duis proident sunt sit.",
					Email:       "test.user@test.io",
					Attachments: []*Attachment{},
					Number:      "000-000-0000",
					LastName:    "user",
					FirstName:   "test",
					Model: Model{
						ID: id,
					},
				},
				Status: OrderStatusUserPending,
				Model: Model{
					ID: id,
				},
				Due:   due,
				Total: 40.00,
			},
			expected: &proto.Order{
				PaymentMethod: proto.PaymentMethod_PaymentMethodPaypal,
				Status:        proto.OrderStatus_UserPending,
				Inquiry: &proto.Inquiry{
					Body:        "Quis incididunt aliqua ex duis proident sunt sit.",
					Attachments: []*proto.Attachment{},
					Email:       "test.user@test.io",
					PhoneNumber: "000-000-0000",
					Id:          id.String(),
					LastName:    "user",
					FirstName:   "test",
				},
				InquiryId: uuid.Nil.String(),
				Id:        id.String(),
				Due:       protodue,
				Total:     40,
			},
		},
	}

	for _, table := range tables {
		order, err := convertOrderToProto(table.order)

		if err != nil {
			t.Error(err)
		}

		table.expected.Due = order.Due

		assert.Equal(t, table.expected, order)
	}
}

func TestConvertOrder(t *testing.T) {

	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		order    *proto.Order
		expected *Order
	}{
		{
			expected: &Order{
				PaymentMethod: PaymentMethodPaypal,
				Inquiry: &Inquiry{
					Description: "Quis incididunt aliqua ex duis proident sunt sit.",
					Email:       "test.user@test.io",
					Attachments: []*Attachment{},
					Number:      "000-000-0000",
					LastName:    "user",
					FirstName:   "test",
					Model: Model{
						ID: id,
					},
				},
				Status: OrderStatusUserPending,
				Model: Model{
					ID: id,
				},
				Due:   due,
				Total: 40.00,
			},
			order: &proto.Order{
				PaymentMethod: proto.PaymentMethod_PaymentMethodPaypal,
				Status:        proto.OrderStatus_UserPending,
				Inquiry: &proto.Inquiry{
					Body:        "Quis incididunt aliqua ex duis proident sunt sit.",
					Email:       "test.user@test.io",
					Attachments: []*proto.Attachment{},
					PhoneNumber: "000-000-0000",
					Id:          id.String(),
					LastName:    "user",
					FirstName:   "test",
				},
				InquiryId: uuid.Nil.String(),
				Id:        id.String(),
				Due:       protodue,
				Total:     40,
			},
		},
	}

	for _, table := range tables {
		order, err := convertOrder(table.order)

		if err != nil {
			t.Error(err)
		}

		table.expected.Due = order.Due

		assert.Equal(t, table.expected, order)
	}

}

func TestConvertInquiry(t *testing.T) {

	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		expected *Inquiry
		inquiry  *proto.Inquiry
	}{
		{
			expected: &Inquiry{
				Description: "Quis incididunt aliqua ex duis proident sunt sit.",
				Email:       "test.user@test.io",
				Number:      "000-000-0000",
				Attachments: []*Attachment{},
				LastName:    "user",
				FirstName:   "test",
				Model: Model{
					ID: id,
				},
			},
			inquiry: &proto.Inquiry{
				Body:        "Quis incididunt aliqua ex duis proident sunt sit.",
				Email:       "test.user@test.io",
				Attachments: []*proto.Attachment{},
				PhoneNumber: "000-000-0000",
				Id:          id.String(),
				LastName:    "user",
				FirstName:   "test",
			},
		},
	}

	for _, table := range tables {
		inquiry := convertInquiry(table.inquiry)

		assert.Equal(t, table.expected, inquiry)
	}
}

func TestConvertInquiryToProto(t *testing.T) {

	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		inquiry  *Inquiry
		expected *proto.Inquiry
	}{
		{
			inquiry: &Inquiry{
				Description: "Quis incididunt aliqua ex duis proident sunt sit.",
				Email:       "test.user@test.io",
				Attachments: []*Attachment{},
				Number:      "000-000-0000",
				LastName:    "user",
				FirstName:   "test",
				Model: Model{
					ID: id,
				},
			},
			expected: &proto.Inquiry{
				Body:        "Quis incididunt aliqua ex duis proident sunt sit.",
				Email:       "test.user@test.io",
				Attachments: []*proto.Attachment{},
				PhoneNumber: "000-000-0000",
				Id:          id.String(),
				LastName:    "user",
				FirstName:   "test",
			},
		},
	}

	for _, table := range tables {
		inquiry := convertInquiryToProto(table.inquiry)

		assert.Equal(t, table.expected, inquiry)
	}
}

func TestConvertAttachment(t *testing.T) {
	if err != nil {
		t.Error(err)
	}

	url := "https://github.com/cryptnode-software/pisces"

	tables := []struct {
		attachment *proto.Attachment
		expected   *Attachment
	}{
		{
			attachment: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeFile,
				Url:  url,
			},
			expected: &Attachment{
				Type: AttachmentTypeFile,
				URL:  url,
			},
		},
		{
			attachment: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeImage,
				Url:  url,
			},
			expected: &Attachment{
				Type: AttachmentTypeImage,
				URL:  url,
			},
		},
		{
			attachment: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeNotImplemented,
				Url:  url,
			},
			expected: &Attachment{
				Type: AttachmentTypeNotImplemented,
				URL:  url,
			},
		},
	}

	for _, table := range tables {
		attachment := convertAttachment(table.attachment)

		assert.Equal(t, table.expected, attachment)
	}
}

func TestConvertAttachmentToProto(t *testing.T) {
	if err != nil {
		t.Error(err)
	}

	url := "https://github.com/cryptnode-software/pisces"

	tables := []struct {
		expected   *proto.Attachment
		attachment *Attachment
	}{
		{
			expected: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeFile,
				Url:  url,
			},
			attachment: &Attachment{
				Type: AttachmentTypeFile,
				URL:  url,
			},
		},
		{
			expected: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeImage,
				Url:  url,
			},
			attachment: &Attachment{
				Type: AttachmentTypeImage,
				URL:  url,
			},
		},
		{
			expected: &proto.Attachment{
				Type: proto.AttachmentType_AttachmentTypeNotImplemented,
				Url:  url,
			},
			attachment: &Attachment{
				Type: AttachmentTypeNotImplemented,
				URL:  url,
			},
		},
	}

	for _, table := range tables {
		attachment := convertAttachmentToProto(table.attachment)

		assert.Equal(t, table.expected, attachment)
	}
}
