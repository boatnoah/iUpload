package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrorDuplicateUserName = errors.New("Duplicate user name")
)

type UserStore struct {
	db *sql.DB
}

type User struct {
	ID             uuid.UUID `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	UserName       string    `json:"user_name"`
	HashedPassword Password  `json:"-"`
	CreatedAt      string    `json:"created_at"`
}

type Password struct {
	text *string
	hash []byte
}

func (p *Password) Set(text string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &text
	p.hash = hash
	return nil
}

func (p *Password) Compare(text string) error {
	return bcrypt.CompareHashAndPassword(p.hash, []byte(text))
}

func (u *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		insert into users (first_name, last_name, user_name, hashed_password)
		values ($1, $2, $3, $4) returning id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := u.db.QueryRowContext(
		ctx,
		query,
		user.FirstName,
		user.LastName,
		user.UserName,
		user.HashedPassword.hash).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrorDuplicateUserName
		}
	}
	return nil
}

func (u *UserStore) GetUserByUserName(ctx context.Context, username string) (*User, error) {

	query := `
		select id, first_name, last_name, user_name, hashed_password, created_at
		from users
		where user_name = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := u.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.HashedPassword,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (u *UserStore) GetUserById(ctx context.Context, uuid string) (*User, error) {

	query := `
		select id, first_name, last_name, user_name, hashed_password, created_at
		from users
		where id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User

	err := u.db.QueryRowContext(ctx, query, uuid).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.UserName,
		&user.HashedPassword,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil

}
