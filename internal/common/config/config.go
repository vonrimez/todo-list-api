package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL     string
	PORT      string
	JWTSecret string
}

func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		DBURL:     os.Getenv("DB_URL"),
		PORT:      os.Getenv("PORT"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}
