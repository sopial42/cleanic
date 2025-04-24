package auth

import (
	"context"

	user "github.com/sopial42/cleanic/internal/domains/user"
	jwtUtils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
)

type Service interface {
	Register(ctx context.Context, newUser user.User) (user.User, error)
	Login(ctx context.Context, loginUser user.User) (user.User, jwtUtils.SignedJWT, error)
}

type UserClient interface {
	Create(ctx context.Context, newUser user.User) (user.User, error)
	GetUserByEmail(ctx context.Context, email user.Email) (user.User, error)
}

