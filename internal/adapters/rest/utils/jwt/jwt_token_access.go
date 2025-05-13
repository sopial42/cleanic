package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sopial42/cleanic/internal/domains/user"
)

const AccessTokenType = "Bearer"

type AccessToken struct {
	SignedToken           SignedAccessToken
	ExpirationDurationMin time.Duration
	Type                  string
	Claims                AccessTokenClaims
}

type SignedAccessToken string

type AccessTokenSecret []byte

type AccessTokenTTLMin int

type AccessTokenAudience string

type AccessTokenClaims struct {
	Subject   user.ID
	ExpiresAt int64
	IssuedAt  int64
	Roles     user.Roles
}

func NewAccessToken(userID user.ID, roles user.Roles, secret AccessTokenSecret, tokenTTLMin AccessTokenTTLMin, audience AccessTokenAudience) (AccessToken, error) {
	claims := generateAccessTokenClaims(userID, roles, audience, tokenTTLMin)
	token, err := generateSignedAccessToken(claims, secret)
	if err != nil {
		return AccessToken{}, fmt.Errorf("unable to generate access token: %w", err)
	}

	return AccessToken{
		SignedToken:                 token,
		ExpirationDurationMin: time.Duration(tokenTTLMin) * time.Minute,
		Type:                  AccessTokenType,
		Claims:                claims,
	}, nil
}

func generateAccessTokenClaims(userID user.ID, roles user.Roles, audience AccessTokenAudience, expirationMinutes AccessTokenTTLMin) AccessTokenClaims {
	return AccessTokenClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Duration(expirationMinutes) * time.Minute).Unix(),
		IssuedAt:  time.Now().Unix(),
		Roles:     roles,
	}
}

func generateSignedAccessToken(claims AccessTokenClaims, secret []byte) (SignedAccessToken, error) {
	jwtClaims := jwt.MapClaims{
		string(SubjectKey):  claims.Subject,
		string(ExpireAtKey): claims.ExpiresAt,
		string(IssuedAtKey): claims.IssuedAt,
		string(RolesKey):    claims.Roles,
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := tokenWithClaims.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return SignedAccessToken(signedToken), nil
}

func ParseBearerHeader(fullHeaderValue string) (string, error) {
	if len(fullHeaderValue) == 0 {
		return "", fmt.Errorf("empty Bearer header")
	}

	token := strings.TrimPrefix(fullHeaderValue, "Bearer ")
	if token == fullHeaderValue {
		return "", fmt.Errorf("unable to parse Bearer token")
	}

	return token, nil
}

func BearerHeaderStringFromAccessToken(token AccessToken) string {
	return fmt.Sprintf("Bearer %s", string(token.SignedToken))
}

func ParseAccessClaims(tokenToParse string, secret []byte) (AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenToParse, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || token == nil {
		return AccessTokenClaims{}, fmt.Errorf("auth token not valid: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return AccessTokenClaims{}, errors.New("invalid token claims")
	}

	// Subject
	idFloat, ok := mapClaims[string(SubjectKey)].(float64)
	if !ok {
		return AccessTokenClaims{}, errors.New("user ID is not a valid float64")
	}

	// ExpiresAt
	rawExp, exists := mapClaims[string(ExpireAtKey)]
	if !exists {
		return AccessTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	expFloat, ok := rawExp.(float64)
	if !ok {
		return AccessTokenClaims{}, errors.New("expiration date is not a float64")
	}

	expTime := time.Unix(int64(expFloat), 0)

	// IssuedAt
	rawIAT, exists := mapClaims[string(IssuedAtKey)]
	if !exists {
		return AccessTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	iatFloat, ok := rawIAT.(float64)
	if !ok {
		return AccessTokenClaims{}, errors.New("expiration date is not a float64")
	}

	iatTime := time.Unix(int64(iatFloat), 0)

	// Parse Roles
	rolesStr, ok := mapClaims[string(RolesKey)].(string)
	if !ok {
		return AccessTokenClaims{}, errors.New("roles are not a valid string")
	}

	currentRoles, err := user.NewRolesFromRolesString(rolesStr)
	if err != nil {
		return AccessTokenClaims{}, fmt.Errorf("unable to parse roles: %w", err)
	}

	return AccessTokenClaims{
		Subject:   user.ID(idFloat),
		ExpiresAt: expTime.Unix(),
		IssuedAt:  iatTime.Unix(),
		Roles:     currentRoles,
	}, nil
}
