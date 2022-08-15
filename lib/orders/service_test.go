package orders_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/orders"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/stretchr/testify/assert"
)

var (
	service, err = orders.NewService(utility.NewEnv(utility.NewLogger()))

	inquiry = &lib.Inquiry{
		Description: "some test description",
		Email:       "test@test.com",
		Number:      "111-222-3333",
		FirstName:   "first_name",
		LastName:    "last_name",
	}

	order = &lib.Order{
		PaymentMethod: lib.PaymentMethodNotImplemented,
		Status:        lib.OrderStatusNotImplemented,
		Due:           time.Now().Add(60 * 24),
	}

	ctx = context.Background()
)

func TestInquiryFunctionality(t *testing.T) {
	new, err := service.SaveInquiry(ctx, inquiry)

	if err != nil {
		t.Error(err)
	}

	inquiry.ID = new.ID

	assert.Equal(t, inquiry, new)
}

func TestSaveInquiry(t *testing.T) {
	tables := []struct {
		inquiry  lib.Inquiry
		expected lib.Inquiry
	}{
		{
			expected: lib.Inquiry{
				Description: "Magna ipsum culpa labore pariatur elit commodo consequat esse est.",
				Email:       "test.user@test.io",
				Number:      "000-000-0000",
				FirstName:   "test",
				LastName:    "user",
			},
			inquiry: lib.Inquiry{
				Description: "Magna ipsum culpa labore pariatur elit commodo consequat esse est.",
				Email:       "test.user@test.io",
				Number:      "000-000-0000",
				FirstName:   "test",
				LastName:    "user",
			},
		},
	}

	for _, table := range tables {
		inquiry, err := service.SaveInquiry(ctx, &table.inquiry)
		if err != nil {
			t.Error(err)
			continue
		}

		table.expected.Model = inquiry.Model

		assert.Equal(t, table.expected, *inquiry)
	}
}

func TestSaveOrder(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		expected *lib.Order
		order    lib.Order
	}{
		{
			expected: &lib.Order{
				PaymentMethod: lib.PaymentMethodNotImplemented,
				Status:        lib.OrderStatusNotImplemented,
				Inquiry: &lib.Inquiry{
					Description: "Magna ipsum culpa labore pariatur elit commodo consequat esse est.",
					Email:       "test.user@test.io",
					Number:      "000-000-0000",
					FirstName:   "test",
					LastName:    "user",
				},
			},
			order: lib.Order{
				PaymentMethod: lib.PaymentMethodNotImplemented,
				Status:        lib.OrderStatusNotImplemented,
				Due:           time.Now().Add(60 * 24),
				Inquiry: &lib.Inquiry{
					Description: "Magna ipsum culpa labore pariatur elit commodo consequat esse est.",
					Email:       "test.user@test.io",
					Number:      "000-000-0000",
					FirstName:   "test",
					LastName:    "user",
				},
			},
		},
	}

	for _, table := range tables {
		order, err := service.SaveOrder(ctx, &table.order)

		if err != nil {
			t.Error(err)
			return
		}

		table.expected.Inquiry.Model = order.Inquiry.Model
		table.expected.InquiryID = order.InquiryID
		table.expected.Model = order.Model
		table.expected.Due = order.Due

		assert.Equal(t, table.expected, order)
	}
}

func TestFailSaveOrder(t *testing.T) {
	tables := []struct {
		order lib.Order
	}{
		{
			order: lib.Order{
				PaymentMethod: lib.PaymentMethodNotImplemented,
				Status:        lib.OrderStatusNotImplemented,
				Due:           time.Now().Add(60 * 24),
			},
		},
	}

	for _, table := range tables {
		_, err := service.SaveOrder(ctx, &table.order)
		if err == nil {
			t.Error(errors.New("order successfully saved when it was suppose to fail"))
		}
	}
}

func TestGetOrders(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		order lib.Order
	}{
		{
			order: lib.Order{},
		},
	}

	for _, table := range tables {
		expected := table.order

		order, err := service.SaveOrder(ctx, &table.order)

		if err != nil {
			t.Error(err)
			continue
		}

		expected.Model = order.Model

		order, err = service.GetOrder(ctx, order.ID)
	}
}
