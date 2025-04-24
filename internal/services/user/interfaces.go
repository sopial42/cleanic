package user

import (
	"context"

	user "github.com/sopial42/cleanic/internal/domains/user"
)

type Service interface {
	Create(ctx context.Context, newUser user.User) (user.User, error)
	GetUsers(ctx context.Context) ([]user.User, error)
	GetUserByID(ctx context.Context, id user.ID) (user.User, error)
	GetUserByEmail(ctx context.Context, email user.Email) (user.User, error)
	UpdateUser(ctx context.Context, reqUserID user.ID, updatedUser user.User) (user.User, error)
	UpdateUserRoles(ctx context.Context, reqUserID user.ID, updatedUser user.User) (user.User, error)
	DeleteUser(ctx context.Context, reqUserID user.ID, userIDToDelete user.ID) error
}

type Persistence interface {
	Insert(ctx context.Context, newUser user.User) (user.User, error)
	ListUsers(ctx context.Context) ([]user.User, error)
	GetUserByID(ctx context.Context, id user.ID) (user.User, error)
	GetUserByEmail(ctx context.Context, email user.Email) (user.User, error)
	UpdateUser(ctx context.Context, updatedUser user.User) (user.User, error)
	UpdateUserRoles(ctx context.Context, updatedUserRoles user.User) (user.User, error)
	DeleteUser(ctx context.Context, userIDToDelete user.ID) error
}
