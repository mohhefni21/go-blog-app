package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"mohhefni/go-blog-app/apps/auth/entity"
	"mohhefni/go-blog-app/apps/auth/repository"
	"mohhefni/go-blog-app/apps/auth/request"
	"mohhefni/go-blog-app/infra/errorpkg"
	"mohhefni/go-blog-app/internal/config"
	"mohhefni/go-blog-app/utility"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type Usecase interface {
	RegisterUser(ctx context.Context, req request.RegisterRequestPayload) (email string, err error)
	LoginUser(ctx context.Context, req request.LoginRequestPayload) (accessToken string, refreshToken string, err error)
	RegenerateAccessToken(ctx context.Context, req request.RegenerateAccessTokenRequestPayload) (accessToken string, err error)
	LogoutUser(ctx context.Context, req request.LogoutRequestPayload) (err error)
	AuthWithGoogle(ctx context.Context) (redirectUrl string, err error)
	AuthWithGoogleCallback(ctx context.Context, req request.OauthGoogleRequestPayload) (userEntity entity.UserEntity, accessToken string, refreshToken string, err error)
	UpdateProfileOnboarding(ctx context.Context, req request.UpdateProfileOnboardingRequestPayload) (accessToken string, refreshToken string, err error)
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

	accessToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return
	}

	refreshToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.RefreshTokenKey, config.Cfg.AuthConfig.RefreshTokenExpiration)
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

	PublicId, _, err := utility.ValidateToken(req.RefreshToken, config.Cfg.AuthConfig.RefreshTokenKey)
	if err != nil {
		return
	}

	err = u.repo.VerifyAvailableToken(ctx, req.RefreshToken)
	if err != nil {
		return
	}

	uuidPublicId, err := utility.ParseUUID(PublicId)
	if err != nil {
		return
	}

	userEntity, err = u.repo.GetUserByPublicId(ctx, uuidPublicId)
	if err != nil {
		return
	}

	accessToken, err = utility.GenerateToken(uuidPublicId, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
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

func (u *usecase) AuthWithGoogle(ctx context.Context) (redirectUrl string, err error) {
	googleConfig := utility.ConfigGoogle(config.Cfg.OAuthConfig)

	redirectUrl = googleConfig.AuthCodeURL(config.Cfg.OAuthConfig.GoogleStateToken, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "select_account"))

	return
}

func (u *usecase) AuthWithGoogleCallback(ctx context.Context, req request.OauthGoogleRequestPayload) (userEntity entity.UserEntity, accessToken string, refreshToken string, err error) {
	googleConfig := utility.ConfigGoogle(config.Cfg.OAuthConfig)

	token, err := googleConfig.Exchange(ctx, req.Code)
	if err != nil {
		return
	}

	client := googleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	userPayload := entity.OauthGoogleUserPayload{}
	err = json.NewDecoder(resp.Body).Decode(&userPayload)
	if err != nil {
		return
	}

	userEntity, err = u.repo.GetUserByEmail(ctx, userPayload.Email)
	if err != nil {
		if err == errorpkg.ErrorNotFound {
			newUser := entity.UserEntity{
				PublicId:  uuid.New(),
				Username:  userPayload.Name,
				Fullname:  userPayload.Name,
				Email:     userPayload.Email,
				Role:      entity.ROLE_USER,
				Picture:   sql.NullString{String: userPayload.Picture, Valid: true},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			newUser.GenerateUsernameOauth(userPayload.Id)

			_, err = u.repo.AddUser(ctx, newUser)
			if err != nil {
				return
			}

			return newUser, "", "", nil
		}
		return
	}

	accessToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return entity.UserEntity{}, "", "", err
	}

	refreshToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.RefreshTokenKey, config.Cfg.AuthConfig.RefreshTokenExpiration)
	if err != nil {
		return entity.UserEntity{}, "", "", err
	}

	err = u.repo.DeleteAuthenticationById(ctx, userEntity.UserId)
	if err != nil {
		return entity.UserEntity{}, "", "", err
	}

	authEntity := entity.NewFromLoginRequestToAuth(userEntity.UserId, refreshToken)
	authEntity.GetRefreshTokenExpiration(config.Cfg.AuthConfig.RefreshTokenExpiration)

	err = u.repo.AddAuthentication(ctx, authEntity)
	if err != nil {
		return entity.UserEntity{}, "", "", err
	}

	return entity.UserEntity{}, accessToken, refreshToken, nil
}

func (u *usecase) UpdateProfileOnboarding(ctx context.Context, req request.UpdateProfileOnboardingRequestPayload) (accessToken string, refreshToken string, err error) {
	var fileName string
	if req.Picture != nil {
		filePicture := req.Picture

		fileName, err = utility.UploadFile(filePicture, "static/profile")
		if err != nil {
			return
		}
	}

	userEntity, err := u.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return
	}

	if req.Picture == nil {
		fileName = userEntity.Picture.String
	}

	userEntity.UsernameValidate()

	err = u.repo.VerifyAvailableUsernameByEmail(ctx, userEntity.Username, req.Email)
	if err != nil {
		return
	}

	userEntity.Username = req.Username
	userEntity.Picture.String = fileName
	userEntity.Bio.String = req.Bio

	err = u.repo.UpdateProfileOnboarding(ctx, req.Email, userEntity)
	if err != nil {
		return
	}

	accessToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.AccessTokenKey, config.Cfg.AuthConfig.AccessTokenExpiration)
	if err != nil {
		return
	}

	refreshToken, err = utility.GenerateToken(userEntity.PublicId, string(userEntity.Role), config.Cfg.AuthConfig.RefreshTokenKey, config.Cfg.AuthConfig.RefreshTokenExpiration)
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
