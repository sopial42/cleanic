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
	Token              SignedAccessToken
	ExpirationDuration time.Duration
	Type               string
}

type SignedAccessToken string

type AccessTokenSecret []byte

type AccessTokenTTLMin int

type AccessTokenAudience string

type accessTokenClaims struct {
	Subject   user.ID
	ExpiresAt time.Time
	IssuedAt  time.Time
	Audience  string
	Roles     user.Roles
}

func NewAccessToken(userID user.ID, roles user.Roles, secret AccessTokenSecret, tokenTTLMin AccessTokenTTLMin, audience AccessTokenAudience) (AccessToken, error) {
	claims := generateAccessTokenClaims(userID, roles, audience, tokenTTLMin)
	token, err := generateSignedAccessToken(claims, secret)
	if err != nil {
		return AccessToken{}, fmt.Errorf("unable to generate access token: %w", err)
	}

	return AccessToken{
		Token:              token,
		ExpirationDuration: time.Duration(tokenTTLMin) * time.Minute,
		Type:               AccessTokenType,
	}, nil
}

func generateAccessTokenClaims(userID user.ID, roles user.Roles, audience AccessTokenAudience, expirationMinutes AccessTokenTTLMin) accessTokenClaims {
	return accessTokenClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Duration(expirationMinutes) * time.Minute),
		IssuedAt:  time.Now(),
		Audience:  string(audience),
		Roles:     roles,
	}
}

func generateSignedAccessToken(claims accessTokenClaims, secret []byte) (SignedAccessToken, error) {
	jwtClaims := jwt.MapClaims{
		string(SubjectKey):  claims.Subject,
		string(ExpireAtKey): claims.ExpiresAt,
		string(IssuedAtKey): claims.IssuedAt,
		string(AudienceKey): claims.Audience,
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
	return fmt.Sprintf("Bearer %s", string(token.Token))
}

func ParseAccessClaims(tokenToParse string, secret []byte) (accessTokenClaims, error) {
	token, err := jwt.Parse(tokenToParse, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || token == nil {
		return accessTokenClaims{}, fmt.Errorf("auth token not valid: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return accessTokenClaims{}, errors.New("invalid token claims")
	}

	// Subject
	idFloat, ok := mapClaims[string(AudienceKey)].(float64)
	if !ok {
		return accessTokenClaims{}, errors.New("user ID is not a valid float64")
	}

	// ExpiresAt
	rawExp, exists := mapClaims[string(ExpireAtKey)]
	if !exists {
		return accessTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	expFloat, ok := rawExp.(float64)
	if !ok {
		return accessTokenClaims{}, errors.New("expiration date is not a float64")
	}

	expTime := time.Unix(int64(expFloat), 0)

	// IssuedAt
	rawIAT, exists := mapClaims[string(IssuedAtKey)]
	if !exists {
		return accessTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	iatFloat, ok := rawIAT.(float64)
	if !ok {
		return accessTokenClaims{}, errors.New("expiration date is not a float64")
	}

	iatTime := time.Unix(int64(iatFloat), 0)

	// Audience
	audience, ok := mapClaims[string(AudienceKey)].(string)
	if !ok {
		return accessTokenClaims{}, errors.New("audience is not a valid string")
	}

	// audience := regexp.MustCompile(",").Split(audienceStr, -1)

	// Parse Roles
	rolesStr, ok := mapClaims[string(RolesKey)].(string)
	if !ok {
		return accessTokenClaims{}, errors.New("roles are not a valid string")
	}

	currentRoles, err := user.NewRolesFromRolesString(rolesStr)
	if err != nil {
		return accessTokenClaims{}, fmt.Errorf("unable to parse roles: %w", err)
	}

	return accessTokenClaims{
		Subject:   user.ID(idFloat),
		ExpiresAt: expTime,
		IssuedAt:  iatTime,
		Audience:  audience,
		Roles:     currentRoles,
	}, nil
}
