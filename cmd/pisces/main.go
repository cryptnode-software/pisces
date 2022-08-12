package main

import (
	"context"
	"net/http"

	// "github.com/aws/aws-sdk-go-v3/aws/session"
	// "github.com/aws/aws-sdk-go-v3/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	clib "github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/auth"
	"github.com/cryptnode-software/pisces/lib/cart"
	order "github.com/cryptnode-software/pisces/lib/orders"
	"github.com/cryptnode-software/pisces/lib/paypal"
	"github.com/cryptnode-software/pisces/lib/product"
	upload "github.com/cryptnode-software/pisces/lib/upload"
	"github.com/gocraft/dbr/v2"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	paylib "github.com/plutov/paypal"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"

	"flag"
	"fmt"
	"log"
	"os"
)

var (
	// Version holds app version
	Version string
	// Build holds the build datetime
	Build string
	// GitHash holds the current git hash
	GitHash string
)

const (
	envDatabaseURL string = "DB_CONNECTION"
	env            string = "ENV"

	envPaypalClientID string = "PAYPAL_CLIENT_ID"
	envPaypalSecretID string = "PAYPAL_SECRET_ID"

	envJWTSecret string = "JWT_SECRET"

	envS3Bucket    string = "S3_BUCKET"
	envS3AccessKey string = "AWS_ACCESS_KEY_ID"
	envS3SecretKey string = "AWS_SECRET_ACCESS_KEY"
	envS3Region    string = "AWS_REGION"
)

var (
	exempt = map[string]bool{
		"/pisces.Pisces/GeneratePaypalClientToken": true,
		"/pisces.Pisces/AuthorizeOrder":            true,
		"/pisces.Pisces/GetTotalCost":              true,
		"/pisces.Pisces/SaveInquiry":               true,
		"/pisces.Pisces/GetInquires":               true,
		"/pisces.Pisces/SaveOrder":                 true,
		"/pisces.Pisces/GetOrders":                 true,
		"/pisces.Pisces/SaveCart":                  true,
		"/pisces.Pisces/Login":                     true,
		"/pisces.Pisces/CreateUser":                true,
	}

	admin = map[string]bool{
		"/pisces.Pisces/GetInquires": true,
	}
)

func main() {

	port := flag.Int("port", 4081, "grpc port")

	flag.Parse()
	environment := NewEnv(NewLogger())

	services := &clib.Services{
		ProductService: NewProductService(environment),
		PaypalService:  NewPaypalService(environment),
		OrderService:   NewOrderService(environment),
		AuthService:    NewAuthService(environment),
		CartService:    NewCartService(environment),
	}

	gw, err := clib.NewGateway(environment, services)
	if err != nil {
		panic(err)
	}

	logger := environment.Log
	logger.Info("starting container...")

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_recovery.UnaryServerInterceptor(
					grpc_recovery.WithRecoveryHandlerContext(
						func(ctx context.Context, p interface{}) error {
							logger.Error("grpc_recovery", p, ctx)

							return p.(error)
						},
					),
				),
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
					logger.Info(info.FullMethod)
					if exempt[info.FullMethod] {
						return handler(ctx, req)
					}

					if admin[info.FullMethod] {
						if _, err := gw.AuthenticateAdmin(ctx); err != nil {
							return nil, err
						}

						return handler(ctx, req)
					}

					if _, err := gw.AuthenticateToken(ctx); err != nil {
						return nil, err
					}
					return handler(ctx, req)
				},
			),
		),
	}

	grpcServer := grpc.NewServer(opts...)
	proto.RegisterPiscesServer(grpcServer, gw)

	server := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(str string) bool {
			return true // change this
		}),
	)

	handler := func(resp http.ResponseWriter, req *http.Request) {
		server.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: http.HandlerFunc(handler),
	}

	logger.Info("listening on port :4081")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
}

//NewLogger returns a new logger based off the current environment
func NewLogger() clib.Logger {
	environ := os.Getenv(env)
	if environ == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
		return nil
	}

	return clib.NewZapper(clib.Environment(environ))
}

