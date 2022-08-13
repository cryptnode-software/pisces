package lib

import (
	"github.com/gocraft/dbr/v2"
)

type Environment string

const (
	// EnvDev for development environment
	EnvDev Environment = "dev"

	// EnvUAT for UAT environment
	EnvUAT Environment = "uat"

	// EnvProd for production environment
	EnvProd Environment = "prod"
)

// Env ...
type Env struct {
	DB          *dbr.Connection
	Log         Logger
	Environment Environment
	PaypalEnv   *PaypalEnv
	JWTEnv      *JWTEnv
	AWSEnv      *AWSEnv
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

//UploadType the primitive type that all of upload configurations support
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
