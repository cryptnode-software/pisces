package utility

import (
	"log"
	"os"

	"github.com/cryptnode-software/pisces/lib"
	clib "github.com/cryptnode-software/pisces/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	paylib "github.com/plutov/paypal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	envDatabaseURL string = "DB_CONNECTION"
	envBase        string = "ENV"

	envPaypalClientID string = "PAYPAL_CLIENT_ID"
	envPaypalSecretID string = "PAYPAL_SECRET_ID"

	envJWTSecret string = "JWT_SECRET"

	envS3Bucket    string = "S3_BUCKET"
	envS3AccessKey string = "AWS_ACCESS_KEY_ID"
	envS3SecretKey string = "AWS_SECRET_ACCESS_KEY"
	envS3Region    string = "AWS_REGION"
	envS3Endpoint  string = "AWS_ENDPOINT"
)

//NewLogger returns a new logger based off the current environment
func NewLogger() clib.Logger {
	env := os.Getenv(envBase)
	if env == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
	}

	return clib.NewZapper(clib.Environment(env))
}

//NewEnv returns a new environment pre-populated with the provided logger
func NewEnv(logger clib.Logger) *clib.Env {
	env := os.Getenv(envBase)
	if env == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
	}

	result := &clib.Env{
		Environment: clib.Environment(env),
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

	{

		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN:                       os.Getenv(envDatabaseURL),
			SkipInitializeWithVersion: false,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			DefaultStringSize:         256,
		}))

		if err != nil {
			log.Fatal(err)
			return nil
		}

		result.GormDB = db

		db.AutoMigrate(
			new(lib.Inquiry),
			new(lib.Order),
		)

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

		//optional aws config
		endpoint := os.Getenv(envS3Endpoint)
		if config.Endpoint = &endpoint; config.Endpoint == nil {
			logger.Error("%s is not set, defaulting to aws endpoint", envS3Endpoint)
		}
		result.AWSEnv = config
	}

	return result
}
