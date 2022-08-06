package errors

import "errors"

var (
	//ErrLinodeEnvNull is the generic error that we return when the server is first initializing, linode
	//is the primary upload client, and the configuration for linode is set as a null pointer
	ErrLinodeEnvNull = errors.New("linode environment is null but is set to the primary upload provider, please configure it properly")
	//ErrLinodePersonalAccessKeyInvalid should be returned when the linode personal access key is invalid
	ErrLinodePersonalAccessKeyInvalid = errors.New("linode personal access key invalid")
	//ErrLinodeAccessKeyInvalid should be returned when the linode access key is invalid
	ErrLinodeAccessKeyInvalid = errors.New("linode access key invalid")
	//ErrLinodeSecretKeyInvalid should be returned when the linode secret key is invalid
	ErrLinodeSecretKeyInvalid = errors.New("linode secret key invalid")
	//ErrLinodeEndpointInvalid should be returned when the linode endpoint is invalid
	ErrLinodeEndpointInvalid = errors.New("linode endpoint invalid")
)

//ErrLinodeEnvInvalid is returned when one of the required properties is
type ErrLinodeEnvInvalid struct {
	PersonalAccessKey bool
	AccessKey         bool
	SecretKey         bool
	Endpoint          bool
}
