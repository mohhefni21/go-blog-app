package handler

import (
	"mohhefni/go-blog-app/apps/tag/response"
	"mohhefni/go-blog-app/apps/tag/usecase"
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

func (h *handler) GetTagName(c echo.Context) error {
	tagSearch := c.QueryParam("search")

	tags, err := h.ucs.GetTags(c.Request().Context(), tagSearch)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	tagsList := response.NewTagsListResponse(tags)

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"tags": tagsList,
		}),
	).Send(c)
}
