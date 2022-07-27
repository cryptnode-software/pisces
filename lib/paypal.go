package lib

import "context"

//PaypalService represents how the paypal service should be structured
type PaypalService interface {
	GenerateClientToken(context.Context) (*GenerateClientTokenResponse, error)
	CreateOrder(context.Context, *Order) (*Order, error)
}

//GenerateClientTokenResponse ...
type GenerateClientTokenResponse struct {
	Token string `json:"client_token"`
}
