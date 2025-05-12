package auth

import (
	"context"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	user "github.com/sopial42/cleanic/internal/domains/user"
)

type Service interface {
	Signup(ctx context.Context, newUser user.User) (user.User, error)
	Login(ctx context.Context, loginUser user.User) (LoginResponse, error)
	Refresh(ctx context.Context, reqUserID user.ID) (LoginResponse, error)
}

type Persistence interface {
	RegisterRefreshToken(ctx context.Context, token utils.RefreshToken) (utils.RefreshToken, error)
}

type UserClient interface {
	Create(ctx context.Context, newUser user.User) (user.User, error)
	GetUserByEmail(ctx context.Context, email user.Email) (user.User, error)
	GetUserByID(ctx context.Context, userID user.ID) (user.User, error)
}
