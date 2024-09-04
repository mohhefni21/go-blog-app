package middleware

import (
	"mohhefni/go-blog-app/infra/errorpkg"
	"mohhefni/go-blog-app/infra/responsepkg"
	"mohhefni/go-blog-app/internal/config"
	"mohhefni/go-blog-app/utility"
	"strings"

	"github.com/labstack/echo/v4"
)

func ChechAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authorization := c.Request().Header.Get("Authorization")
		if authorization == "" {
			return responsepkg.NewResponse(
				responsepkg.WithStatus(errorpkg.ErrUnauthorized),
			).Send(c)
		}

		splitToken := strings.Split(authorization, "Bearer ")
		if len(splitToken) != 2 {
			return responsepkg.NewResponse(
				responsepkg.WithStatus(errorpkg.ErrUnauthorized),
			).Send(c)
		}

		token := splitToken[1]

		username, role, err := utility.ValidateToken(token, config.Cfg.AuthConfig.AccessTokenKey)
		if err != nil {
			return responsepkg.NewResponse(
				responsepkg.WithStatus(errorpkg.ErrUnauthorized),
			).Send(c)
		}

		c.Set("username", username)
		c.Set("role", role)

		return next(c)
	}
}
