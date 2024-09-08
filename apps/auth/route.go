package auth

import (
	"mohhefni/go-blog-app/apps/auth/handler"
	"mohhefni/go-blog-app/apps/auth/repository"
	"mohhefni/go-blog-app/apps/auth/usecase"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)

	v1 := e.Group("api/v1/auth")

	v1.POST("/register", handler.PostRegisterUser)
	v1.POST("/login", handler.PostLoginUser)
	v1.POST("/refresh-token", handler.PostRegenerateAccessToken)
	v1.DELETE("/logout", handler.DeleteLogoutUser)
	v1.GET("/google", handler.GetGoogleLogin)
	v1.GET("/google/callback", handler.GetGoogleLoginCallback)
	v1.PUT("/onboarding", handler.PutUpdateProfilOnboarding)
}
