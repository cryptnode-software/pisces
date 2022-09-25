package main

import (
	"context"
	"log"
	"net/http"
	"os"

	commons "github.com/cryptnode-software/commons/pkg"
	pisces "github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/services"
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
		"/pisces.Pisces/StartUpload":               true,
	}

	admin = map[string]bool{
		"/pisces.Pisces/GetInquires": true,
	}
)

func main() {

	port := flag.Int("port", 4081, "grpc port")

	flag.Parse()

	environ := commons.Environment(os.Getenv(env))
	if environ == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
		return
	}

	environment := pisces.NewEnv(commons.NewLogger(environ))

	gw, err := pisces.NewGateway(environment, services.New(environment))
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
