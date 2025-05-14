package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	authMiddleware "github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
	user "github.com/sopial42/cleanic/internal/domains/user"
	authSVC "github.com/sopial42/cleanic/internal/services/auth"
)

type authHandler struct {
	authService   authSVC.Service
	cookiesConfig config.CookieStoreConfig
}

type AccessTokenResponse struct {
	Token            utils.SignedAccessToken `json:"access_token"`
	Type             string                  `json:"token_type"`
	ExpiresInSeconds int64                   `json:"expires_in"`
}

// UserUpdateInput is used to parse input for multiples reasons
// do not handle roles as it is implemented on a dedicated safe route
// avoid the used default tag `json:"-"` in case user needs a pwd update
type UserUpdateInput struct {
	ID       user.ID       `json:"id"`
	Email    user.Email    `json:"email"`
	Password user.Password `json:"password"`
}

func SetHandler(e *echo.Echo, config config.CookieStoreConfig, service authSVC.Service, refreshMiddleware authMiddleware.AuthRefreshMiddleware) {
	u := &authHandler{
		service,
		config,
	}
	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/auth/signup", u.register)
		apiV1.POST("/auth/login", u.login)
		apiV1.POST("/auth/refresh", u.refresh, refreshMiddleware.RequireRefreshToken())
		apiV1.POST("/auth/logout", u.logout, refreshMiddleware.RequireRefreshToken())
	}
}

func (a *authHandler) register(context echo.Context) error {
	ctx := context.Request().Context()
	newUserInput := new(UserUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	newUser := user.User{
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	userCreated, err := a.authService.Signup(ctx, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusCreated, userCreated)
}

func (a *authHandler) login(context echo.Context) error {
	ctx := context.Request().Context()
	newUserInput := new(UserUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return context.JSON(http.StatusBadRequest, fmt.Errorf("unable to parse input: %w", err))
	}

	newUser := user.User{
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	refreshToken, accessToken, err := a.authService.Login(ctx, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to login: %w", err))
	}

	sess, err := session.Get(authMiddleware.SessionName, context)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to get session: %w", err))
	}

	sess.Options = &sessions.Options{
		Domain:   a.cookiesConfig.Domain,
		HttpOnly: true,
		MaxAge:   a.cookiesConfig.MaxAgeSeconds,
		Path:     a.cookiesConfig.Domain,
		SameSite: http.SameSite(a.cookiesConfig.SameSite),
	}

	sess.Values[authMiddleware.RefreshTokenCookieName] = string(refreshToken.SignedToken)
	if err := sess.Save(context.Request(), context.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to save session: %w", err))
	}

	return context.JSON(http.StatusOK, AccessTokenResponse{
		Token:            accessToken.SignedToken,
		Type:             accessToken.Type,
		ExpiresInSeconds: int64(accessToken.ExpirationDurationMin.Seconds()),
	})
}

func (a *authHandler) refresh(context echo.Context) error {
	ctx := context.Request().Context()
	sess, err := session.Get(authMiddleware.SessionName, context)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to get session: %w", err))
	}

	currentSignedRefreshTokenValue := sess.Values[authMiddleware.RefreshTokenCookieName]
	if currentSignedRefreshTokenValue == nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("handler unable to get refreshToken from session, not found / nil token"))
	}

	currentSignedRefreshToken, ok := currentSignedRefreshTokenValue.(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to parse refreshToken from session: %w", err))
	}

	// Refresh token in the database
	refreshToken, accessToken, err := a.authService.Refresh(ctx, utils.SignedRefreshToken(currentSignedRefreshToken))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to refresh tokens: %w", err))
	}

	sess.Options = &sessions.Options{
		Domain:   a.cookiesConfig.Domain,
		HttpOnly: true,
		MaxAge:   a.cookiesConfig.MaxAgeSeconds,
		Path:     a.cookiesConfig.Domain,
		SameSite: http.SameSite(a.cookiesConfig.SameSite),
		Secure:   a.cookiesConfig.Secure,
	}

	sess.Values[authMiddleware.RefreshTokenCookieName] = string(refreshToken.SignedToken)
	if err := sess.Save(context.Request(), context.Response()); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to save refresh tokens: %w", err))
	}

	return context.JSON(http.StatusOK, AccessTokenResponse{
		Token:            accessToken.SignedToken,
		Type:             accessToken.Type,
		ExpiresInSeconds: int64(accessToken.ExpirationDurationMin.Seconds()),
	})
}

// Refresh handler is protected by the jwtRefresh middleware
func (a *authHandler) logout(context echo.Context) error {
	ctx := context.Request().Context()
	reqUserID, err := contextUtils.GetUserIDFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to authenticate user: %w", err))
	}

	if err := a.authService.Logout(ctx, reqUserID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to logout user: %w", err))
	}

	sess, err := session.Get(authMiddleware.SessionName, context)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		MaxAge: -1,
	}

	if err := sess.Save(context.Request(), context.Response()); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, "logged out")
}
