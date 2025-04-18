package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig DBConfig
	Port     string
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

	return &Config{
		DBConfig: DBConfig{
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
