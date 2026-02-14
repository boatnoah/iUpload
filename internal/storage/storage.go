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
	ErrorNotFound        = errors.New("Did not find record with that id")
	QueryTimeoutDuration = time.Second * 5
)

type Storage struct {
	ImageStorage interface {
		Create(context.Context, uuid.UUID, string, string) (*Image, error)
		GetById(context.Context, uuid.UUID) (*Image, error)
		DeleteById(context.Context, string) error
	}
	UserStorage interface {
		Create(context.Context, *User) error
		GetUserById(context.Context, string) (*User, error)
		GetUserByUserName(context.Context, string) (*User, error)
	}
	SessionStorage interface {
		CreateSession(context.Context, uuid.UUID) (*Session, error)
		GetIDbyToken(context.Context, string) (string, error)
	}
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{
		ImageStorage:   &Images{db},
		UserStorage:    &UserStore{db},
		SessionStorage: &Sessions{db},
	}
}
