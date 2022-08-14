package paypal

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/plutov/paypal"
)

//Service represents the paypal structure and all of the calls that we currently
//support
type Service struct {
	client *paypal.Client
	env    *lib.Env
}

//NewService returns a new paypal service that we can use in
//order to interact with paypal directly
func NewService(env *lib.Env) (*Service, error) {

	client, err := paypal.NewClient(env.PaypalEnv.ClientID, env.PaypalEnv.SecretID, env.PaypalEnv.Host)
	if err != nil {
		panic(err)
	}

	_, err = client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	return &Service{
		client,
		env,
	}, nil

}

//CreateOrder creates a paypal order
func (service *Service) CreateOrder(ctx context.Context, order *lib.Order) (*lib.Order, error) {

	if order.ID == nil {
		return nil, errors.New("no local order id associated with paypal order")
	}

	porder, err := service.client.CreateOrder(
		paypal.OrderIntentAuthorize, []paypal.PurchaseUnitRequest{
			{
				ReferenceID: string(*order.ID),
				Amount: &paypal.PurchaseUnitAmount{
					Currency: "USD",
					Value:    fmt.Sprintf("%.2f", order.Total),
				},
			},
		},
		nil, nil)

	if err != nil {
		return nil, err
	}

	order.ExtID = porder.ID

	return order, nil
}

//GenerateClientToken generates a client token for frontend rendering
func (service *Service) GenerateClientToken(ctx context.Context) (*lib.GenerateClientTokenResponse, error) {
	buf := bytes.NewBuffer([]byte{})
	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s%s", service.client.APIBase, "/v1/identity/generate-token"), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Language", "en_US")

	//get access token before creating client token
	_, err = service.client.GetAccessToken()
	if err != nil {
		return nil, err
	}

	res := &lib.GenerateClientTokenResponse{}

	err = service.client.SendWithAuth(req, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
