package paypal_test

import (
	"context"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/paypal"
	"github.com/google/uuid"
)

var (
	env = lib.NewEnv(lib.NewLogger(lib.EnvDev))

	service, err = paypal.NewService(env)

	id = uuid.New()

	order = &lib.Order{
		Total: 40.00,
		Model: lib.Model{
			ID: id,
		},
	}

	ctx = context.Background()
)

func TestGenerateClientSideToken(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}
	response, err := service.GenerateClientToken(ctx)
	if err != nil {
		t.Error(err)
		return
	}

	if response.Token == "" {
		t.Error("no client token was returned from response")
		return
	}
}

func TestCreateOrder(t *testing.T) {
	if err != nil {
		t.Error(err)
		return
	}

	_, err := service.CreateOrder(ctx, order)
	if err != nil {
		t.Error(err)
		return
	}
}
