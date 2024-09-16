package post

import (
	"mohhefni/go-blog-app/apps/post/handler"
	"mohhefni/go-blog-app/apps/post/repository"
	"mohhefni/go-blog-app/apps/post/usecase"
	"mohhefni/go-blog-app/infra/middleware"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo, db *sqlx.DB) {
	repo := repository.NewRepository(db)
	usecase := usecase.NewUsecase(repo)
	handler := handler.NewHandler(usecase)

	v1 := e.Group("api/v1/posts")

	v1.Static("/cover", "static/cover")
	v1.Static("/content-image", "static/content-image")
	v1.POST("", handler.PostAddPost, middleware.ChechAuth)
	v1.PUT("/cover", handler.PutUpdateCover, middleware.ChechAuth)
	v1.GET("", handler.GetPosts)
	v1.GET("/detail/:slug", handler.GetDetailPost)
	v1.GET("/user/:username", handler.GetPostsByUsername)
	v1.GET("/admin/dashboard", handler.GetPostsByUserLogin, middleware.ChechAuth)
	v1.DELETE("/:slug", handler.DeletePost, middleware.ChechAuth)
	v1.PUT("/:idPost", handler.PutUpdatePost, middleware.ChechAuth)
	v1.POST("/content-image/:idPost", handler.PostUploadContentImage, middleware.ChechAuth)
}
