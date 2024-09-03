package usecase

import (
	"context"
	"mohhefni/go-blog-app/apps/auth/entity"
	"mohhefni/go-blog-app/apps/auth/repository"
	"mohhefni/go-blog-app/apps/auth/request"
	"mohhefni/go-blog-app/infra/errorpkg"
	"mohhefni/go-blog-app/internal/config"
	"mohhefni/go-blog-app/utility"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req request.RegisterRequestPayload) (email string, err error)
	LoginUser(ctx context.Context, req request.LoginRequestPayload) (accessToken string, refreshToken string, err error)
}

type usecase struct {
	repo repository.Repository
}

func NewUsecase(repo repository.Repository) Usecase {
	return &usecase{
		repo: repo,
	}
}

func (u *usecase) RegisterUser(ctx context.Context, req request.RegisterRequestPayload) (email string, err error) {
	userEntity := entity.NewFromRegisterRequest(req)

	err = userEntity.RegisterValidate()
	if err != nil {
		return
	}

	err = u.repo.VerifyAvailableEmail(ctx, userEntity.Email)
	if err != nil {
		return
	}

	err = u.repo.VerifyAvailableUsername(ctx, userEntity.Email)
	if err != nil {
		return
	}

	userEntity.Password, err = utility.EncryptPassword(req.Password)
	if err != nil {
		return
	}

	email, err = u.repo.AddUser(ctx, userEntity)

	return
}

func (u *usecase) LoginUser(ctx context.Context, req request.LoginRequestPayload) (accessToken string, refreshToken string, err error) {
	authEntity := entity.NewFromLoginRequest(req)

	err = authEntity.ValidateLogin()
	if err != nil {
		return
	}

	authEntity, err = u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return
	}

	err = utility.VerifyPasswordFormPlainn(authEntity.Password, req.Password)
	if err != nil {
		err = errorpkg.ErrPasswordNotMatch
		return
	}

	accessToken, err = utility.GenerateToken(authEntity.Username, string(authEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return
	}

	accessToken, err = utility.GenerateToken(authEntity.Username, string(authEntity.Role), config.Cfg.AuthConfig.RefreshTokenKey, config.Cfg.AuthConfig.RefreshTokenExpiration)
	if err != nil {
		return
	}

	return
}