//NewEnv returns a new environment pre-populated with the provided logger
func NewEnv(logger clib.Logger) *clib.Env {
	environ := os.Getenv(env)
	if environ == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
		return nil
	}

	result := &clib.Env{
		Environment: clib.Environment(environ),
		Log:         logger,
	}

	//paypal config
	{
		client := os.Getenv(envPaypalClientID)
		if client == "" {
			log.Fatalf("paypal client id not provided please provide %s env variable", envPaypalClientID)
		}

		secret := os.Getenv(envPaypalSecretID)
		if secret == "" {
			log.Fatalf("paypal secret id not provided please provide %s env variable", envPaypalSecretID)
		}

		result.PaypalEnv = &clib.PaypalEnv{
			ClientID: client,
			SecretID: secret,
		}

		switch result.Environment {
		case clib.EnvProd:
			result.PaypalEnv.Host = paylib.APIBaseLive
		default:
			result.PaypalEnv.Host = paylib.APIBaseSandBox
		}
	}

	//auth config
	{
		jwtSecret := os.Getenv(envJWTSecret)
		if jwtSecret == "" {
			log.Fatalf("%s not set, if not properly set jwt tokens will be unsafe to use", envJWTSecret)
		}

		result.JWTEnv = &clib.JWTEnv{
			Secret: jwtSecret,
		}
	}

	//database config
	{
		sql, err := dbr.Open("mysql", os.Getenv(envDatabaseURL), nil)
		if err != nil {
			log.Fatal(err)
		}
		result.DB = sql
	}

	//aws config
	{
		config := new(clib.AWSEnv)
		if config.Region = os.Getenv(envS3Region); config.Region == "" {
			log.Fatalf("%s not, and required for s3 configuration", envS3Region)
		}
		if config.AccessKey = os.Getenv(envS3AccessKey); config.AccessKey == "" {
			log.Fatalf("%s not set and required for s3 configuration", envS3AccessKey)
		}
		if config.Bucket = os.Getenv(envS3Bucket); config.Bucket == "" {
			log.Fatalf("%s not set and required for s3 configuration", envS3Bucket)
		}
		if config.SecretKey = os.Getenv(envS3SecretKey); config.SecretKey == "" {
			log.Fatalf("%s not set and required for s3 configuration", envS3SecretKey)
		}
		result.AWSEnv = config
	}

	//upload config
	{
		result.Upload = new(clib.UploadEnv)

		switch result.Environment {
		case clib.EnvDev:
			result.Upload.Type = clib.UploadTypeLinode
			result.Upload.Linode = &clib.LinodeEnv{}
		default:
			result.Upload.Type = clib.UploadTypeMemory
		}

	}

	return result
}

//NewPaypalService returns a service that satisfies the clib.PaypalService interface
func NewPaypalService(env *clib.Env) clib.PaypalService {
	paypal, err := paypal.NewService(env)
	if err != nil {
		panic(err)
	}
	return paypal
}

//NewAuthService returns a service that satisfies the clib.AuthService interface
func NewAuthService(env *clib.Env) clib.AuthService {
	service, err := auth.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

//NewOrderService returns a service that satisfies the clib.OrderService
func NewOrderService(env *clib.Env) clib.OrderService {
	order, err := order.NewService(env)
	if err != nil {
		panic(err)
	}
	return order
}

//NewProductService returns a new product service
func NewProductService(env *clib.Env) clib.ProductService {
	service, err := product.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

//NewCartService returns a new cart service
func NewCartService(env *clib.Env) clib.CartService {
	service, err := cart.NewService(env)
	if err != nil {
		panic(err)
	}
	return service
}

func NewS3Client(env *clib.Env) (client *s3.Client) {
	client = s3.NewFromConfig(aws.Config{
		Region: env.AWSEnv.Region,
		Credentials: aws.CredentialsProviderFunc(func(ctx context.Context) (creds aws.Credentials, err error) {
			creds = aws.Credentials{
				AccessKeyID:     env.AWSEnv.AccessKey,
				SecretAccessKey: env.AWSEnv.SecretKey,
			}
			return
		}),
	})
	return
}

//NewUploadService will return a `upload` service on success,
//otherwise it will panic and close the application
func NewUploadService(env *clib.Env) clib.UploadService {
	service, err := upload.NewService(env)
	if err != nil {
		panic(err)
	}

	return service
}
