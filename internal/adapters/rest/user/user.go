package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	user "github.com/sopial42/cleanic/internal/domains/user"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

type userHandler struct {
	uService userSVC.Service
}

func SetHandler(e *echo.Echo, service userSVC.Service) {
	u := &userHandler{
		service,
	}
	apiV1 := e.Group("/api/v1")
	{
		apiV1.POST("/user", u.register)
		apiV1.POST("/user/login", u.login)
		apiV1.GET("/users", u.getUsers, middleware.RequireRoles(
			user.Roles{user.RoleAdmin}))
		apiV1.GET("/user/:id", u.getUserByID, middleware.RequireRoles(
			user.Roles{user.RoleAdmin}))
		apiV1.PATCH("/user", u.updateUser, middleware.RequireRoles(
			user.Roles{user.RoleDoctor}))
		apiV1.PATCH("/user/roles", u.updateUserRoles, middleware.RequireRoles(
			user.Roles{user.RoleAdmin}))
		apiV1.DELETE("/user/:id", u.deleteUser, middleware.RequireRoles(
			user.Roles{user.RoleDoctor}))
	}
}

// UserUpdateInput is used to parse input for multiples reasons
// do not handle roles as it is implemented on a dedicated safe route
// avoid the used default tag `json:"-"` in case user needs a pwd update
type UserUpdateInput struct {
	ID       user.ID    `json:"id"`
	Email    user.Email `json:"email"`
	Password string     `json:"password"`
}

func (u *userHandler) register(context echo.Context) error {
	ctx := context.Request().Context()
	newUserInput := new(UserUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	newUser := user.User{
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	userCreated, err := u.uService.Register(ctx, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusCreated, userCreated)
}

func (u *userHandler) login(context echo.Context) error {
	ctx := context.Request().Context()
	newUserInput := new(UserUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return context.JSON(http.StatusBadRequest, err)
	}

	newUser := user.User{
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	userFound, jwt, err := u.uService.Login(ctx, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	context.Response().Header().Set("Authorization", jwt.ToAuthHeaderString())
	return context.JSON(http.StatusOK, userFound)
}

func (u *userHandler) getUsers(context echo.Context) error {
	ctx := context.Request().Context()
	users, err := u.uService.GetUsers(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, users)
}

func (u *userHandler) getUserByID(context echo.Context) error {
	ctx := context.Request().Context()
	id := context.Param("id")
	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	userFound, err := u.uService.GetUserByID(ctx, user.ID(idInt))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, userFound)
}

func (u *userHandler) updateUser(context echo.Context) error {
	ctx := context.Request().Context()
	ctxUserIDValue := context.Get(string(userSVC.UserIDKey))
	reqUserID, ok := ctxUserIDValue.(user.ID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("unable to authenticate user"))
	}

	newUserInput := new(UserUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to parse user input: %w", err))
	}

	newUser := user.User{
		ID:       newUserInput.ID,
		Email:    newUserInput.Email,
		Password: newUserInput.Password,
	}

	userUpdated, err := u.uService.UpdateUser(ctx, reqUserID, newUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to update user: %w", err))
	}

	return context.JSON(http.StatusOK, userUpdated)
}

// UserRolesUpdateInput only handle roles input for safety
type UserRolesUpdateInput struct {
	ID    user.ID    `json:"id"`
	Roles user.Roles `json:"roles"`
}

// updateUserRoles only handle roles for safety
func (u *userHandler) updateUserRoles(context echo.Context) error {
	ctx := context.Request().Context()
	ctxUserIDValue := context.Get(string(userSVC.UserIDKey))
	reqUserID, ok := ctxUserIDValue.(user.ID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("unable to authenticate user"))
	}

	newUserInput := new(UserRolesUpdateInput)
	if err := context.Bind(newUserInput); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("unable to parse user role input: %w", err))
	}

	newUserRoles := user.User{
		ID:    newUserInput.ID,
		Roles: newUserInput.Roles,
	}

	userUpdated, err := u.uService.UpdateUserRoles(ctx, reqUserID, newUserRoles)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to update user: %w", err))
	}

	return context.JSON(http.StatusOK, userUpdated)
}

func (u *userHandler) deleteUser(context echo.Context) error {
	ctx := context.Request().Context()
	idParam := context.Param("id")
	idToDelete, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	ctxUserIDValue := context.Get(string(userSVC.UserIDKey))
	reqUserID, ok := ctxUserIDValue.(user.ID)
	if !ok {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.New("unable to authenticate user"))
	}

	err = u.uService.DeleteUser(ctx, reqUserID, user.ID(idToDelete))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to delete user: %w", err))
	}

	return context.NoContent(http.StatusNoContent)
}
