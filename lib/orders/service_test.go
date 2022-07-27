package orders

import (
	"context"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/stretchr/testify/assert"
)

var (
	service, err = NewService(utility.NewEnv(utility.NewLogger()))

	inquiry = &lib.Inquiry{
		Description: "some test description",
		Email:       "test@test.com",
		Number:      "1112223333",
		FirstName:   "first_name",
		LastName:    "last_name",
	}

	order = &lib.Order{
		PaymentMethod: lib.PaymentMethodNotImplemented,
		Status:        lib.OrderStatusNotImplemented,
		Due:           "2022-07-26",
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

func TestOrderFunctionality(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	if order.InquiryID <= 0 {
		inquiry, err = service.SaveInquiry(ctx, inquiry)
		if err != nil {
			t.Error(err)
		}

		order.InquiryID = inquiry.ID
	}

	new, err := service.SaveOrder(ctx, order)

	if err != nil {
		t.Error(err)
	}

	//synthetically replace id to omit it during
	//asserting equality
	order.ID = new.ID

	assert.Equal(t, order, new)

	order, err := service.SaveOrder(ctx, new)

	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, new, order)
}
