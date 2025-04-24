package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestParseClaims_ExpirationDate(t *testing.T) {
	secret := []byte("test_secret")

	tests := []struct {
		name        string
		exp         time.Time
		shouldError bool
		errorText   string
	}{
		{
			name:        "valid token",
			exp:         time.Now().Add(10 * time.Minute),
			shouldError: false,
		},
		{
			name:        "expired token",
			exp:         time.Now().Add(-10 * time.Minute),
			shouldError: true,
			errorText:   "token is expired",
		},
		{
			name:        "missing expiration",
			exp:         time.Time{}, // special case handled separately
			shouldError: true,
			errorText:   "expiration date (exp) is missing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims := jwt.MapClaims{
				string(UserIDKey): float64(1),
				string(RolesKey):  "admin",
			}

			if !tt.exp.IsZero() {
				claims[string(ExpireAtKey)] = float64(tt.exp.Unix())
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			signed, err := token.SignedString(secret)
			require.NoError(t, err)

			s := SignedJWT{token: signed}
			parsedClaims, err := s.ParseClaims(secret)

			if tt.shouldError {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errorText)
			} else {
				require.NoError(t, err)
				require.WithinDuration(t, tt.exp, parsedClaims.ExpDate, time.Second)
			}
		})
	}
}
