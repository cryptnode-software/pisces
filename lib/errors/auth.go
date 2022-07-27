package errors

import (
	"errors"
	"fmt"
)

var (
	//ErrNoUsernameOrEmailProvided ...
	ErrNoUsernameOrEmailProvided = errors.New("no username or email was provided, please provide one or the other")

	//ErrNoUserFound ...
	ErrNoUserFound = errors.New("no user found with the provided email or username, please try another one")

	//ErrInvalidPassword ...
	ErrInvalidPassword = errors.New("the password provided is invalid, please provide a different one")

	//ErrNoMetadata ...
	ErrNoMetadata = errors.New("no metadata was provided in context please provide one")
)

//ErrInvalidHeader ...
type ErrInvalidHeader struct {
	Header string
}

func (err *ErrInvalidHeader) Error() string {
	return fmt.Sprintf("no header with the tag %s was provided, please provide one", err.Header)
}

//ErrNoAdminAccess ...
type ErrNoAdminAccess struct {
	Username string
}

func (err ErrNoAdminAccess) Error() string {
	return fmt.Sprintf("the user %+v doesn't have access to the request route", err.Username)
}
