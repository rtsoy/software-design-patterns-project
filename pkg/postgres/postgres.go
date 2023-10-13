package postgres

import (
	"fmt"

	"github.com/rtsoy/software-design-patterns-project/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config, models []interface{}) (*gorm.DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.PostgresUser, cfg.PostgresUser, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName)

	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
