package storage

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

type Images struct {
	db *sql.DB
}

type Image struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	ObjectKey   string    `json:"object_key"`
	ContentType string    `json:"content_type"`
	CreatedAt   string    `json:"created_at"`
}

func (i *Images) Create(ctx context.Context, userId uuid.UUID, objectKey, contentType string) (*Image, error) {
	query := `
		insert into images (user_id, object_key, content_type)
		values ($1, $2, $3) returning id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var image Image

	image.UserID = userId
	image.ObjectKey = objectKey
	image.ContentType = contentType

	err := i.db.QueryRowContext(
		ctx,
		query,
		userId,
		objectKey,
		contentType,
	).Scan(
		&image.ID,
		&image.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &image, nil

}

func (i *Images) GetById(ctx context.Context, id uuid.UUID) (*Image, error) {
	query := `
		select * from images 
		where id = $1
	`
	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var image Image

	err := i.db.QueryRowContext(ctx, query, id).Scan(
		&image.ID,
		&image.UserID,
		&image.ObjectKey,
		&image.ContentType,
		&image.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &image, nil

}

func (i *Images) DeleteById(ctx context.Context, id string) error {
	query := `
		delete from images
		where id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	res, err := i.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if rows == 0 {
		return ErrorNotFound
	}

	return nil

}
