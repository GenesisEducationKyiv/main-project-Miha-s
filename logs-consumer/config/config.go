package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Credentials struct {
	Username string
	Password string
	Host     string
}

func LoadCredentials() Credentials {
	_ = godotenv.Load()
	return Credentials{
		Username: os.Getenv("RABBIT_USERNAME"),
		Password: os.Getenv("RABBIT_PASSWORD"),
		Host:     os.Getenv("RABBIT_URL_PORT"),
	}
}
