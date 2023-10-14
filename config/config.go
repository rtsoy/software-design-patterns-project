package config

import (
	"log"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

// Config represents configuration settings for application.
type Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`

	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresDBName   string `env:"POSTGRES_DB_NAME"`
}

var (
	cfg  *Config
	once sync.Once
)

// Get returns a singleton instance of the Config struct.
func Get() *Config {
	// Using the sync.Once package to ensure that the instance is created only once.
	once.Do(func() {
		log.Printf("creating a new config instance now")
		cfg = &Config{}

		if err := cleanenv.ReadEnv(cfg); err != nil {
			log.Fatalf("error reading environment variables: %v", err)
		}
	})

	return cfg
}
