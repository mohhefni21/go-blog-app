package handler

import (
	"mohhefni/go-blog-app/apps/post/request"
	"mohhefni/go-blog-app/apps/post/usecase"
	"mohhefni/go-blog-app/infra/errorpkg"
	"mohhefni/go-blog-app/infra/responsepkg"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	ucs usecase.Usecase
}

func NewHandler(usecase usecase.Usecase) handler {
	return handler{
		ucs: usecase,
	}
}

func (h *handler) PostAddPost(c echo.Context) error {
	req := request.AddPostRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idPost, err := h.ucs.CreatePost(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"idPost": idPost,
		}),
	).Send(c)
}

func (h *handler) PutUpdateCover(c echo.Context) error {
	idPost := c.FormValue("idPost")

	cover, err := c.FormFile("cover")
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	err = h.ucs.UploadCover(c.Request().Context(), cover, idPost)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
	).Send(c)
}
