package post

import (
	"mohhefni/go-blog-app/apps/post/handler"
	"mohhefni/go-blog-app/apps/post/repository"
	"mohhefni/go-blog-app/apps/post/usecase"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)

	v1 := e.Group("api/v1/posts")

	v1.POST("", handler.PostAddPost)
	v1.PUT("/cover", handler.PutUpdateCover)
}
