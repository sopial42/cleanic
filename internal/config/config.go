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

func getOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
