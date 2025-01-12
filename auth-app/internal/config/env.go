package config

import (
	"github.com/caarlos0/env/v11"
	"github.com/rs/zerolog/log"
)

func Load() *Config {
	var cfg Config
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	return &cfg
}

type Config struct {
	App      App
	Database Database
}

type App struct {
	Port      int    `env:"APP_PORT" envDefault:"8080"`
	JwtSecret string `env:"JWT_SECRET_KEY"`
}

type Database struct {
	Host     string `env:"DB_HOST"`
	Port     int    `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	DbName   string `env:"DB_NAME"`
}
