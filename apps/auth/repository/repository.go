package repository

import (
	"context"
	"database/sql"
	"errors"
	"mohhefni/go-blog-app/apps/auth/entity"
	"mohhefni/go-blog-app/infra/errorpkg"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	AddUser(ctx context.Context, model entity.UserEntity) (email string, err error)
	VerifyAvailableUsername(ctx context.Context, username string) (err error)
	VerifyAvailableUsernameByEmail(ctx context.Context, email string, username string) (err error)
	VerifyAvailableEmail(ctx context.Context, email string) (err error)
	VerifyAvailableToken(ctx context.Context, refreshToken string) (err error)
	GetUserByEmail(ctx context.Context, email string) (model entity.UserEntity, err error)
	GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error)
	AddAuthentication(ctx context.Context, model entity.AuthEntity) (err error)
	DeleteAuthenticationById(ctx context.Context, idUser int) (err error)
	DeleteAuthenticationRefreshToken(ctx context.Context, refreshToken string) (err error)
	UpdateProfileOnboarding(ctx context.Context, email string, model entity.UserEntity) (err error)
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
			public_id, username, fullname, email, password, role, picture, created_at, updated_at
		) VALUES (
			:public_id, :username, :fullname, :email, :password, :role, :picture, :created_at, :updated_at
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

func (r *repository) VerifyAvailableUsernameByEmail(ctx context.Context, username string, email string) (err error) {
	query := `
		SELECT 
			1
		FROM users
		WHERE username=$1 AND email != $2
	`

	var exists int
	err = r.db.QueryRowContext(ctx, query, username, email).Scan(&exists)
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

func (r *repository) GetUserByPublicId(ctx context.Context, publicId uuid.UUID) (model entity.UserEntity, err error) {
	query := `
		SELECT
			*
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

func (r *repository) VerifyAvailableToken(ctx context.Context, refreshToken string) (err error) {
	query := `
	SELECT 
			1
		FROM authentications
		WHERE refresh_token=$1
	`

	var exists int
	err = r.db.QueryRowContext(ctx, query, refreshToken).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errorpkg.ErrUnauthorized
		}
		return
	}

	return
}

func (r *repository) DeleteAuthenticationById(ctx context.Context, idUser int) (err error) {
	query := `
		DELETE 
		FROM 
			authentications
		WHERE user_id=$1
	`

	_, err = r.db.ExecContext(ctx, query, idUser)
	if err != nil {
		return
	}

	return
}

func (r *repository) DeleteAuthenticationRefreshToken(ctx context.Context, refreshToken string) (err error) {
	query := `
		DELETE 
		FROM 
			authentications
		WHERE refresh_token=$1
	`

	_, err = r.db.ExecContext(ctx, query, refreshToken)
	if err != nil {
		return
	}

	return
}

func (r *repository) UpdateProfileOnboarding(ctx context.Context, email string, model entity.UserEntity) (err error) {
	query := `
		UPDATE users 
			SET username=$1, picture=$2, bio=$3
		WHERE email=$4
	`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return
	}

	_, err = stmt.Stmt.ExecContext(ctx, model.Username, model.Picture, model.Bio, email)
	if err != nil {
		return
	}

	return
}
