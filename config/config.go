package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDBName   string `env:"POSTGRES_DB_NAME"`
}

func New() (*Config, error) {
	cfg := &Config{}

	return cfg, cleanenv.ReadEnv(cfg)
}
