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

type AuthMiddleware struct {
	secret []byte
}

func NewAuthMiddleware(jwtConfig config.JWTConfig) AuthMiddleware {
	return AuthMiddleware{
		secret: jwtConfig.GetSecret(),
	}
}

func (a *AuthMiddleware) RequireRoles(requiredRoles user.Roles) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			token, err := jwtUtils.NewJWTFromAuthHeaderString(header)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse authorization header: %w", err))
			}

			userID, userRoles, err := token.ParseClaims(a.secret)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse auth token: %w", err))
			}

			if err := userRoles.ValidateRequiredRoles(requiredRoles); err != nil {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("unauthorized resource: %w", err))
			}

			contextUtils.SetUserIDAndRolesToContext(c, userID, userRoles)
			return next(c)
		}
	}
}
