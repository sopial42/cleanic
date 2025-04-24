package jwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/sopial42/cleanic/internal/domains/user"
)

var (
	UserIDKey   = ClaimsKey("user_id")
	RolesKey    = ClaimsKey("roles")
	ExpireAtKey = ClaimsKey("exp")
)

type CustomClaims struct {
	UserID    user.ID
	UserRoles user.Roles
	ExpDate   time.Time
}

func newClaims(id user.ID, roles user.Roles) jwt.MapClaims {
	claims := jwt.MapClaims{
		string(UserIDKey):   id,
		string(RolesKey):    roles.String(),
		string(ExpireAtKey): time.Now().Add(24 * time.Hour).Unix(),
	}

	return claims
}

type ClaimsKey string

type SignedJWT struct {
	token string
}

func NewJWTFromAuthHeaderString(fullHeaderValue string) (SignedJWT, error) {
	var signedJWT SignedJWT

	if len(fullHeaderValue) == 0 {
		return SignedJWT{}, fmt.Errorf("empty token")
	}

	signedJWT.token = strings.TrimPrefix(fullHeaderValue, "Bearer ")
	if signedJWT.token == fullHeaderValue {
		return SignedJWT{}, fmt.Errorf("unable to parse Bearer token")
	}

	return signedJWT, nil
}

func (s SignedJWT) ToAuthHeaderString() string {
	return fmt.Sprintf("Bearer %s", string(s.token))
}

func GenerateJWT(userID user.ID, roles user.Roles, secret []byte) (SignedJWT, error) {
	claims := newClaims(userID, roles)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString(secret)
	if err != nil {
		return SignedJWT{}, fmt.Errorf("unable to sign token: %w", err)
	}

	return SignedJWT{token: jwt}, nil
}
func (s SignedJWT) ParseClaims(secret []byte) (CustomClaims, error) {
	var claims CustomClaims

	token, err := jwt.Parse(string(s.token), func(t *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || token == nil {
		return claims, fmt.Errorf("auth token not valid: %w", err)
	}

	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return claims, errors.New("invalid token claims")
	}

	// Parse and validate expiration date
	rawExp, exists := mapClaims[string(ExpireAtKey)]
	if !exists {
		return claims, errors.New("expiration date (exp) is missing")
	}

	expFloat, ok := rawExp.(float64)
	if !ok {
		return claims, errors.New("expiration date is not a float64")
	}

	expTime := time.Unix(int64(expFloat), 0)
	if time.Now().After(expTime) {
		return claims, errors.New("token is expired")
	}

	claims.ExpDate = expTime
	// Parse User ID
	idFloat, ok := mapClaims[string(UserIDKey)].(float64)
	if !ok {
		return claims, errors.New("user ID is not a valid float64")
	}

	claims.UserID = user.ID(idFloat)
	// Parse Roles
	rolesStr, ok := mapClaims[string(RolesKey)].(string)
	if !ok {
		return claims, errors.New("roles are not a valid string")
	}

	currentRoles, err := user.NewRolesFromRolesString(rolesStr)
	if err != nil {
		return claims, fmt.Errorf("unable to parse roles: %w", err)
	}

	claims.UserRoles = currentRoles
	return claims, nil
}
