package middleware

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sopial42/cleanic/internal/domains/user"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

func RequireRoles(requiredRoles user.Roles) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			token, err := userSVC.NewJWTFromAuthHeaderString(auth)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse authorization header: %w", err))
			}

			userID, userRoles, err := token.ParseClaims()
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, fmt.Errorf("unable to parse auth token: %w", err))
			}

			if err := userRoles.ValidateRequiredRoles(requiredRoles); err != nil {
				return echo.NewHTTPError(http.StatusForbidden, fmt.Errorf("unauthorized resource: %w", err))
			}

			c.Set("roles", userRoles.String())
			c.Set("user_id", userID)
			return next(c)
		}
	}
}
