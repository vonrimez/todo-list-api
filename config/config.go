package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URL     string
	PORT       string
	JWT_SECRET string
}

func LoadConfig() *Config {
	godotenv.Load()
	return &Config{
		DB_URL:     os.Getenv("DB_URL"),
		PORT:       os.Getenv("PORT"),
		JWT_SECRET: os.Getenv("JWT_SECRET"),
	}
}
