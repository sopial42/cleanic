package context

import (
	"context"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/sopial42/cleanic/internal/domains/user"
)

var (
	UserIDKey    = ContextKey("user_id")
	UserRolesKey = ContextKey("roles")
)

type ContextKey string

func GetUserIDFromContext(ctx context.Context) (user.ID, error) {
	raw := ctx.Value(UserIDKey)
	if raw == nil {
		return 0, fmt.Errorf("user id not found in context")
	}

	id, ok := raw.(user.ID)
	if !ok {
		return 0, fmt.Errorf("user id is not a valid format (int64)")
	}

	return id, nil
}

func SetUserIDToContext(ctxEcho echo.Context, userID user.ID) {
	ctx := ctxEcho.Request().Context()
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctxEcho.SetRequest(ctxEcho.Request().WithContext(ctx))
}

func SetUserIDAndRolesToContext(ctxEcho echo.Context, userID user.ID, userRoles user.Roles) {
	ctx := ctxEcho.Request().Context()
	ctx = context.WithValue(ctx, UserIDKey, userID)
	ctx = context.WithValue(ctx, UserRolesKey, userRoles)
	ctxEcho.SetRequest(ctxEcho.Request().WithContext(ctx))
}
