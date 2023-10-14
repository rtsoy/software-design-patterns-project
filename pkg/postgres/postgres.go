package postgres

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string, models []interface{}) (*gorm.DB, error) {
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
