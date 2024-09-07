package handler

import (
	"mohhefni/go-blog-app/apps/auth/request"
	"mohhefni/go-blog-app/apps/auth/usecase"
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

func (h *handler) PostRegisterUser(c echo.Context) error {
	req := request.RegisterRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	email, err := h.ucs.RegisterUser(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"emailUser": email,
		}),
	).Send(c)
}

func (h *handler) PostLoginUser(c echo.Context) error {
	req := request.LoginRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	accessToken, refreshToken, err := h.ucs.LoginUser(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}),
	).Send(c)
}

func (h *handler) PostRegenerateAccessToken(c echo.Context) error {
	req := request.RegenerateAccessTokenRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	accessToken, err := h.ucs.RegenerateAccessToken(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusCreated),
		responsepkg.WithData(map[string]interface{}{
			"accessToken": accessToken,
		}),
	).Send(c)
}

func (h *handler) DeleteLogoutUser(c echo.Context) error {
	req := request.LogoutRequestPayload{}

	err := c.Bind(&req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	err = h.ucs.LogoutUser(c.Request().Context(), req)
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithMessage("Logout berhasil"),
	).Send(c)
}

func (h *handler) GetGoogleLogin(c echo.Context) error {
	redirectUrl, err := h.ucs.AuthWithGoogle(c.Request().Context())
	if err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	return c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (h *handler) GetGoogleLoginCallback(c echo.Context) error {
	req := request.OauthGoogleRequestPayload{}

	// Bind the request payload
	if err := c.Bind(&req); err != nil {
		return responsepkg.NewResponse(
			responsepkg.WithStatus(errorpkg.ErrorBadRequest),
		).Send(c)
	}

	// Call usecase method to handle the logic
	accessToken, refreshToken, err := h.ucs.AuthWithGoogleCallback(c.Request().Context(), req)
	if err != nil {
		// Handle the error and return a proper response
		return responsepkg.NewResponse(
			responsepkg.WithStatus(err),
		).Send(c)
	}

	// If it's a new user, return a message to redirect to profile update
	if accessToken == "" && refreshToken == "" {
		return responsepkg.NewResponse(
			responsepkg.WithHttpCode(http.StatusCreated),
			responsepkg.WithMessage("User registered successfully, redirect to profile update"),
		).Send(c)
	}

	// Return the tokens for existing user
	return responsepkg.NewResponse(
		responsepkg.WithHttpCode(http.StatusOK),
		responsepkg.WithData(map[string]interface{}{
			"accessToken":  accessToken,
			"refreshToken": refreshToken,
		}),
	).Send(c)
}
