package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/cryptnode-software/pisces/lib"
	"github.com/cryptnode-software/pisces/lib/errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gopkg.in/hlandau/passlib.v1"
	"gorm.io/gorm"
)

var (
	tables = struct {
		auth  string
		users string
	}{
		users: "users",
	}
)

// Service the auth structure handles anything that may be auth related
//
// Login handles the functionality to properly check and login a user
type Service struct {
	*lib.Env
	repo RepoI
}

// NewService creates a new paypal service that satisfies the PaypalService interface
func NewService(env *lib.Env) (lib.AuthService, error) {
	return &Service{
		env,
		&repo{
			env.GormDB,
		},
	}, nil
}

// Login accepts a login response with a valid username and password if they match then the jwt
// is hashed and returned in order to properly user the application
func (s *Service) Login(ctx context.Context, req *lib.LoginRequest) (*lib.User, error) {
	return s.repo.Login(ctx, req)
}

// CreateUser creates a user in the
func (s *Service) CreateUser(ctx context.Context, user *lib.User, password string) (*lib.User, error) {
	return s.repo.CreateUser(ctx, user, password)
}

// GenerateJWT creates a jwt and signs it with the secret that is collected from the JWTSecret env property
func (s *Service) GenerateJWT(ctx context.Context, user *lib.User) (string, error) {
	claims := struct {
		User *lib.User `json:"user"`
		jwt.RegisteredClaims
	}{
		user,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.Env.JWTEnv.Secret))
}

// DecodeJWT decodes a jwt
func (s *Service) DecodeJWT(ctx context.Context, token string) (*lib.User, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.Env.JWTEnv.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	var result *lib.User

	if claims, ok := t.Claims.(jwt.MapClaims); ok && t.Valid {
		switch u := claims["user"].(type) {
		case map[string]interface{}:

			result = new(lib.User)

			uuid, err := uuid.Parse(u["id"].(string))
			if err != nil {
				return nil, err
			}

			result.Username = u["username"].(string)
			result.Email = u["email"].(string)
			result.Admin = u["admin"].(bool)
			result.ID = uuid
		}

	} else {
		return nil, err
	}

	return result, nil
}

func (s *Service) DeleteUser(ctx context.Context, user *lib.User, conditions *lib.DeleteConditions) error {
	if conditions != nil {
		if conditions.HardDelete {
			return s.repo.HardDelete(ctx, user)
		}
	}
	return s.repo.SoftDelete(ctx, user)
}

// AuthenticateToken makes sure a token is valid and isn't expired otherwise it
// will raise an exception.
func (s *Service) AuthenticateToken(ctx context.Context) (*lib.User, error) {
	token, err := lib.GetAuthFromContext(ctx)
	if err != nil {
		return nil, err
	}

	return s.DecodeJWT(ctx, token)
}

// AuthenticateAdmin authenticates a request that is only to be used by admin personal
// doesn't only user the jwt token but double checks with the database information befor
// approval
func (s *Service) AuthenticateAdmin(ctx context.Context) (*lib.User, error) {
	user, err := s.AuthenticateToken(ctx)

	if err != nil {
		return nil, err
	}

	user, err = s.repo.FindUser(ctx, user.Username, user.Email)

	if err != nil {
		return nil, err
	}

	if !user.Admin {
		return nil, errors.ErrNoAdminAccess{Username: user.Username}
	}

	return user, nil
}

type RepoI interface {
	CreateUser(ctx context.Context, user *lib.User, password string) (*lib.User, error)
	FindUser(ctx context.Context, username, email string) (*lib.User, error)
	Login(context.Context, *lib.LoginRequest) (*lib.User, error)
	HardDelete(ctx context.Context, user *lib.User) error
	SoftDelete(ctx context.Context, user *lib.User) error
}

type repo struct {
	*gorm.DB
}

func (r *repo) CreateUser(ctx context.Context, luser *lib.User, password string) (*lib.User, error) {
	if luser.Username == "" {
		return nil, errors.ErrNoUsernameOrEmailProvided
	}

	if luser.Email == "" {
		return nil, errors.ErrNoUsernameOrEmailProvided

	}

	if password == "" {
		return nil, errors.ErrInvalidPassword
	}

	hash, err := passlib.Hash(password)
	if err != nil {
		return nil, err
	}

	entry := new(user)

	entry.Password = hash
	entry.User = luser

	r.DB.Model(new(user)).Create(entry)

	return entry.User, nil

}

func (r *repo) Login(ctx context.Context, req *lib.LoginRequest) (*lib.User, error) {
	user, err := r.FindUser(ctx, req.Username, req.Email)
	if err != nil {
		return nil, err
	}

	hash := ""
	err = r.DB.Model(new(lib.User)).Select("password").Where("id = ?", user.ID).First(&hash).Error
	if err != nil {
		return nil, err
	}

	if hash == "" {
		return nil, errors.ErrNoUserFound
	}

	_, err = passlib.Verify(req.Password, hash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repo) FindUser(ctx context.Context, username, email string) (*lib.User, error) {

	if username == "" && email == "" {
		return nil, errors.ErrNoUsernameOrEmailProvided
	}

	tx := r.DB.Model(new(lib.User))

	if username != "" {
		tx = tx.Where("username = ?", username)
	}

	if email != "" {
		tx = tx.Where("email = ?", email)
	}

	user := new(lib.User)

	err := tx.First(user).Error

	if err != nil {
		return nil, err
	}

	if username != "" {
		if username != user.Username {
			return nil, errors.ErrNoUserFound
		}
	}

	if email != "" {
		if email != user.Email {
			return nil, errors.ErrNoUserFound
		}
	}

	return user, nil
}

func (r *repo) HardDelete(ctx context.Context, user *lib.User) error {
	return r.DB.Unscoped().Delete(user).Error
}
func (r *repo) SoftDelete(ctx context.Context, user *lib.User) error {
	return r.DB.Delete(user).Error
}
