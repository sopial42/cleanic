package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sopial42/cleanic/internal/domains/user"
)

type RefreshToken struct {
	SignedToken  SignedRefreshToken
	Claims RefreshTokenClaims
}

type SignedRefreshToken string

type RefreshTokenSecret []byte

type RefreshTokenTTLDays int

type RefreshTokenAudience string

type RefreshTokenClaims struct {
	ID        uuid.UUID
	Subject   user.ID
	ExpiresAt int64
	IssuedAt  int64
}

func NewRefreshToken(userID user.ID, secret RefreshTokenSecret, tokenTTLDays RefreshTokenTTLDays) (RefreshToken, error) {
	claims := generateRefreshTokenClaims(userID, tokenTTLDays)
	token, err := generateSignedRefreshToken(claims, secret)
	if err != nil {
		return RefreshToken{}, fmt.Errorf("unable to sign refresh with claims: %w", err)
	}

	return RefreshToken{
		SignedToken:  token,
		Claims: claims,
	}, nil
}

func generateRefreshTokenClaims(userID user.ID, tokenTTLDays RefreshTokenTTLDays) RefreshTokenClaims {
	return RefreshTokenClaims{
		ID:        uuid.New(),
		Subject:   userID,
		ExpiresAt: time.Now().Add(time.Duration(tokenTTLDays) * 24 * time.Hour).Unix(),
		IssuedAt:  time.Now().Unix(),
	}
}

func generateSignedRefreshToken(claims RefreshTokenClaims, secret RefreshTokenSecret) (SignedRefreshToken, error) {
	jwtClaims := jwt.MapClaims{
		string(IDKey):       claims.ID,
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

func ParseRefreshClaims(signedToken SignedRefreshToken, secret RefreshTokenSecret) (RefreshTokenClaims, error) {
	token, err := jwt.Parse(string(signedToken), func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || token == nil {
		return RefreshTokenClaims{}, fmt.Errorf("refresh token not valid: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return RefreshTokenClaims{}, errors.New("invalid refresh token claims")
	}

	// Parse the token ID
	rawID, exists := mapClaims[string(IDKey)]
	if !exists {
		return RefreshTokenClaims{}, errors.New("token ID is missing")
	}

	id, ok := rawID.(string)
	if !ok {
		return RefreshTokenClaims{}, errors.New("token ID is not a valid string")
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return RefreshTokenClaims{}, fmt.Errorf("token ID is not a valid UUID: %w", err)
	}

	// Subject
	rawSubject, exists := mapClaims[string(SubjectKey)]
	if !exists {
		return RefreshTokenClaims{}, errors.New("subject (user ID) is missing")
	}

	subjectIDFloat, ok := rawSubject.(float64)
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
		ID:        parsedID,
		Subject:   user.ID(subjectIDFloat),
		ExpiresAt: expTime.Unix(),
		IssuedAt:  iatTime.Unix(),
	}, nil
}
