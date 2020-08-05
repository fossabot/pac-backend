package config

import (
	"encoding/json"
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	BindAddress    string
	LogLevel       string
	LogPersistence bool
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load(".env")

	if err != nil {
		return nil, err
	}

	bindAddress := os.Getenv("BIND_ADDRESS")
	logLevel := os.Getenv("LOG_LEVEL")
	logPersistence, _ := strconv.ParseBool(os.Getenv("LOG_PERSISTENCE"))

	return &Config{
		BindAddress:    bindAddress,
		LogLevel:       logLevel,
		LogPersistence: logPersistence,
	}, nil
}

func (c *Config) String() string {
	s, _ := json.MarshalIndent(c, "", "\t")
	return string(s)
}
