package main

import (
	"context"
	"log"
	"net/http"
	"os"

	pisces "github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/utility"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	proto "go.buf.build/grpc/go/thenewlebowski/pisces/general/v1"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"

	"flag"
	"fmt"
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

	envS3SecretKey string = "AWS_SECRET_ACCESS_KEY"
	envS3AccessKey string = "AWS_ACCESS_KEY_ID"
	envS3Endpoint  string = "AWS_ENDPOINT"
	envS3Region    string = "AWS_REGION"
	envS3Bucket    string = "S3_BUCKET"
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
		"/pisces.Pisces/StartUpload":               true,
	}

	admin = map[string]bool{
		"/pisces.Pisces/GetInquires": true,
	}
)

func main() {

	port := flag.Int("port", 4081, "grpc port")

	flag.Parse()

	environ := pisces.Environment(os.Getenv(env))
	if environ == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
		return
	}

	environment := pisces.NewEnv(pisces.NewLogger(environ))

	gw, err := pisces.NewGateway(environment, utility.Services(environment))
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
// func NewLogger() clib.Logger {
// 	environ := os.Getenv(env)
// 	if environ == "" {
// 		log.Fatalf("environment is not provided: please provide %s variable", env)
// 		return nil
// 	}

// 	return clib.NewZapper(clib.Environment(environ))
// }

// //NewEnv returns a new environment pre-populated with the provided logger
// func NewEnv(logger clib.Logger) *clib.Env {
// 	environ := os.Getenv(env)
// 	if environ == "" {
// 		log.Fatalf("environment is not provided: please provide %s variable", env)
// 		return nil
// 	}

// 	result := &clib.Env{
// 		Environment: clib.Environment(environ),
// 		Log:         logger,
// 	}

// 	//paypal config
// 	{
// 		client := os.Getenv(envPaypalClientID)
// 		if client == "" {
// 			log.Fatalf("paypal client id not provided please provide %s env variable", envPaypalClientID)
// 		}

// 		secret := os.Getenv(envPaypalSecretID)
// 		if secret == "" {
// 			log.Fatalf("paypal secret id not provided please provide %s env variable", envPaypalSecretID)
// 		}

// 		result.PaypalEnv = &clib.PaypalEnv{
// 			ClientID: client,
// 			SecretID: secret,
// 		}

// 		switch result.Environment {
// 		case clib.EnvProd:
// 			result.PaypalEnv.Host = paylib.APIBaseLive
// 		default:
// 			result.PaypalEnv.Host = paylib.APIBaseSandBox
// 		}
// 	}

// 	//auth config
// 	{
// 		jwtSecret := os.Getenv(envJWTSecret)
// 		if jwtSecret == "" {
// 			log.Fatalf("%s not set, if not properly set jwt tokens will be unsafe to use", envJWTSecret)
// 		}

// 		result.JWTEnv = &clib.JWTEnv{
// 			Secret: jwtSecret,
// 		}
// 	}

// 	//database config
// 	{

// 		var err error
// 		if result.GormDB, err = pgorm.NewDatabase(os.Getenv(envDatabaseURL)); err != nil {
// 			log.Fatalf("%+v", err)
// 			return nil
// 		}
// 	}

// 	//aws config
// 	{
// 		config := new(clib.AWSEnv)
// 		if config.Region = os.Getenv(envS3Region); config.Region == "" {
// 			log.Fatalf("%s not, and required for s3 configuration", envS3Region)
// 		}
// 		if config.AccessKey = os.Getenv(envS3AccessKey); config.AccessKey == "" {
// 			log.Fatalf("%s not set and required for s3 configuration", envS3AccessKey)
// 		}
// 		if config.Bucket = os.Getenv(envS3Bucket); config.Bucket == "" {
// 			log.Fatalf("%s not set and required for s3 configuration", envS3Bucket)
// 		}
// 		if config.SecretKey = os.Getenv(envS3SecretKey); config.SecretKey == "" {
// 			log.Fatalf("%s not set and required for s3 configuration", envS3SecretKey)
// 		}

// 		//optional aws configuration
// 		endpoint := os.Getenv(envS3Endpoint)
// 		if config.Endpoint = &endpoint; config.Endpoint == nil {
// 			logger.Error("%s is not set, defaulting to aws endpoint", envS3Endpoint)
// 		}

// 		result.AWSEnv = config
// 	}

// 	return result
// }
