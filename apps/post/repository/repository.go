package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/post/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	VerifyAvailableTitle(ctx context.Context, title string) (err error)
	AddPost(ctx context.Context, model entity.PostEntity) (idPost int, err error)
	UpdateCover(ctx context.Context, cover string, idPost int) (err error)
	GetDataPosts(ctx context.Context, model entity.PostsPaginationEntity) (posts []entity.GetListPostsEntity, err error)
	GetDetailPostBySLug(ctx context.Context, slug string) (postDetail entity.GetDetailPostResponseEntity, err error)
	GetPostById(ctx context.Context, idPost int) (postDetail entity.PostEntity, err error)
	VerifyAvailableUsername(ctx context.Context, username string) (err error)
	GetDataPostsByUsername(ctx context.Context, model entity.PostsPaginationEntity, username string) (posts []entity.GetListPostsEntity, err error)
	GetDataPostsByUserLogin(ctx context.Context, publicId uuid.UUID) (posts []entity.GetListPostsByUserLoginEntity, err error)
	DeletePostById(ctx context.Context, idPost int) (err error)
	UpdatePostById(ctx context.Context, model entity.PostEntity) (err error)
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

func (r *repository) GetDetailPostBySLug(ctx context.Context, slug string) (postDetail entity.GetDetailPostResponseEntity, err error) {
	query := `
        SELECT
            posts.post_id, posts.cover, posts.title, posts.content, posts.published_at, 
			users.fullname AS "author.fullname", 
			users.username AS "author.username", 
			users.picture AS "author.picture"
        FROM 
            posts
        INNER JOIN
            users ON posts.user_id = users.user_id
        WHERE 
            posts.slug=$1
        `

	err = r.db.GetContext(ctx, &postDetail, query, slug)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errorpkg.ErrNotFound
			return
		}
		return
	}

	return
}

func (r *repository) VerifyAvailableUsername(ctx context.Context, username string) (err error) {
	query := `
		SELECT 
			1
		FROM users
		WHERE username=$1
	`

	var exists int
	err = r.db.QueryRowContext(ctx, query, username).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorpkg.ErrNotFound
		}
		return err
	}

	return nil
}

func (r *repository) GetDataPostsByUsername(ctx context.Context, model entity.PostsPaginationEntity, username string) (posts []entity.GetListPostsEntity, err error) {
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
				 users.username=$1
			AND
				posts.post_id > $2
			ORDER BY 
				posts.post_id DESC
			LIMIT $3
			`
		err := r.db.SelectContext(ctx, &posts, query, username, model.Cursor, model.Limit)
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
				 users.username=$1
			AND
				posts.post_id > $2
			ORDER BY 
				posts.post_id DESC
			LIMIT $3
			`
		searchParam := "%" + model.Search + "%"
		err := r.db.SelectContext(ctx, &posts, query, username, model.Cursor, searchParam, model.Limit)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			return nil, err
		}
	}

	return posts, nil
}

func (r *repository) GetDataPostsByUserLogin(ctx context.Context, publicId uuid.UUID) (posts []entity.GetListPostsByUserLoginEntity, err error) {
	query := `
			SELECT
				posts.post_id, posts.cover, posts.title, posts.slug, posts.status, 
				posts.published_at, posts.created_at
			FROM 
				posts
			INNER JOIN
				users ON posts.user_id = users.user_id
			WHERE
				 users.public_id=$1
			`
	err = r.db.SelectContext(ctx, &posts, query, publicId)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return posts, nil
}

func (r *repository) DeletePostById(ctx context.Context, idPost int) (err error) {
	query := `
			DELETE 
			FROM 
				posts
			WHERE
				post_id=$1
			`

	_, err = r.db.ExecContext(ctx, query, idPost)
	if err != nil {
		return
	}

	return
}

func (r *repository) GetPostById(ctx context.Context, idPost int) (postDetail entity.PostEntity, err error) {
	query := `
        SELECT
            post_id, user_id, cover, title, slug, excerpt, content, published_at, status, created_at, updated_at
        FROM 
            posts
        WHERE 
            post_id=$1
        `

	err = r.db.GetContext(ctx, &postDetail, query, idPost)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errorpkg.ErrNotFound
			return
		}
		return
	}

	return
}

func (r *repository) UpdatePostById(ctx context.Context, model entity.PostEntity) (err error) {
	query := `
        UPDATE
			posts
		SET
            title=:title, excerpt=:excerpt, slug=:slug, content=:content, status=:status, published_at=:published_at, updated_at=:updated_at
        WHERE 
            post_id=:post_id
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
