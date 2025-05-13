package config

import (
	"net/http"
	"strings"
)

type JWTConfig struct {
	AccessTokenConfig  AccessTokenConfig
	RefreshTokenConfig RefreshTokenConfig
	CookieStoreConfig  CookieStoreConfig
}

type CookieStoreConfig struct {
	Domain        string
	MaxAgeSeconds int
	Path          string
	// SameSiteDefaultMode = 0/1
	// SameSiteLaxMode = 2
	// SameSiteStrictMode = 3
	// SameSiteNoneMode = 4
	SameSite http.SameSite
	Secret   []byte
	Secure   bool
}

type AccessTokenConfig struct {
	secret                 []byte
	TokenExpirationMinutes int
	Audience               string
}

func (a *AccessTokenConfig) GetSecret() []byte {
	return a.secret
}

type RefreshTokenConfig struct {
	secret                 []byte
	TokenExpirationMinutes int
	Audience               string
}

func (r *RefreshTokenConfig) GetSecret() []byte {
	return r.secret
}

func parseSameSite(value string) http.SameSite {
	switch strings.ToLower(value) {
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	case "default":
		return http.SameSiteDefaultMode
	default:
		return -1
	}
}
