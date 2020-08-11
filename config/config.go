package config

import (
	"encoding/json"
	"github.com/spf13/viper"
)

type Config struct {
	BindAddress       string
	LogLevel          string
	LogPersistence    bool
	// DB connection
	DbDriver          string
	DbHost            string
	DbPort            string
	DbName            string
	DbUser            string
	DbPassword        string
	// Oauth
	OAuthEnable       bool
	OAuthIssuer       string
	OAuthClientId     string
	OAuthClientSecret string
	OAuthRedirectUrl  string
}

// Default config for running the service locally
var Defaults = map[string]string{
	"BIND_ADDRESS":    ":9090",
	"LOG_LEVEL":       "DEBUG",
	"LOG_PERSISTENCE": "true",
	"DB_DRIVER":       "sqlite3",
	"DB_NAME":         "test.db",
	"ENABLE_OAUTH":    "false",
}

func LoadConfig() (*Config, error) {
	configReader := viper.New()

	// 1) Set defaults
	configReader.SetDefault("BIND_ADDRESS", Defaults["BIND_ADDRESS"])
	configReader.SetDefault("LOG_LEVEL", Defaults["LOG_LEVEL"])
	configReader.SetDefault("LOG_PERSISTENCE", Defaults["LOG_LEVEL"])
	configReader.SetDefault("DB_DRIVER", Defaults["DB_DRIVER"])
	configReader.SetDefault("DB_NAME", Defaults["DB_NAME"])
	configReader.SetDefault("ENABLE_OAUTH", Defaults["ENABLE_OAUTH"])

	// 2) Load the environment variables
	configReader.AutomaticEnv()

	// 3) Look for .env file in the working directory
	configReader.SetConfigType("env")
	configReader.SetConfigName(".env")
	configReader.AddConfigPath(".")

	if err := configReader.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore, as it is optional
		} else {
			return nil, err
		}
	}

	config := Config{}

	config.BindAddress = configReader.GetString("BIND_ADDRESS")
	config.LogLevel = configReader.GetString("LOG_LEVEL")
	config.LogPersistence = configReader.GetBool("LOG_PERSISTENCE")

	config.DbDriver = configReader.GetString("DB_DRIVER")
	config.DbHost = configReader.GetString("DB_HOST")
	config.DbPort = configReader.GetString("DB_PORT")
	config.DbName = configReader.GetString("DB_NAME")
	config.DbUser = configReader.GetString("DB_USER")
	config.DbPassword = configReader.GetString("DB_PASSWORD")

	config.OAuthEnable = configReader.GetBool("ENABLE_OAUTH")
	config.OAuthIssuer = configReader.GetString("OAUTH_ISSUER")
	config.OAuthClientId = configReader.GetString("OAUTH_CLIENT_ID")
	config.OAuthClientSecret = configReader.GetString("OAUTH_CLIENT_SECRET")
	config.OAuthRedirectUrl = configReader.GetString("OAUTH_REDIRECT_URL")

	return &config, nil
}

func (c *Config) String() string {
	s, _ := json.MarshalIndent(c, "", "\t")
	return string(s)
}
