package auth

import (
	"context"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	user "github.com/sopial42/cleanic/internal/domains/user"
)

type Service interface {
	Signup(ctx context.Context, newUser user.User) (user.User, error)
	Login(ctx context.Context, loginUser user.User) (utils.RefreshToken, utils.AccessToken, error)
	Logout(ctx context.Context, userID user.ID) error
	Refresh(ctx context.Context, signedToken utils.SignedRefreshToken) (utils.RefreshToken, utils.AccessToken, error)
}

type Persistence interface {
	// StoreRefreshToken create or rotate the current token associated to a userID
	StoreRefreshTokenClaims(ctx context.Context, claims utils.RefreshTokenClaims) error
	GetRefreshTokenClaimsByUserID(ctx context.Context, userID user.ID) (utils.RefreshTokenClaims, error)
	DeleteRefreshTokenClaims(ctx context.Context, userID user.ID) error
}

type UserClient interface {
	Create(ctx context.Context, newUser user.User) (user.User, error)
	GetUserByEmail(ctx context.Context, email user.Email) (user.User, error)
	GetUserByID(ctx context.Context, userID user.ID) (user.User, error)
}
