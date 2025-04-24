package auth

import (
	"context"
	"fmt"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
	user "github.com/sopial42/cleanic/internal/domains/user"
	"golang.org/x/crypto/bcrypt"
)

type authSVC struct {
	uClient   UserClient
	jwtConfig config.JWTConfig
}

func NewAuthService(uClient UserClient, jwtConfig config.JWTConfig) Service {
	return &authSVC{
		uClient:   uClient,
		jwtConfig: jwtConfig,
	}
}

func (a *authSVC) Register(ctx context.Context, newUser user.User) (user.User, error) {
	userCreated, err := a.uClient.Create(ctx, newUser)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to create a user: %w", err)
	}

	return userCreated, nil
}

func (a *authSVC) Login(ctx context.Context, loginUser user.User) (user.User, utils.SignedJWT, error) {
	userFound, err := a.uClient.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return user.User{}, utils.SignedJWT{}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		return user.User{}, utils.SignedJWT{}, fmt.Errorf("invalid password: %w", err)
	}

	jwt, err := utils.GenerateJWT(userFound.ID, userFound.Roles, a.jwtConfig.GetSecret())
	if err != nil {
		return user.User{}, utils.SignedJWT{}, fmt.Errorf("unable to generate token: %w", err)
	}

	return userFound, jwt, nil
}
