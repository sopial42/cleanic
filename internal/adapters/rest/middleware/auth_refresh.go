package middleware

import (
	"fmt"
	"log"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	jwtUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
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

func (a *AuthRefreshMiddleware) RequireRefreshToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get("session", c)
			if err != nil {
				return fmt.Errorf("unable to get refreshToken from session: %w", err)
			}

			token := sess.Values["refresh_token"]
			// safeparse token as string
			if token == nil {
				return fmt.Errorf("unable to get refreshToken from session, nil token")
			}

			tokenStr, ok := token.(string)
			if !ok {
				return fmt.Errorf("unable to parse refreshToken from session: %w", err)
			}

			log.Printf("refresh token: %s", tokenStr)
			claims, err := jwtUtils.ParseRefreshClaims(tokenStr, a.secret)
			if err != nil {
				return fmt.Errorf("unable to parse refreshToken: %w", err)
			}

			contextUtils.SetUserIDToContext(c, claims.Subject)
			return next(c)
		}
	}
}
