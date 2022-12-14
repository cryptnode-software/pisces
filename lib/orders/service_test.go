package orders_test

import (
	"context"
	"errors"
	"testing"
	"time"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/orders"
	"github.com/stretchr/testify/assert"
)

var (
	service, err = orders.NewService(lib.NewEnv(commons.NewLogger(commons.EnvDev)))

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

		if err := deseed([]*lib.Inquiry{
			inquiry,
		}); err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestGetInquiry(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		inquiry *lib.Inquiry
	}{
		{
			inquiry: &lib.Inquiry{
				Description: "Magna ipsum culpa labore pariatur elit commodo consequat esse est.",
				Email:       "test.user@test.io",
				Number:      "000-000-0000",
				FirstName:   "test",
				LastName:    "user",
			},
		},
	}

	for _, table := range tables {
		if err := seed([]*lib.Inquiry{
			table.inquiry,
		}); err != nil {
			t.Error(err)
			continue
		}

		inquiry, err := service.GetInquiry(ctx, table.inquiry.ID)
		if err != nil {
			t.Error(err)
			return
		}

		table.inquiry.CreatedAt = inquiry.CreatedAt
		table.inquiry.UpdatedAt = inquiry.UpdatedAt

		assert.Equal(t, table.inquiry, inquiry)

		if err := deseed([]*lib.Inquiry{table.inquiry}); err != nil {
			t.Error(err)
			continue
		}
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
		order, err := service.SaveOrder(ctx, &table.order, nil)

		if err != nil {
			t.Error(err)
			return
		}

		table.expected.Inquiry.Model = order.Inquiry.Model
		table.expected.InquiryID = order.InquiryID
		table.expected.Model = order.Model
		table.expected.Due = order.Due

		assert.Equal(t, table.expected, order)

		if err := deseed([]*lib.Order{
			order,
		}); err != nil {
			t.Error(err)
			return
		}

		if err := deseed([]*lib.Inquiry{
			order.Inquiry,
		}); err != nil {
			t.Error(err)
			continue
		}
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
		_, err := service.SaveOrder(ctx, &table.order, nil)
		if err == nil {
			t.Error(errors.New("order successfully saved when it was suppose to fail"))
		}
	}
}

func TestGetOrder(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		order lib.Order
	}{
		{
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
				Cart: make([]*lib.Cart, 0),
			},
		},
	}

	for _, table := range tables {
		expected := &table.order

		if err := seed([]*lib.Order{
			expected,
		}); err != nil {
			t.Error(err)
			continue
		}

		order, err = service.GetOrder(ctx, expected.ID)

		if err != nil {
			t.Error(err)
			continue
		}

		expected.Inquiry.Model = order.Inquiry.Model
		expected.InquiryID = order.InquiryID

		expected.Model = order.Model

		expected.Due = order.Due

		assert.Equal(t, expected, order)

		deseed([]*lib.Order{
			expected,
		})

		if err := deseed([]*lib.Inquiry{
			order.Inquiry,
		}); err != nil {
			t.Error(err)
			continue
		}
	}
}

func TestGetOrders(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	tables := []struct {
		expected []*lib.Order
		status   lib.OrderStatus
	}{
		{
			expected: []*lib.Order{
				{
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
			status: lib.OrderStatusNotImplemented,
		},
		{
			expected: []*lib.Order{
				{
					PaymentMethod: lib.PaymentMethodNotImplemented,
					Status:        lib.OrderStatusAccepted,
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
			status: lib.OrderStatusAccepted,
		},
	}

	for _, table := range tables {
		if err := seed(table.expected); err != nil {
			t.Error(err)
			return
		}

		orders, err := service.GetOrders(ctx, &lib.OrderConditions{
			Status: table.status,
		})
		if err != nil {
			t.Error(err)
			return
		}

		for i, o := range orders {
			orders[i].Due = o.Due
			assert.Equal(t, orders[i], o)
		}

		if err := deseed(table.expected); err != nil {
			t.Error(err)
			return
		}

		for _, order := range table.expected {
			if err := deseed([]*lib.Inquiry{
				order.Inquiry,
			}); err != nil {
				t.Error(err)
				continue
			}
		}

	}

}

func seed[T *lib.Order | *lib.Inquiry](models []T) error {
	for _, model := range models {
		switch model := any(model).(type) {
		case *lib.Inquiry:
			inquiry, err := service.SaveInquiry(ctx, model)
			model = inquiry
			if err != nil {
				return err
			}
		case *lib.Order:
			inquiry, err := service.SaveInquiry(ctx, model.Inquiry)
			if err != nil {
				return err
			}
			model.Inquiry = inquiry

			order, err := service.SaveOrder(ctx, model, nil)
			model = order

			if err != nil {
				return err
			}
		}
	}
	return nil
}

func deseed[T *lib.Order | *lib.Inquiry](models []T) error {
	for _, model := range models {
		switch model := any(model).(type) {
		case *lib.Order:
			err := service.DeleteOrder(ctx, model, &lib.DeleteConditions{
				HardDelete: true,
			})

			if err != nil {
				return err
			}

			if model.Inquiry != nil {
				err = service.DeleteInquiry(ctx, model.Inquiry, &lib.DeleteConditions{
					HardDelete: true,
				})

				if err != nil {
					return err
				}
			}

		case *lib.Inquiry:
			err := service.DeleteInquiry(ctx, model, &lib.DeleteConditions{
				HardDelete: true,
			})

			if err != nil {
				return err
			}
		}
	}
	return nil
}
