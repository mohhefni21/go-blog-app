package handler

import (
	"mohhefni/go-blog-app/apps/interaction/request"
	"mohhefni/go-blog-app/apps/interaction/usecase"
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

func (h *handler) PostAddInteractionLike(c echo.Context) error {
	req := request.AddInteractionRequestPayload{}
	publicId := c.Get("public_id").(string)

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idInteraction, err := h.ucs.CreateInteractionLike(c.Request().Context(), req, publicId)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"idInteraction": idInteraction,
		}),
	).Send(c)
}

func (h *handler) PostAddInteractionShare(c echo.Context) error {
	req := request.AddInteractionRequestPayload{}
	publicId := c.Get("public_id").(string)

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idInteraction, err := h.ucs.CreateInteractionShare(c.Request().Context(), req, publicId)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"idInteraction": idInteraction,
		}),
	).Send(c)
}

func (h *handler) PostAddInteractionBookmark(c echo.Context) error {
	req := request.AddInteractionRequestPayload{}
	publicId := c.Get("public_id").(string)

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	idInteraction, err := h.ucs.CreateInteractionBookmark(c.Request().Context(), req, publicId)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"idInteraction": idInteraction,
		}),
	).Send(c)
}

func (h *handler) DeleteInteraction(c echo.Context) error {
	idInteraction := c.Param("idInteraction")

	err := h.ucs.DeleteInteraction(c.Request().Context(), idInteraction)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithMessage("Interaction berhasil dihapus"),
	).Send(c)
}
