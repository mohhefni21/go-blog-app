package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/comment/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AddComment(ctx context.Context, model entity.CommentEntity) (idComment int, err error)
	GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error)
	UpdateCommentById(ctx context.Context, model entity.CommentEntity) (err error)
	DeleteCommentById(ctx context.Context, commentId int) (err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddComment(ctx context.Context, model entity.CommentEntity) (idComment int, err error) {
	query := `
		INSERT INTO comments (
			post_id, user_id, parent_id, content, created_at, updated_at
		) VALUES (
			:post_id, :user_id, :parent_id, :content, :created_at, :updated_at
		) RETURNING comment_id
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &idComment, model)
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

func (r *repository) UpdateCommentById(ctx context.Context, model entity.CommentEntity) (err error) {
	query := `
        UPDATE
			comments
		SET
            content=:content, updated_at=:updated_at
        WHERE 
            comment_id=:comment_id
        `

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, model)
	if err != nil {
		return
	}

	return
}

func (r *repository) DeleteCommentById(ctx context.Context, commentId int) (err error) {
	query := `
        DELETE FROM
			comments
        WHERE 
            comment_id=$1
		OR
			parent_id=$1
        `

	_, err = r.db.ExecContext(ctx, query, commentId)
	if err != nil {
		return
	}

	return
}
