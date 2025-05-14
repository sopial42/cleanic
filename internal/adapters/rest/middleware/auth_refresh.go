package middleware

import (
	"fmt"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	jwtUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
)

const RefreshTokenCookieName = "refresh_token"
const SessionName = "session"

type AuthRefreshMiddleware struct {
	secret                 []byte
	TokenExpirationMinutes int
}

func NewAuthRefreshMiddleware(config config.RefreshTokenConfig) AuthRefreshMiddleware {
	return AuthRefreshMiddleware{
		secret:                 config.GetSecret(),
		TokenExpirationMinutes: config.TokenExpirationMinutes,
	}
}

func (a *AuthRefreshMiddleware) RequireRefreshToken() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			sess, err := session.Get(SessionName, c)
			if err != nil {
				return fmt.Errorf("unable to get refreshToken from session: %w", err)
			}

			token := sess.Values[RefreshTokenCookieName]
			if token == nil {
				return fmt.Errorf("middlware unable to get refreshToken from session, not found / nil token")
			}

			tokenStr, ok := token.(string)
			if !ok {
				return fmt.Errorf("unable to parse signed refresh token from session: %w", err)
			}

			claims, err := jwtUtils.ParseRefreshClaims(jwtUtils.SignedRefreshToken(tokenStr), a.secret)
			if err != nil {
				return fmt.Errorf("unable to parse refreshToken: %w", err)
			}

			if claims.ExpiresAt < time.Now().Unix() {
				return fmt.Errorf("refreshToken expired, need login")
			}

			contextUtils.SetUserIDToContext(c, claims.Subject)
			return next(c)
		}
	}
}
