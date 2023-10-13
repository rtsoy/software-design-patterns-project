package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
}

func New() (*Config, error) {
	cfg := &Config{}

	return cfg, cleanenv.ReadEnv(cfg)
}
