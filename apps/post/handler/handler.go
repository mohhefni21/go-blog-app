package handler

import (
	"errors"
	"mohhefni/go-blog-app/apps/post/request"
	"mohhefni/go-blog-app/apps/post/response"
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

	coverPicture, err := c.FormFile("cover")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			return responsepkg.NewResponse(
				responsepkg.WithStatus(errorpkg.ErrorBadRequest),
			).Send(c)
		}

		coverPicture = nil
	}

	err = h.ucs.UploadCover(c.Request().Context(), coverPicture, idPost)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
	).Send(c)
}

func (h *handler) GetPosts(c echo.Context) error {
	req := request.GetPostsRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	posts, err := h.ucs.GetDataPosts(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	postsList := response.NewListPostsResponse(posts)

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithData(map[string]interface{}{
			"posts": postsList,
		}),
		responsepkg.WithQuery(req),
	).Send(c)
}
