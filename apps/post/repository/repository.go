package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/post/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	VerifyAvailableTitle(ctx context.Context, title string) (err error)
	AddPost(ctx context.Context, model entity.PostEntity) (idPost int, err error)
	UpdateCover(ctx context.Context, cover string, idPost int) (err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) VerifyAvailableTitle(ctx context.Context, title string) (err error) {
	query := `
		SELECT
			1
		FROM posts
		WHERE title=$1
	`
	var exits int8
	err = r.db.QueryRowContext(ctx, query, title).Scan(&exits)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return
	}

	return errorpkg.ErrTitleAlreadyUsed
}

func (r *repository) AddPost(ctx context.Context, model entity.PostEntity) (idPost int, err error) {
	query := `
		INSERT INTO posts (
			user_id, title, slug, excerpt, content, published_at, status, created_at, updated_at
		) VALUES (
			:user_id, :title, :slug, :excerpt, :content, :published_at, :status, :created_at, :updated_at
		) RETURNING post_id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &idPost, &model)
	if err != nil {
		return
	}

	return
}

func (r *repository) UpdateCover(ctx context.Context, cover string, idPost int) (err error) {
	query := `
		UPDATE posts
		SET cover=$1
		WHERE post_id=$2
	`

	_, err = r.db.ExecContext(ctx, query, cover, idPost)
	if err != nil {
		return
	}

	return
}
