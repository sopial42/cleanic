package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

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

func (a *authSVC) Login(ctx context.Context, loginUser user.User) (utils.RefreshToken, utils.AccessToken, error) {
	userFound, err := a.uClient.GetUserByEmail(ctx, loginUser.Email)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, err
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(loginUser.Password)); err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("invalid password: %w", err)
	}

	refreshToken, accessToken, err := generateTokens(userFound, a.jwtConfig)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to generate tokens: %w", err)
	}

	if err := a.persistence.StoreRefreshTokenClaims(ctx, refreshToken.Claims); err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to store refresh token: %w", err)
	}

	return refreshToken, accessToken, nil
}

func (a *authSVC) Logout(ctx context.Context, userID user.ID) error {
	if err := a.persistence.DeleteRefreshTokenClaims(ctx, userID); err != nil {
		return fmt.Errorf("unable to delete refresh token: %w", err)
	}

	return nil
}

// Refresh ensure the received token is the one associated to the user in DB
// If yes, rotate tokens and return the new ones
// Only one token per user is valid at a time
func (a *authSVC) Refresh(ctx context.Context, signedTokenCandidate utils.SignedRefreshToken) (utils.RefreshToken, utils.AccessToken, error) {
	claimsCandidate, err := utils.ParseRefreshClaims(signedTokenCandidate, a.jwtConfig.RefreshTokenConfig.GetSecret())
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to parse refresh token: %w", err)
	}

	currentUserID := claimsCandidate.Subject

	// The current refresh token is expired
	// Login has to be done again
	if claimsCandidate.ExpiresAt < time.Now().Unix() {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("candidate refresh token expired")
	}

	// Get the refresh token from DB,
	// Only one is valid per user at a time
	// If the token is not found, it means the user has logged out
	// or the token has been revoked
	// then he need to login again
	peristedClaims, err := a.persistence.GetRefreshTokenClaimsByUserID(ctx, currentUserID)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to get refresh token from DB: %w", err)
	}

	// If the refresh token in DB is expired, the user need to login again
	if peristedClaims.ExpiresAt < time.Now().Unix() {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("stored refresh token expired")
	}

	// Check if the token ID in DB is the same as the one received
	// If not, that means the user may have sent a token from another device
	// then he need to login again on the new device
	// We could add much more security here by checking the IP address, user agent, etc. and ensure the reason for an old token to be used is legit
	if peristedClaims.ID != claimsCandidate.ID {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("refresh token candidate is not anymore associated to the user")
	}

	// Get user details to get roles
	userFound, err := a.uClient.GetUserByID(ctx, currentUserID)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, err
	}

	// Rotate tokens
	refreshToken, accessToken, err := generateTokens(userFound, a.jwtConfig)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to generate tokens: %w", err)
	}

	if err := a.persistence.StoreRefreshTokenClaims(ctx, refreshToken.Claims); err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to store refresh token: %w", err)
	}

	return refreshToken, accessToken, nil
}

func generateTokens(user user.User, config config.JWTConfig) (utils.RefreshToken, utils.AccessToken, error) {
	refreshToken, err := utils.NewRefreshToken(
		user.ID, config.RefreshTokenConfig.GetSecret(),
		utils.RefreshTokenTTLDays(config.RefreshTokenConfig.TokenExpirationMinutes),
	)
	if err != nil {
		return utils.RefreshToken{}, utils.AccessToken{}, fmt.Errorf("unable to generate refresh token: %w", err)
	}

	accessToken, err := utils.NewAccessToken(
		user.ID, user.Roles, config.AccessTokenConfig.GetSecret(),
		utils.AccessTokenTTLMin(config.AccessTokenConfig.TokenExpirationMinutes),
	)

	return refreshToken, accessToken, err
}
