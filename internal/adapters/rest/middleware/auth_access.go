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

type AuthAccessMiddleware struct {
	secret                 []byte
	TokenExpirationMinutes int
}

func NewAuthAccessMiddleware(config config.AccessTokenConfig) AuthAccessMiddleware {
	return AuthAccessMiddleware{
		secret:                 config.GetSecret(),
		TokenExpirationMinutes: config.TokenExpirationMinutes,
	}
}

func (a *AuthAccessMiddleware) RequireRoles(requiredRoles user.Roles) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			token, err := jwtUtils.ParseBearerHeader(header)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse authorization header: %w", err))
			}

			claims, err := jwtUtils.ParseAccessClaims(token, a.secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse auth token: %w", err))
			}

			if err := user.ValidateRequiredRoles(requiredRoles, claims.Roles); err != nil {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("unauthorized resource: %w", err))
			}

			contextUtils.SetUserIDAndRolesToContext(c, claims.Subject, claims.Roles)
			return next(c)
		}
	}
}
