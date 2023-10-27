package repository

import (
	"errors"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

var ErrForbidden = errors.New("forbidden")

type FuelPumpRepositoryPostgres struct {
	db *gorm.DB
}

func NewFuelPumpRepositoryPostgres(db *gorm.DB) *FuelPumpRepositoryPostgres {
	return &FuelPumpRepositoryPostgres{
		db: db,
	}
}

func (r *FuelPumpRepositoryPostgres) ReleaseAll(telegramID int64) ([]uint64, error) {
	var result []uint64

	err := r.db.Model(&model.FuelPump{}).Select("id").
		Where("telegram_id = ?", telegramID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	err = r.db.Model(&model.FuelPump{}).Where("telegram_id = ?", telegramID).Update(
		"is_available", true).Update(
		"telegram_id", 0).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *FuelPumpRepositoryPostgres) TakeOrRelease(telegramID int64, pumpID uint64) error {
	var pump model.FuelPump

	err := r.db.First(&pump, pumpID).Error
	if err != nil {
		return err
	}

	if !pump.IsAvailable && pump.TelegramID != telegramID {
		return ErrForbidden
	}

	if !pump.IsAvailable {
		err = r.db.Model(&model.FuelPump{}).Where("id = ?", pumpID).Update(
			"is_available", true).Update(
			"telegram_id", 0).Error

		return err
	}

	err = r.db.Model(&model.FuelPump{}).Where("id = ?", pumpID).Update(
		"is_available", false).Update(
		"telegram_id", telegramID).Error

	return err
}

func (r *FuelPumpRepositoryPostgres) GetAll(stationID uint64) ([]model.FuelPump, error) {
	var result []model.FuelPump

	err := r.db.Where("gas_station_id = ?", stationID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}
