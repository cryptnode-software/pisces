package auth

import "github.com/cryptnode-software/pisces/lib"

type user struct {
	Password string
	lib.User
}
