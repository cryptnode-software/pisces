package utility

import (
	"log"
	"os"

	clib "github.com/cryptnode-software/pisces/lib"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr/v2"
	paylib "github.com/plutov/paypal"
)

const (
	envDatabaseURL string = "DB_CONNECTION"
	envBase        string = "ENV"

	envPaypalClientID string = "PAYPAL_CLIENT_ID"
	envPaypalSecretID string = "PAYPAL_SECRET_ID"

	envJWTSecret string = "JWT_SECRET"
)

//NewLogger returns a new logger based off the current environment
func NewLogger() clib.Logger {
	env := os.Getenv(envBase)
	if env == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
	}

	return clib.NewZapper(env)
}

//NewEnv returns a new environment pre-populated with the provided logger
func NewEnv(logger clib.Logger) *clib.Env {
	env := os.Getenv(envBase)
	if env == "" {
		log.Fatalf("environment is not provided: please provide %s variable", env)
	}

	result := &clib.Env{
		Environment: env,
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

		switch env {
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
		sql, err := dbr.Open("mysql", os.Getenv(envDatabaseURL), nil)
		if err != nil {
			log.Fatal(err)
		}
		result.DB = sql
	}

	return result
}
