package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type FuelPumpRepositoryPostgres struct {
	db *gorm.DB
}

func NewFuelPumpRepositoryPostgres(db *gorm.DB) *FuelPumpRepositoryPostgres {
	return &FuelPumpRepositoryPostgres{
		db: db,
	}
}

func (r *FuelPumpRepositoryPostgres) TakeOrRelease(pumpID uint) error {
	return nil
}

func (r *FuelPumpRepositoryPostgres) GetAll(stationID uint) ([]model.FuelPump, error) {
	var result []model.FuelPump

	err := r.db.Where("gas_station_id = ?", stationID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
