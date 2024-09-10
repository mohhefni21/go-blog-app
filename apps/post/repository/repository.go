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
	GetDataPosts(ctx context.Context, model entity.PostsPaginationEntity) (posts []entity.GetListPostsEntity, err error)
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

func (r *repository) GetDataPosts(ctx context.Context, model entity.PostsPaginationEntity) (posts []entity.GetListPostsEntity, err error) {
	var query string
	if model.Search == "" {
		query = `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
            posts.published_at, users.fullname, users.username, users.picture
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.post_id > $1
        ORDER BY 
            posts.post_id DESC
        LIMIT $2
        `
		err := r.db.SelectContext(ctx, &posts, query, model.Cursor, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	} else {
		query = `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.slug, posts.excerpt, 
            posts.published_at, users.fullname, users.username, users.picture
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.post_id > $1
        AND
            posts.title ILIKE $2
        ORDER BY 
            posts.post_id DESC
        LIMIT $3
        `
		searchParam := "%" + model.Search + "%"
		err := r.db.SelectContext(ctx, &posts, query, model.Cursor, searchParam, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}

	return posts, nil
}
