package lib

import (
	"github.com/gocraft/dbr/v2"
)

const (
	// EnvDev for development environment
	EnvDev = "dev"

	// EnvUAT for UAT environment
	EnvUAT = "uat"

	// EnvProd for production environment
	EnvProd = "prod"
)

// Env ...
type Env struct {
	DB          *dbr.Connection
	Log         Logger
	Environment string
	PaypalEnv   *PaypalEnv
	JWTEnv      *JWTEnv
}

//PaypalEnv the structure for the paypal environment
type PaypalEnv struct {
	ClientID string
	SecretID string
	Host     string
}

//JWTEnv the structure that is required for JWT configuration
type JWTEnv struct {
	Secret string
}
