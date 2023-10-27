package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type ObserverRepositoryPostgres struct {
	db *gorm.DB
}

func NewObserverRepositoryPostgres(db *gorm.DB) *ObserverRepositoryPostgres {
	return &ObserverRepositoryPostgres{
		db: db,
	}
}

func (r *ObserverRepositoryPostgres) GetAll(pumpID uint64) ([]int64, error) {
	var result []int64

	err := r.db.Model(model.Observer{}).Select("telegram_id").
		Where("fuel_pump_id = ?", pumpID).Find(&result).Error

	return result, err
}

func (r *ObserverRepositoryPostgres) Add(telegramID int64, pumpID uint64) error {
	err := r.db.Create(&model.Observer{
		TelegramID: telegramID,
		FuelPumpID: pumpID,
	}).Error

	return err
}

func (r *ObserverRepositoryPostgres) Delete(telegramID int64, pumpID uint64) error {
	err := r.db.Where("telegram_id = ? AND fuel_pump_id = ?",
		telegramID, pumpID).Delete(&model.Observer{}).Error

	return err
}
