package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	Environment string
}

func Load() (Config, error) {
	appPort := os.Getenv("AUTH_SERVICE_PORT")
	if appPort == "" {
		appPort = "8081"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	environment := os.Getenv("APP_ENV")
	if environment == "" {
		environment = "development"
	}

	config := Config{
		AppPort:     appPort,
		DatabaseURL: databaseURL,
		Environment: environment,
	}

	if err := config.Validate(); err != nil {
		return Config{}, err
	}

	return config, nil
}

func (c Config) Validate() error {
	if c.AppPort == "" {
		return fmt.Errorf("AUTH_SERVICE_PORT cannot be empty")
	}

	if c.DatabaseURL == "" {
		return fmt.Errorf("DATABASE_URL cannot be empty")
	}

	switch c.Environment {
	case "development", "testing", "staging", "production":
	default:
		return fmt.Errorf("invalid APP_ENV: %s", c.Environment)
	}

	return nil
}
