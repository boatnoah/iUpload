package storage

import (
	"context"
	"database/sql"
)

// TODO model the actual table

type Images struct {
	db *sql.DB
}

func (i *Images) GetById(ctx context.Context, id string) error {
	return nil
}
func (i *Images) DeleteById(ctx context.Context, id string) error {
	return nil
}
