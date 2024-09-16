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
	publicId := c.Get("public_id").(string)

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idPost, err := h.ucs.CreatePost(c.Request().Context(), req, publicId)
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
		responsepkg.WithMessage("Cover post berhasil diubah"),
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

	posts, err := h.ucs.GetListPosts(c.Request().Context(), req)
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
		responsepkg.WithQuery(req.DefaultValuePagination()),
	).Send(c)
}

func (h *handler) GetDetailPost(c echo.Context) error {
	slug := c.Param("slug")

	post, comment, err := h.ucs.GetDetailPost(c.Request().Context(), slug)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	detailPost := response.NewDetailPostResponse(post, comment)

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithData(map[string]interface{}{
			"post": detailPost,
		}),
	).Send(c)
}

func (h *handler) GetPostsByUsername(c echo.Context) error {
	username := c.Param("username")
	req := request.GetPostsRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	posts, err := h.ucs.GetListPostsByUsername(c.Request().Context(), req, username)
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
		responsepkg.WithQuery(req.DefaultValuePagination()),
	).Send(c)
}

func (h *handler) GetPostsByUserLogin(c echo.Context) error {
	publicId := c.Get("public_id").(string)

	posts, err := h.ucs.GetListPostsByUserLogin(c.Request().Context(), publicId)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	postsList := response.NewListPostsByUserLoginResponse(posts)

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithData(map[string]interface{}{
			"posts": postsList,
		}),
	).Send(c)
}

func (h *handler) DeletePost(c echo.Context) error {
	slug := c.Param("slug")

	err := h.ucs.DeletePost(c.Request().Context(), slug)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithMessage("Post berhasil dihapus"),
	).Send(c)
}

func (h *handler) PutUpdatePost(c echo.Context) error {
	idPost := c.Param("idPost")
	req := request.UpdatePostRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	err = h.ucs.UpdatePost(c.Request().Context(), req, idPost)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithMessage("Post berhasil diubah"),
	).Send(c)
}

func (h *handler) PostUploadContentImage(c echo.Context) error {
	idPost := c.Param("idPost")

	contentImage, err := c.FormFile("upload")
	if err != nil {
		if !errors.Is(err, http.ErrMissingFile) {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": map[string]string{
					"message": err.Error(),
				},
			})
		}
	}

	url, err := h.ucs.UpdateImageContent(c.Request().Context(), idPost, contentImage)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": map[string]string{
				"message": err.Error(),
			},
		})
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithUrlCKEditor(url),
	).Send(c)
}
