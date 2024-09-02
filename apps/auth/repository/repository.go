package repository

import (
	"context"
	"mohhefni/go-blog-app/apps/auth/entity"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AddUser(ctx context.Context, model entity.AuthEntity) (email string, err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddUser(ctx context.Context, model entity.AuthEntity) (email string, err error) {
	return
}
