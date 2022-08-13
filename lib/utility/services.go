package utility

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/auth"
	"github.com/cryptnode-software/pisces/lib/cart"
	"github.com/cryptnode-software/pisces/lib/orders"
	"github.com/cryptnode-software/pisces/lib/paypal"
	"github.com/cryptnode-software/pisces/lib/product"
)

func Services(env *lib.Env) (services *lib.Services) {
	return &lib.Services{
		ProductService: productservice(env),
		PaypalService:  paypalservice(env),
		OrderService:   orderservice(env),
		CartService:    cartservice(env),
		AuthService:    authservice(env),
		S3Client:       s3client(env),
	}
}

//NewPaypalService returns a service that satisfies the clib.PaypalService interface
func paypalservice(env *lib.Env) lib.PaypalService {
	paypal, err := paypal.NewService(env)
	if err != nil {
		panic(err)
	}
	return paypal
}

//NewAuthService returns a service that satisfies the lib.AuthService interface
func authservice(env *lib.Env) lib.AuthService {
	service, err := auth.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

//NewOrderService returns a service that satisfies the lib.OrderService
func orderservice(env *lib.Env) lib.OrderService {
	order, err := orders.NewService(env)
	if err != nil {
		panic(err)
	}
	return order
}

//NewProductService returns a new product service
func productservice(env *lib.Env) lib.ProductService {
	service, err := product.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

//NewCartService returns a new cart service
func cartservice(env *lib.Env) lib.CartService {
	service, err := cart.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

func s3client(env *lib.Env) (client *s3.Client) {
	client = s3.NewFromConfig(aws.Config{
		Region: env.AWSEnv.Region,
		EndpointResolver: aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: fmt.Sprintf("https://%s.linodeobjects.com", env.AWSEnv.Region),
			}, nil
		}),
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (creds aws.Credentials, err error) {
			creds = aws.Credentials{
				AccessKeyID:     env.AWSEnv.AccessKey,
				SecretAccessKey: env.AWSEnv.SecretKey,
			}
			return
		}),
	},
	)
	return
}
