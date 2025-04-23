package user

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	user "github.com/sopial42/cleanic/internal/domains/user"
)

var (
	jwtSecret = []byte("your_secret_key")
	UserIDKey = ClaimsKey("user_id")
	RolesKey  = ClaimsKey("roles")
)

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

func (s SignedJWT) ParseClaims() (user.ID, user.Roles, error) {
	token, err := jwt.Parse(string(s.token), func(t *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if token == nil {
		return 0, user.Roles{}, errors.New("auth token not valid")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idFloat, ok := claims[string(UserIDKey)].(float64)
		if !ok {
			return 0, user.Roles{}, fmt.Errorf("user id is not a valid format")
		}

		id := user.ID(idFloat)
		rolesStr, ok := claims[string(RolesKey)].(string)
		if !ok {
			return 0, user.Roles{}, fmt.Errorf("roles are not a valid format")
		}

		currentRoles, err := user.NewRolesFromRolesString(rolesStr)
		return id, currentRoles, err
	}

	return 0, user.Roles{}, fmt.Errorf("unable to parse token: %w", err)
}

func generateJWT(userID user.ID, roles user.Roles) (SignedJWT, error) {
	claims := jwt.MapClaims{
		string(UserIDKey): userID,
		string(RolesKey):  roles.String(),
		"exp":             time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwt, err := token.SignedString(jwtSecret)
	if err != nil {
		return SignedJWT{}, fmt.Errorf("unable to sign token: %w", err)
	}

	return SignedJWT{token: jwt}, nil
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	raw := ctx.Value(UserIDKey)
	if raw == nil {
		return 0, fmt.Errorf("user id not found in context")
	}

	str, ok := raw.(string)
	if !ok {
		return 0, fmt.Errorf("user id is not a valid format (not string)")
	}

	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("user id is not a valid format (not int)")
	}

	return id, nil
}
