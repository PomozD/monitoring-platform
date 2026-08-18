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

	return Config{
		AppPort:     appPort,
		DatabaseURL: databaseURL,
		Environment: environment,
	}, nil
}
