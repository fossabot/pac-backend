package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	BindAddress string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	return &Config{
		BindAddress: os.Getenv("BIND_ADDRESS"),
	}, nil
}
