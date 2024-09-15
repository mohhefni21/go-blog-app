package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/interaction/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AddInteractions(ctx context.Context, model entity.InteractionEntity) (idInteraction int, err error)
	GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error)
	DeleteInteractionById(ctx context.Context, idInteraction int) (err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddInteractions(ctx context.Context, model entity.InteractionEntity) (idInteraction int, err error) {
	query := `
		INSERT INTO interactions (
			post_id, user_id, type, created_at, updated_at
		) VALUES (
			:post_id, :user_id, :type, :created_at, :updated_at
		) RETURNING interaction_id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &idInteraction, model)
	if err != nil {
		return
	}

	return
}

func (r *repository) GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error) {
	query := `
		SELECT
			user_id, public_id, username, email
		FROM users
		WHERE public_id=$1
	`

	err = r.db.GetContext(ctx, &model, query, publicId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserEntity{}, errorpkg.ErrorNotFound
		}
		return
	}

	return
}

func (r *repository) DeleteInteractionById(ctx context.Context, idInteraction int) (err error) {
	query := `
			DELETE 
			FROM 
				interactions
			WHERE
				interaction_id=$1
			`

	_, err = r.db.ExecContext(ctx, query, idInteraction)
	if err != nil {
		if err == sql.ErrNoRows {
			return errorpkg.ErrorNotFound
		}
		return
	}

	return
}
