package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	jwtUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
	"github.com/sopial42/cleanic/internal/domains/user"
)

type AuthRefreshMiddleware struct {
	secret                 []byte
	TokenExpirationMinutes int
}

func NewAuthRefreshMiddleware(config config.AccessTokenConfig) AuthRefreshMiddleware {
	return AuthRefreshMiddleware{
		secret:                 config.GetSecret(),
		TokenExpirationMinutes: config.TokenExpirationMinutes,
	}
}

func (a *AuthRefreshMiddleware) RequireRoles(requiredRoles user.Roles) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			token, err := jwtUtils.ParseBearerHeader(header)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse authorization header: %w", err))
			}

			claims, err := jwtUtils.ParseRefreshClaims(token, a.secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse auth token: %w", err))
			}

			contextUtils.SetUserIDToContext(c, claims.Subject)
			return next(c)
		}
	}
}
