package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string `env:"DATABASE_USER" env-required:"true"`
	DBPassword string `env:"DATABASE_PASSWORD" env-required:"true"`
	DBHost     string `env:"DATABASE_HOST" env-required:"true"`
	DBPort     string `env:"DATABASE_PORT" env-default:"5432"`
	DBName     string `env:"DATABASE_NAME" env-required:"true"`

	Port string `env:"SERVER_PORT" env-default:"8080"`

	JWTKey []byte `env:"JWT_KEY" env-required:"true"`
}

func Load() (*Config, error) {
	var cfg Config

	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
