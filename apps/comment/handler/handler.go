package handler

import (
	"mohhefni/go-blog-app/apps/comment/request"
	"mohhefni/go-blog-app/apps/comment/usecase"
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

func (h *handler) PostAddComment(c echo.Context) error {
	req := request.AddCommentPayload{}
	publicId := c.Get("public_id").(string)

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idComment, err := h.ucs.CreateComment(c.Request().Context(), req, publicId)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"idComment": idComment,
		}),
	).Send(c)
}

func (h *handler) PutUpdateComment(c echo.Context) error {
	idComment := c.Param("idComment")
	req := request.UpdateCommentPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	err = h.ucs.UpdateComment(c.Request().Context(), req, idComment)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithMessage("Comment berhasil diubah"),
	).Send(c)
}

func (h *handler) DeleteComment(c echo.Context) error {
	idComment := c.Param("idComment")

	err := h.ucs.DeleteComment(c.Request().Context(), idComment)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithMessage("Comment berhasil dihapus"),
	).Send(c)
}
