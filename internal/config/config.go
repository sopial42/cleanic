package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	JWT  JWTConfig
	DB   DBConfig
	Port string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func Load() *Config {
	_ = godotenv.Load()

	// jwt token
	jwtAccessTokenSecret := mustGet("JWT_ACCESS_TOKEN_SECRET")
	jwtAccessTokenExpirationMinutesEnv := mustGet("JWT_ACCESS_TOKEN_TTL_MIN")
	jwtRefreshTokenSecret := mustGet("JWT_REFRESH_TOKEN_SECRET")
	jwtRefreshTokenExpirationMinutesEnv := mustGet("JWT_REFRESH_TOKEN_TTL_DAYS")

	atTTLMinutes, err := strconv.Atoi(jwtAccessTokenExpirationMinutesEnv)
	if err != nil {
		log.Fatalf("unable to cast access token ttl :%v", err)
	}

	etTTLDays, err := strconv.Atoi(jwtRefreshTokenExpirationMinutesEnv)
	if err != nil {
		log.Fatalf("unable to cast refresh token ttl :%v", err)
	}

	// JWT Cookies
	cookieSecret := mustGet("JWT_COOKIE_SECRET")
	cookieDomain := mustGet("JWT_COOKIE_DOMAIN")
	cookiePath := mustGet("JWT_COOKIE_PATH")
	sameSite := mustGet("JWT_COOKIE_SAME_SITE")
	sameSiteIota := parseSameSite(sameSite)
	if sameSiteIota == -1 {
		log.Fatalf("unable to cast cookie same site :%v", sameSite)
	}

	cookieMaxAgeSeconds, err := strconv.Atoi(mustGet("JWT_COOKIE_MAX_AGE_SECONDS"))
	if err != nil {
		log.Fatalf("unable to cast cookie max age seconds :%v", err)
	}

	secure := mustGet("JWT_COOKIE_SECURE")
	secureBool, err := strconv.ParseBool(secure)
	if err != nil {
		log.Fatalf("unable to cast cookie secure :%v", err)
	}

	return &Config{
		JWT: JWTConfig{
			AccessTokenConfig: AccessTokenConfig{
				secret:                 []byte(jwtAccessTokenSecret),
				TokenExpirationMinutes: atTTLMinutes,
			},
			RefreshTokenConfig: RefreshTokenConfig{
				secret:                 []byte(jwtRefreshTokenSecret),
				TokenExpirationMinutes: etTTLDays,
			},
			CookieStoreConfig: CookieStoreConfig{
				Domain:        cookieDomain,
				MaxAgeSeconds: cookieMaxAgeSeconds,
				Path:          cookiePath,
				SameSite:      sameSiteIota,
				Secret:        []byte(cookieSecret),
				Secure:        secureBool,
			},
		},
		DB: DBConfig{
			Host:     mustGet("DB_HOST"),
			Port:     mustGet("DB_PORT"),
			User:     mustGet("DB_USER"),
			Password: mustGet("DB_PASSWORD"),
			DBName:   mustGet("DB_NAME"),
		},
		Port: mustGet("PORT"),
	}
}

func mustGet(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return val
}
