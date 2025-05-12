package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"

	authMiddleware "github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	user "github.com/sopial42/cleanic/internal/domains/user"
	authSVC "github.com/sopial42/cleanic/internal/services/auth"
)

type authHandler struct {
	authService authSVC.Service
}

func SetHandler(e *echo.Echo, service authSVC.Service, refreshMiddleware authMiddleware.AuthRefreshMiddleware) {
	u := &authHandler{
		service,
	}
	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/auth/signup", u.register)
		apiV1.POST("/auth/login", u.login)
		apiV1.POST("/auth/refresh", u.refresh, refreshMiddleware.RequireRefreshToken())
	}
}

// UserUpdateInput is used to parse input for multiples reasons
// do not handle roles as it is implemented on a dedicated safe route
// avoid the used default tag `json:"-"` in case user needs a pwd update
type UserUpdateInput struct {
	ID       user.ID       `json:"id"`
	Email    user.Email    `json:"email"`
	Password user.Password `json:"password"`
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
		return context.JSON(http.StatusBadRequest, err)
	}

	newUser := user.User{
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	loginResponse, err := a.authService.Login(ctx, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	sess, err := session.Get("session", context)
	if err != nil {
		return err
	}
	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["refresh_token"] = string(loginResponse.RefreshTokenResponse.Token)
	if err := sess.Save(context.Request(), context.Response()); err != nil {
		return err
	}

	return context.JSON(http.StatusOK, loginResponse.AccessTokenResponse)
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh handler is protected by the jwtRefresh middleware
func (a *authHandler) refresh(context echo.Context) error {
	_ = context.Request().Context()
	refreshToken := new(refreshRequest)
	if err := context.Bind(refreshToken); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctx := context.Request().Context()
	reqUserID, err := contextUtils.GetUserIDFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to authenticate user: %w", err))
	}

	newTokens, err := a.authService.Refresh(ctx, reqUserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to refresh tokens: %w", err))
	}

	sess, err := session.Get("session", context)
	if err != nil {
		return err
	}

	sess.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
	}

	sess.Values["refresh_token"] = string(newTokens.RefreshTokenResponse.Token)
	if err := sess.Save(context.Request(), context.Response()); err != nil {
		return err
	}

	return context.JSON(http.StatusCreated, "newTokens")
}
