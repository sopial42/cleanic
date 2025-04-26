package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sopial42/cleanic/internal/domains/user"
)

type RefreshToken struct {
	Token  SignedRefreshToken
	Claims RefreshTokenClaims
}

type SignedRefreshToken string

type RefreshTokenSecret []byte

type RefreshTokenTTLDays int

type RefreshTokenAudience string

type RefreshTokenClaims struct {
	Subject   user.ID
	ExpiresAt time.Time
	IssuedAt  time.Time
}

func NewRefreshToken(userID user.ID, secret RefreshTokenSecret, tokenTTLDays RefreshTokenTTLDays, audience RefreshTokenAudience) (RefreshToken, error) {
	claims := generateRefreshTokenClaims(userID, audience, tokenTTLDays)
	token, err := generateRefreshToken(claims, secret)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("unable to sign refresh with claims: %w", err)
	}

	return RefreshToken{
		Token:  token,
		Claims: claims,
	}, nil
}

func generateRefreshTokenClaims(userID user.ID, audience RefreshTokenAudience, tokenTTLDays RefreshTokenTTLDays) RefreshTokenClaims {
	return RefreshTokenClaims{
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Duration(tokenTTLDays) * time.Hour),
		IssuedAt:  time.Now(),
	}
}

func generateRefreshToken(claims RefreshTokenClaims, secret RefreshTokenSecret) (SignedRefreshToken, error) {
	jwtClaims := jwt.MapClaims{
		string(SubjectKey):  claims.Subject,
		string(ExpireAtKey): claims.ExpiresAt,
		string(IssuedAtKey): claims.IssuedAt,
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := tokenWithClaims.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return SignedRefreshToken(signedToken), nil

}

func ParseRefreshClaims(tokenToParse string, secret RefreshTokenSecret) (RefreshTokenClaims, error) {
	token, err := jwt.Parse(tokenToParse, func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	if err != nil || token == nil {
		return RefreshTokenClaims{}, fmt.Errorf("auth token not valid: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return RefreshTokenClaims{}, errors.New("invalid token claims")
	}

	// Subject
	idFloat, ok := mapClaims[string(AudienceKey)].(float64)
	if !ok {
		return RefreshTokenClaims{}, errors.New("user ID is not a valid float64")
	}

	// ExpiresAt
	rawExp, exists := mapClaims[string(ExpireAtKey)]
	if !exists {
		return RefreshTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	expFloat, ok := rawExp.(float64)
	if !ok {
		return RefreshTokenClaims{}, errors.New("expiration date is not a float64")
	}

	expTime := time.Unix(int64(expFloat), 0)

	// IssuedAt
	rawIAT, exists := mapClaims[string(IssuedAtKey)]
	if !exists {
		return RefreshTokenClaims{}, errors.New("expiration date (exp) is missing")
	}

	iatFloat, ok := rawIAT.(float64)
	if !ok {
		return RefreshTokenClaims{}, errors.New("expiration date is not a float64")
	}

	iatTime := time.Unix(int64(iatFloat), 0)
	return RefreshTokenClaims{
		Subject:   user.ID(idFloat),
		ExpiresAt: expTime,
		IssuedAt:  iatTime,
	}, nil
}
