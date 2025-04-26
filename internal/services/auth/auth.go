package auth

import (
	"context"
	"errors"
	"fmt"

	utils "github.com/sopial42/cleanic/internal/adapters/rest/utils/jwt"
	"github.com/sopial42/cleanic/internal/config"
	user "github.com/sopial42/cleanic/internal/domains/user"
	"golang.org/x/crypto/bcrypt"
)

type authSVC struct {
	uClient     UserClient
	jwtConfig   config.JWTConfig
	persistence Persistence
}

func NewAuthService(uClient UserClient, jwtConfig config.JWTConfig, persistence Persistence) Service {
	return &authSVC{
		uClient:     uClient,
		jwtConfig:   jwtConfig,
		persistence: persistence,
	}
}

func (a *authSVC) Signup(ctx context.Context, newUser user.User) (user.User, error) {
	if !newUser.Password.IsValid() {
		return user.User{}, errors.New("invalid password")
	}

	if !newUser.Email.IsValid() {
		return user.User{}, errors.New("invalid email")
	}

	userCreated, err := a.uClient.Create(ctx, newUser)
	if err != nil {
		return user.User{}, fmt.Errorf("unable to create a user: %w", err)
	}

	return userCreated, nil
}

type LoginResponse struct {
	AccessTokenResponse
	RefreshTokenResponse
}

type RefreshTokenResponse struct {
	Token utils.SignedRefreshToken `json:"refresh_token"`
}

type AccessTokenResponse struct {
	Token            utils.SignedAccessToken `json:"access_token"`
	Type             string                  `json:"token_type"`
	ExpiresInSeconds int64                   `json:"expires_in"`
}

func (a *authSVC) Login(ctx context.Context, loginUser user.User) (LoginResponse, error) {
	userFound, err := a.uClient.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return LoginResponse{}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		return LoginResponse{}, fmt.Errorf("invalid password: %w", err)
	}

	refreshToken, err := utils.NewRefreshToken(
		userFound.ID, a.jwtConfig.RefreshTokenConfig.GetSecret(),
		utils.RefreshTokenTTLDays(a.jwtConfig.AccessTokenConfig.TokenExpirationMinutes),
		utils.RefreshTokenAudience(
			a.jwtConfig.RefreshTokenConfig.Audience))
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unable to generate refresh token: %w", err)
	}

	accessToken, err := utils.NewAccessToken(
		userFound.ID, userFound.Roles, a.jwtConfig.AccessTokenConfig.GetSecret(),
		utils.AccessTokenTTLMin(a.jwtConfig.AccessTokenConfig.TokenExpirationMinutes),
		utils.AccessTokenAudience(
			a.jwtConfig.AccessTokenConfig.Audience,
		),
	)

	if err != nil {
		return LoginResponse{}, fmt.Errorf("unable to generate access token: %w", err)
	}

	_, err = a.persistence.RegisterRefreshToken(ctx, refreshToken)
	if err != nil {
		return LoginResponse{}, fmt.Errorf("unable to register refresh token: %w", err)
	}
	return LoginResponse{
		AccessTokenResponse: AccessTokenResponse{
			Token:            accessToken.Token,
			Type:             accessToken.Type,
			ExpiresInSeconds: int64(accessToken.ExpirationDuration.Seconds()),
		},
		RefreshTokenResponse: RefreshTokenResponse{
			Token: refreshToken.Token,
		},
	}, nil
}
