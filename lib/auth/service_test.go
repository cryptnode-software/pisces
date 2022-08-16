package auth_test

import (
	"context"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/auth"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	password         = "$s2$16384$8$1$RumzdAfkvTNNHOPWev+d0L5c$AWSyfV0gxaF+yrQHzG5kRE/gzHH/waRTjMwg1qQachs="
	unhashedpassword = "testpassword"

	user = &lib.User{
		Email:    "testuser@test.com",
		Username: "testuser",
		Admin:    false,
		Model: lib.Model{
			ID: uuid.New(),
		},
	}

	newuser = struct {
		password string
		*lib.User
	}{}
)

var env = utility.NewEnv(utility.NewLogger())
var service, err = auth.NewService(env)

func TestGenerateAndDecodeJWT(t *testing.T) {
	ctx := context.Background()

	token, err := service.GenerateJWT(ctx, user)

	if err != nil {
		t.Error(err)
		return
	}

	if token == "" {
		t.Error("token was returned empty")
	}

	u, err := service.DecodeJWT(ctx, token)
	if err != nil {
		t.Error(err)
		return
	}

	if u == nil {
		t.Errorf("no user was returned from token")
		return
	}

	assert.Equal(t, user, u)
}

func TestLoginUser(t *testing.T) {

	req := &lib.LoginRequest{
		Password: unhashedpassword,
		Username: user.Username,
	}

	ctx := context.Background()

	user, err := service.Login(ctx, req)
	if err != nil {
		t.Error(err)
		return
	}

	if user == nil {
		t.Error("no user returned")
		return
	}
}

func TestFailedLogin(t *testing.T) {

	req := &lib.LoginRequest{
		Password: "fakepassword",
		Username: user.Username,
	}

	ctx := context.Background()

	_, err := service.Login(ctx, req)
	if err == nil {
		t.Errorf("login failed to fail")
		return
	}

}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	user, err := service.CreateUser(ctx, user, unhashedpassword)

	if err != nil {
		t.Error(err)
		return
	}

	if user == nil {
		t.Error("failed to create user")
		return
	}
}
