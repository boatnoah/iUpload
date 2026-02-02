package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("resource not found")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	ImageStorage interface {
		GetById(context.Context, string) error
		DeleteById(context.Context, string) error
	}
	UserStorage interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, string) (*User, error)
		GetUserByUserName(context.Context, string) (*User, error)
	}
	SessionStorage interface {
		CreateSession(context.Context, uuid.UUID) (*Session, error)
		FindToken(context.Context, string) (bool, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		ImageStorage:   &Images{db},
		UserStorage:    &UserStore{db},
		SessionStorage: &Sessions{db},
	}
}
