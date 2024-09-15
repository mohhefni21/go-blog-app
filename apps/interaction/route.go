package interaction

import (
	"mohhefni/go-blog-app/apps/interaction/handler"
	"mohhefni/go-blog-app/apps/interaction/repository"
	"mohhefni/go-blog-app/apps/interaction/usecase"
	"mohhefni/go-blog-app/infra/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)

	v1 := e.Group("api/v1/interactions")

	v1.POST("/like", handler.PostAddInteractionLike, middleware.ChechAuth)
	v1.POST("/share", handler.PostAddInteractionShare, middleware.ChechAuth)
	v1.POST("/bookmark", handler.PostAddInteractionBookmark, middleware.ChechAuth)
	v1.DELETE("/:idInteraction", handler.DeleteInteraction, middleware.ChechAuth)
}
