package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/auth/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AddUser(ctx context.Context, model entity.UserEntity) (email string, err error)
	VerifyAvailableUsername(ctx context.Context, username string) (err error)
	VerifyAvailableEmail(ctx context.Context, email string) (err error)
	GetUserByEmail(ctx context.Context, email string) (model entity.UserEntity, err error)
	AddAuthentication(ctx context.Context, model entity.AuthEntity) (err error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) AddUser(ctx context.Context, model entity.UserEntity) (email string, err error) {
	query := `
		INSERT INTO users (
			username, fullname, email, password, created_at, updated_at
		) VALUES (
			:username, :fullname, :email, :password, :created_at, :updated_at
		) RETURNING email
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	err = stmt.GetContext(ctx, &email, model)
	if err != nil {
		return
	}

	return
}

func (r *repository) VerifyAvailableEmail(ctx context.Context, email string) (err error) {
	query := `
		SELECT 
			1
		FROM users
		WHERE email=$1
	`

	var exists int
	err = r.db.QueryRowContext(ctx, query, email).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err
	}

	return errorpkg.ErrEmailAlreadyUsed
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
			return nil
		}
		return err
	}

	return errorpkg.ErrUsernameAlreadyUsed
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (model entity.UserEntity, err error) {
	query := `
	SELECT
			*
		FROM users
		WHERE email=$1
	`

	err = r.db.GetContext(ctx, &model, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserEntity{}, errorpkg.ErrorNotFound
		}
		return
	}

	return
}

func (r *repository) AddAuthentication(ctx context.Context, model entity.AuthEntity) (err error) {
	query := `
		INSERT INTO authentications (
			user_id, refresh_token, refresh_token_expires_at
		) VALUES (
			:user_id, :refresh_token, :refresh_token_expires_at
		)
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.ExecContext(ctx, &model)
	if err != nil {
		return
	}

	return
}
