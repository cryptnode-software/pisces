package auth_test

import (
	"context"
	"testing"

	commons "github.com/cryptnode-software/commons/pkg"
	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/auth"
	"github.com/stretchr/testify/assert"
)

var (
	ctx              = context.Background()
	unhashedpassword = "testpassword"

	testuser = &lib.User{
		Email:    "testuser@test.com",
		Username: "testuser",
		Admin:    false,
	}

	newuser = &user{
		password: unhashedpassword,
		User:     testuser,
	}
)

var env = lib.NewEnv(commons.NewLogger(commons.EnvDev))
var service, err = auth.NewService(env)

func TestGenerateAndDecodeJWT(t *testing.T) {

	token, err := service.GenerateJWT(ctx, testuser)

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

	assert.Equal(t, testuser, u)
}

// TestLoginUser tests against the testuser defined above
func TestLoginUser(t *testing.T) {

	tables := []struct {
		user *user
		fail bool
	}{
		{
			user: &user{
				password: "incorrect_password",
				User: &lib.User{
					Username: newuser.Username,
				},
			},
			fail: true,
		},
		{
			user: &user{
				password: newuser.password,
				User: &lib.User{
					Username: "incorrect.email@incorrect.com",
				},
			},
			fail: true,
		},
		{
			user: &user{
				password: newuser.password,
				User: &lib.User{
					Username: newuser.Username,
				},
			},
			fail: false,
		},
	}

	seed([]*user{
		newuser,
	})

	for _, table := range tables {

		req := &lib.LoginRequest{
			Username: table.user.Username,
			Password: table.user.password,
		}

		ctx := context.Background()

		{
			user, err := service.Login(ctx, req)

			if !table.fail && err != nil {
				t.Error(err)
				return
			}

			if (table.fail && err == nil) || (table.fail && user != nil) {
				t.Error("login user was suppose to fail when it didn't")
			}

			if !table.fail && user == nil {
				t.Error("no user returned")
				return
			}

		}

	}

	deseed([]*user{
		newuser,
	})
}

func TestCreateUser(t *testing.T) {

	tables := []struct {
		user *user
		fail bool
	}{
		{
			user: &user{
				password: "",
				User: &lib.User{
					Username: testuser.Username,
					Email:    testuser.Email,
				},
			},
			fail: true,
		},
		{
			user: &user{
				password: unhashedpassword,
				User: &lib.User{
					Username: "",
					Email:    testuser.Email,
				},
			},
			fail: true,
		},
		{
			user: &user{
				password: unhashedpassword,
				User: &lib.User{
					Username: testuser.Username,
					Email:    "",
				},
			},
			fail: true,
		},
		{
			user: &user{
				password: unhashedpassword,
				User: &lib.User{
					Username: testuser.Username,
					Email:    testuser.Email,
				},
			},
			fail: false,
		},
	}
	for _, table := range tables {
		u, err := service.CreateUser(ctx, table.user.User, table.user.password)

		if (table.fail && err == nil) || (table.fail && u != nil) {
			t.Error("create user was suppose to fail but didn't")
			return
		}

		if !table.fail && err != nil {
			t.Error(err)
			return
		}

		if !table.fail && u == nil {
			t.Error("failed to create user")
			return
		}

		deseed([]*user{
			{
				User: u,
			},
		})
	}
}

func TestAuthenticateAdmin(t *testing.T) {
	ctx := context.Background()

	tables := []struct {
		user *user
		fail bool
	}{
		{
			user: &user{
				User: &lib.User{
					Username: newuser.Username,
					Email:    newuser.Email,
					Admin:    true,
				},
				password: newuser.password,
			},
			fail: false,
		},
		{
			user: &user{
				User: &lib.User{
					Username: newuser.Username,
					Email:    newuser.Email,
					Admin:    false,
				},
				password: newuser.password,
			},
			fail: true,
		},
	}

	for _, table := range tables {
		err := seed([]*user{
			table.user,
		})

		if err != nil {
			t.Error(err)
			return
		}

		token, err := service.GenerateJWT(ctx, table.user.User)

		if err != nil {
			t.Error(err)
			return
		}

		ctx = lib.SetAuthContext(ctx, token)

		_, err = service.AuthenticateAdmin(ctx)

		if table.fail && err == nil {
			t.Error("authenticate user failed to fail")
			return
		}

		if !table.fail && err != nil {
			t.Error(err)
			return
		}

		err = deseed([]*user{
			table.user,
		})

		if err != nil {
			t.Error(err)
			return
		}

	}

	token, err := service.GenerateJWT(ctx, testuser)

	if err != nil {
		t.Error(err)
		return
	}

	if token == "" {
		t.Error("token was returned empty")
	}
}

func TestSoftDeleteUser(t *testing.T) {
	tables := []struct {
		user *user
		fail bool
	}{
		{
			user: newuser,
			fail: false,
		},
	}

	for _, table := range tables {

		seed([]*user{
			table.user,
		})

		err := service.DeleteUser(ctx, table.user.User, nil)

		if !table.fail && err != nil {
			t.Error(err)
			return
		}

		deseed([]*user{
			table.user,
		})
	}
}

func seed(users []*user) error {

	for _, user := range users {
		u, err := service.CreateUser(ctx, user.User, user.password)

		if err != nil {
			return err
		}

		user.User = u
	}

	return nil
}

func deseed(users []*user) error {
	for _, user := range users {
		if err := service.DeleteUser(ctx, user.User, &lib.DeleteConditions{
			HardDelete: true,
		}); err != nil {
			return err
		}
	}
	return nil
}

type user struct {
	password string
	*lib.User
}
