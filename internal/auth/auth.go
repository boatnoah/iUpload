package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/boatnoah/iupload/internal/storage"
)

var (
	ErrorInvalidPayload   = errors.New("Username and password are both required")
	ErrorPasswordMismatch = errors.New("Password does not match")
)

type Auth struct {
	storage *storage.Storage
}

func New(storage *storage.Storage) *Auth {
	return &Auth{storage}
}

type UserPayload struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type UserLoginPayload struct {
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

func normalizeUsername(username string) string {
	return strings.ToLower(strings.TrimSpace(username))
}

func (a *Auth) RegisterUser(ctx context.Context, userPayload UserPayload) (*storage.User, *storage.Session, error) {

	var user storage.User

	user.FirstName = userPayload.FirstName
	user.LastName = userPayload.LastName
	user.UserName = normalizeUsername(userPayload.UserName)

	var password storage.Password

	err := password.Set(userPayload.Password)

	if err != nil {
		return nil, nil, err
	}

	user.HashedPassword = password

	err = a.storage.UserStorage.Create(ctx, &user)

	if err != nil {
		return nil, nil, err
	}

	session, err := a.storage.SessionStorage.CreateSession(ctx, user.ID)

	if err != nil {
		return nil, nil, err
	}

	return &user, session, nil
}

func (a *Auth) LogInUser(ctx context.Context, userLoginPayload UserLoginPayload) (*storage.Session, error) {

	if userLoginPayload.UserName == "" || userLoginPayload.Password == "" {
		return nil, ErrorInvalidPayload
	}

	user, err := a.storage.UserStorage.GetUserByUserName(ctx, normalizeUsername(userLoginPayload.UserName))

	if err != nil {
		return nil, err
	}

	err = user.HashedPassword.Compare(userLoginPayload.Password)
	if err != nil {
		return nil, ErrorPasswordMismatch
	}

	session, err := a.storage.SessionStorage.CreateSession(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return session, nil
}

func (a *Auth) ValidateUser(ctx context.Context, token string) (bool, error) {
	exist, err := a.storage.SessionStorage.FindToken(ctx, token)
	return exist, err
}
