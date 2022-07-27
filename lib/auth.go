package lib

import (
	"context"

	"github.com/cryptnode-software/pisces/lib/errors"
	"google.golang.org/grpc/metadata"
)

const (
	authheader = "auth"
)

//AuthService represents our internal auth service mostly used for authenticating customers and
//our admin employees
type AuthService interface {
	CreateUser(ctx context.Context, user *User, password string) (*User, error)
	DecodeJWT(ctx context.Context, token string) (*User, error)
	GenerateJWT(ctx context.Context, user *User) (string, error)
	AuthenticateToken(ctx context.Context) (*User, error)
	AuthenticateAdmin(ctx context.Context) (*User, error)
	Login(context.Context, *LoginRequest) (*User, error)
}

//LoginRequest holds the values that are required to properly login
type LoginRequest struct {
	Username string
	Email    string
	Password string
}

//User the general structure of a user through out the ecosystem
type User struct {
	Username string `json:"username" db:"username"`
	Email    string `json:"email" db:"email"`
	Admin    bool   `json:"admin" db:"admin"`
	ID       int32  `json:"id" db:"id"`
}

//Auth ...
type Auth struct {
}

//CreateUserResponse ...
type CreateUserResponse struct {
}

//CreateUserRequest ...
type CreateUserRequest struct {
}

//GetAuthFromContext gets the authentication token can be omitted by specifing the route that
//doesn't require authentication during gateway intialization
func GetAuthFromContext(ctx context.Context) (string, error) {
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.ErrNoMetadata
	}

	auth, ok := metadata[authheader]
	if !ok {
		return "", &errors.ErrInvalidHeader{
			Header: authheader,
		}
	}

	if len(auth) > 1 {
		return "", &errors.ErrInvalidHeader{
			Header: authheader,
		}
	}

	return auth[0], nil
}
