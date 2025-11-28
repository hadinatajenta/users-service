package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application level configuration.
type Config struct {
	AppPort string
	DB      DBConfig
}

// DBConfig groups database settings so they can be passed around cleanly.
type DBConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() (Config, error) {
	// Overload to ensure values from .env always replace any already-set env vars.
	err := godotenv.Overload()
	if err != nil {
		return Config{}, errors.New("failed to load env")
	}

	portStr := os.Getenv("DB_PORT")
	if portStr == "" {
		portStr = "5432"
	}
	dbPort, err := strconv.Atoi(portStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	cfg := Config{
		AppPort: os.Getenv("APP_PORT"),
		DB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     dbPort,
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			Name:     os.Getenv("DB_NAME"),
			SSLMode:  os.Getenv("DB_SSLMODE"),
		},
	}

	if cfg.DB.Name == "" {
		return Config{}, errors.New("DB_NAME must be set")
	}

	return cfg, nil
}
