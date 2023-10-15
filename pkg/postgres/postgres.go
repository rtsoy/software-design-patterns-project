package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// New initializes a new PostgreSQL database connection and performs auto-migration for the provided models.
func New(dsn string, models []interface{}) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		return nil, err
	}

	// Auto-migrate each provided model to create corresponding database tables.
	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
