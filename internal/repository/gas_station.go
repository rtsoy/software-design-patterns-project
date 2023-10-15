package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type GasStationRepositoryPostgres struct {
	db *gorm.DB
}

func NewGasStationRepository(db *gorm.DB) *GasStationRepositoryPostgres {
	return &GasStationRepositoryPostgres{
		db: db,
	}
}

func (r *GasStationRepositoryPostgres) GetAll() ([]model.GasStation, error) {
	var result []model.GasStation

	err := r.db.Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
