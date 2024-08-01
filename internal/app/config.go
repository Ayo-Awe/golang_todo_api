package app

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PORT   int    `envconfig:"PORT" default:"8080"`
	DB_URL string `envconfig:"DATABASE_URL" required:"true"`
}

func LoadConfig() (*Config, error) {
	cfg := Config{}

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
