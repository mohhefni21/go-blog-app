package tag

import (
	"mohhefni/go-blog-app/apps/tag/handler"
	"mohhefni/go-blog-app/apps/tag/repository"
	"mohhefni/go-blog-app/apps/tag/usecase"
	"mohhefni/go-blog-app/infra/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)

	v1 := e.Group("api/v1/tags")

	v1.GET("", handler.GetTagName, middleware.ChechAuth)
}
