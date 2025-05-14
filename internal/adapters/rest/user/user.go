package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	contextUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/context"
	user "github.com/sopial42/cleanic/internal/domains/user"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

type userHandler struct {
	uService userSVC.Service
}

func SetHandler(e *echo.Echo, service userSVC.Service, access middleware.AuthAccessMiddleware) {
	u := &userHandler{
		service,
	}

	requireAdmin := access.RequireRoles(user.Roles{user.RoleAdmin})
	requireDoctor := access.RequireRoles(user.Roles{user.RoleDoctor})
	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/users", u.getUsers, requireAdmin)
		apiV1.GET("/user/:id", u.getUserByID, requireAdmin)
		apiV1.PATCH("/user/roles", u.updateUserRoles, requireAdmin)
		apiV1.PATCH("/user", u.updateUser, requireDoctor)
		apiV1.DELETE("/user/:id", u.deleteUser, requireDoctor)
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
	reqUserID, err := contextUtils.GetUserIDFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to authenticate user: %w", err))
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
	reqUserID, err := contextUtils.GetUserIDFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to authenticate user: %w", err))
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
	reqUserID, err := contextUtils.GetUserIDFromContext(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to authenticate user: %w", err))
	}

	idParam := context.Param("id")
	idToDelete, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = u.uService.DeleteUser(ctx, reqUserID, user.ID(idToDelete))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Errorf("unable to delete user: %w", err))
	}

	return context.NoContent(http.StatusNoContent)
}
