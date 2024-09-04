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
	RegenerateAccessToken(ctx context.Context, req request.RegenerateAccessTokenRequestPayload) (accessToken string, err error)
	LogoutUser(ctx context.Context, req request.LogoutRequestPayload) (err error)
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

	err = u.repo.VerifyAvailableUsername(ctx, userEntity.Username)
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
	userEntity := entity.NewFromLoginRequest(req)

	err = userEntity.ValidateLogin()
	if err != nil {
		return
	}

	userEntity, err = u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return
	}

	err = utility.VerifyPasswordFormPlainn(userEntity.Password, req.Password)
	if err != nil {
		err = errorpkg.ErrPasswordNotMatch
		return
	}

	accessToken, err = utility.GenerateToken(userEntity.Username, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return
	}

	refreshToken, err = utility.GenerateToken(userEntity.Username, string(userEntity.Role), config.Cfg.AuthConfig.RefreshTokenKey, config.Cfg.AuthConfig.RefreshTokenExpiration)
	if err != nil {
		return
	}

	authEntity := entity.NewFromLoginRequestToAuth(userEntity.UserId, refreshToken)
	authEntity.GetRefreshTokenExpiration(config.Cfg.AuthConfig.RefreshTokenExpiration)

	err = u.repo.DeleteAuthenticationById(ctx, userEntity.UserId)
	if err != nil {
		return
	}

	err = u.repo.AddAuthentication(ctx, authEntity)
	if err != nil {
		return
	}

	return
}

func (u *usecase) RegenerateAccessToken(ctx context.Context, req request.RegenerateAccessTokenRequestPayload) (accessToken string, err error) {
	userEntity := entity.UserEntity{}

	username, _, err := utility.ValidateToken(req.RefreshToken, config.Cfg.AuthConfig.RefreshTokenKey)
	if err != nil {
		return
	}

	err = u.repo.VerifyAvailableToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}

	userEntity, err = u.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return
	}

	accessToken, err = utility.GenerateToken(userEntity.Username, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return
	}

	return
}

func (u *usecase) LogoutUser(ctx context.Context, req request.LogoutRequestPayload) (err error) {
	_, _, err = utility.ValidateToken(req.RefreshToken, config.Cfg.AuthConfig.RefreshTokenKey)
	if err != nil {
		return
	}

	err = u.repo.DeleteAuthenticationRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}

	return
}
