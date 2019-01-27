package books

import (
	"context"
	"database/sql"
)

type HandlerBooks interface {
	AddBooks(ctx context.Context) error
	EditBooks(ctx context.Context, id string) error
	DeleteBooks(ctx context.Context, id string) error
}

type service struct {
	db   *sql.DB
}

func CreateInstance(db  *sql.DB) HandlerBooks {
	return &service{
		db:db,
	}
}

func (s *service) AddBooks(ctx context.Context) error {
	return nil
}

func (s *service) EditBooks(ctx context.Context, id string) error {
	return nil
}

func (s *service) DeleteBooks(ctx context.Context, id string) error {
	return nil
}