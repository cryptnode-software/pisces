package lib

import (
	"log"
	"os"

	commons "github.com/cryptnode-software/commons/pkg"
	pgorm "github.com/cryptnode-software/pisces/lib/gorm"
	paylib "github.com/plutov/paypal"
	"gorm.io/gorm"
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

// Env ...
type Env struct {
	GormDB      *gorm.DB
	Log         commons.Logger
	Environment commons.Environment
	PaypalEnv   *PaypalEnv
	JWTEnv      *JWTEnv
	AWSEnv      *AWSEnv
}

// PaypalEnv the structure for the paypal environment
type PaypalEnv struct {
	ClientID string
	SecretID string
	Host     string
}

// JWTEnv the structure that is required for JWT configuration
type JWTEnv struct {
	Secret string
}

// UploadType the primitive type that all of upload configurations support
type UploadType string

const (
	//UploadTypeMemory allows us to specify a memory upload, used in our tests
	UploadTypeMemory UploadType = "UPLOAD_MEMORY"
	//UploadTypeLinode allows us to specify a linode upload
	UploadTypeLinode UploadType = "UPLOAD_LINODE"
)

type AWSEnv struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
	Endpoint  *string
}

func NewEnv(logger commons.Logger) (result *Env) {
	environ := os.Getenv(env)
	var err error

	if environ == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
		return nil
	}

	result = &Env{
		Environment: commons.Environment(environ),
		Log:         logger,
	}

	result.PaypalEnv = NewPaypalEnv(
		result.Environment,
		os.Getenv(envPaypalClientID),
		os.Getenv(envPaypalSecretID),
	)

	result.JWTEnv = NewJWTEnv(os.Getenv(envJWTSecret))

	if result.GormDB, err = pgorm.NewDatabase(
		os.Getenv(envDatabaseURL),
	); err != nil {
		log.Fatalf("%+v", err)
		return nil
	}

	result.AWSEnv = NewAWSEnv()

	return
}

func NewPaypalEnv(env commons.Environment, client, secret string) (result *PaypalEnv) {

	if client == "" {
		log.Fatalf("paypal client id not provided please provide %s env variable", envPaypalClientID)
		return
	}

	if secret == "" {
		log.Fatalf("paypal secret id not provided please provide %s env variable", envPaypalSecretID)
		return
	}

	result = &PaypalEnv{
		ClientID: client,
		SecretID: secret,
	}

	switch env {
	case commons.EnvProd:
		result.Host = paylib.APIBaseLive
	default:
		result.Host = paylib.APIBaseSandBox
	}

	return
}

func NewJWTEnv(secret string) (env *JWTEnv) {
	env = new(JWTEnv)

	if secret == "" {
		log.Fatalf("%s not set, if not properly set jwt tokens will be unsafe to use", envJWTSecret)
	}

	return
}

func NewAWSEnv() (env *AWSEnv) {
	env = new(AWSEnv)
	if env.Region = os.Getenv(envS3Region); env.Region == "" {
		log.Fatalf("%s not, and required for s3 configuration", envS3Region)
	}
	if env.AccessKey = os.Getenv(envS3AccessKey); env.AccessKey == "" {
		log.Fatalf("%s not set and required for s3 configuration", envS3AccessKey)
	}
	if env.Bucket = os.Getenv(envS3Bucket); env.Bucket == "" {
		log.Fatalf("%s not set and required for s3 configuration", envS3Bucket)
	}
	if env.SecretKey = os.Getenv(envS3SecretKey); env.SecretKey == "" {
		log.Fatalf("%s not set and required for s3 configuration", envS3SecretKey)
	}

	//optional aws configuration
	endpoint := os.Getenv(envS3Endpoint)
	if env.Endpoint = &endpoint; env.Endpoint == nil {
		log.Printf("%s is not set, defaulting to aws endpoint", envS3Endpoint)
	}

	return
}
