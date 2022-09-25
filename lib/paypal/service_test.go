package paypal_test

import (
	"context"
	"testing"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/paypal"
	"github.com/google/uuid"
)

var (
	env = lib.NewEnv(commons.NewLogger(commons.EnvDev))

	service, err = paypal.NewService(env)

	id = uuid.New()

	order = &lib.Order{
		Total: 40.00,
		Model: commons.Model{
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
