package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"

	authMiddleware "github.com/sopial42/cleanic/internal/adapters/rest/middleware"
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
		apiV1.POST("/auth/refresh", u.refresh)
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

	return context.JSON(http.StatusOK, loginResponse)
}


func (a *authHandler) refresh(context echo.Context) error {
	// ctx := context.Request().Context()
	// newUserInput := new(UserUpdateInput)
	// if err := context.Bind(newUserInput); err != nil {
	// 	return echo.NewHTTPError(http.StatusBadRequest, err)
	// }

	// newUser := user.User{
	// 	Email:    newUserInput.Email,
	// 	Password: newUserInput.Password,
	// }

	// userCreated, err := a.authService.Signup(ctx, newUser)
	// if err != nil {
	// 	return echo.NewHTTPError(http.StatusInternalServerError, err)
	// }

	return context.JSON(http.StatusCreated, "userCreated")
}
