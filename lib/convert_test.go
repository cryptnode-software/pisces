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

}

func TestConvertOrderToProto(t *testing.T) {

}
func TestConvertInquiryToProto(t *testing.T) {

}
