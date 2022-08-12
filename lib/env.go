package lib

import (
	"github.com/cryptnode-software/pisces/lib/errors"
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
	Upload      *UploadEnv
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

//UploadEnv the
type UploadEnv struct {
	Type   UploadType
	Linode *LinodeEnv
}

//LinodeEnv is the req. structure for a proper
//linode upload configuration.
type LinodeEnv struct {
	PersonalAccessKey string
	SecretKey         string
	AccessKey         string
	Endpoint          string
}

//Validate validates our linode environment for the
//required properties. Doesn't do a hard check on
//optional properties
func (env *LinodeEnv) Validate() error {
	if env == nil {
		return errors.ErrLinodeEnvNull
	}

	if env.AccessKey == "" {
		return errors.ErrLinodeAccessKeyInvalid
	}

	if env.SecretKey == "" {
		return errors.ErrLinodeSecretKeyInvalid
	}

	if env.PersonalAccessKey == "" {
		return errors.ErrLinodePersonalAccessKeyInvalid
	}

	if env.Endpoint == "" {
		return errors.ErrLinodeEndpointInvalid
	}

	return nil
}

type AWSEnv struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Region    string
}
