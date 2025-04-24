package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	JWT  JWTConfig
	DB   DBConfig
	Port string
}

type JWTConfig struct {
	secret []byte
}

func (j JWTConfig) GetSecret() []byte {
	return j.secret
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

	jwtSecret := mustGet("JWT_SECRET")
	return &Config{
		JWT: JWTConfig{
			secret: []byte(jwtSecret),
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

// func getOrDefault(key, defaultVal string) string {
// 	val := os.Getenv(key)
// 	if val == "" {
// 		return defaultVal
// 	}
// 	return val
// }
