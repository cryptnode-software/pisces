package auth

import (
	"context"
	"testing"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/utility"
	"github.com/stretchr/testify/assert"
	"gopkg.in/hlandau/passlib.v1"
)

var (
	password         = "$s2$16384$8$1$RumzdAfkvTNNHOPWev+d0L5c$AWSyfV0gxaF+yrQHzG5kRE/gzHH/waRTjMwg1qQachs="
	unhashedpassword = "testpassword"

	user = &lib.User{
		Email:    "testuser@test.com",
		Username: "testuser",
		Admin:    false,
		ID:       0,
	}

	newuser = struct {
		password string
		*lib.User
	}{}
)

var env = utility.NewEnv(utility.NewLogger())
var service, err = NewService(env)

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
	service.repo = &mockrepo{}

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
	service.repo = &mockrepo{}

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
	service.repo = &mockrepo{}
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

type mockrepo struct{}

func (r *mockrepo) CreateUser(ctx context.Context, user *lib.User, password string) (*lib.User, error) {
	newuser = struct {
		password string
		*lib.User
	}{
		password,
		user,
	}

	newuser.ID = 0

	return newuser.User, nil
}

func (r *mockrepo) Login(ctx context.Context, req *lib.LoginRequest) (*lib.User, error) {

	user, err := r.FindUser(ctx, req.Username, "")
	if err != nil {
		return nil, err
	}

	_, err = passlib.Verify(req.Password, password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *mockrepo) FindUser(ctx context.Context, username, email string) (*lib.User, error) {

	if username == user.Username {
		return user, nil
	}

	if email == user.Email {
		return user, nil
	}

	return nil, nil
}
