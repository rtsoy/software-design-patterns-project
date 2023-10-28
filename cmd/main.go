package main

import (
	"fmt"
	"log"

	"github.com/rtsoy/software-design-patterns-project/config"
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
	"github.com/rtsoy/software-design-patterns-project/internal/telebot"
	"github.com/rtsoy/software-design-patterns-project/pkg/postgres"
)

var models = []interface{}{
	model.GasStation{},
	model.FuelPump{},
	model.CreditCard{},
	model.Observer{},
}

func main() {
	// Load configuration settings from environment variables.
	cfg := config.Get()

	// Build the PostgreSQL DSN for database connection.
	dbDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser, cfg.PostgresUser, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName)

	// Initialize a connection to the PostgreSQL database and create tables for specified models.
	db, err := postgres.New(dbDSN, models)
	if err != nil {
		log.Fatalf("db initializition failure: %v", err)
	}

	// Create a repository using the database connection.
	repo := repository.NewRepository(db)

	// Initialize a Telegram bot with the repository.
	bot, err := telebot.New(repo)
	if err != nil {
		log.Fatalf("failed to initiazize a telegram bot: %v", err)
	}

	// TODO: Implement graceful shutdown logic.
	_ = bot
}
