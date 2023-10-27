package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type ObserverRepository interface {
	GetAll(pumpID uint64) ([]int64, error)
	Add(telegramID int64, pumpID uint64) error
	Delete(telegramID int64, pumpID uint64) error
}

type FuelPumpRepository interface {
	GetAll(stationID uint64) ([]model.FuelPump, error)
	TakeOrRelease(telegramID int64, pumpID uint64) error
	ReleaseAll(telegramID int64) ([]uint64, error)
}

type GasStationRepository interface {
	GetAll() ([]model.GasStation, error)
}

type CreditCardRepository interface {
	GetAll(telegramID int64) ([]model.CreditCard, error)
	Create(card model.CreditCard) error
}

type Repository struct {
	CreditCardRepository
	GasStationRepository
	FuelPumpRepository
	ObserverRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		CreditCardRepository: NewCreditCardRepository(db),
		GasStationRepository: NewGasStationRepository(db),
		FuelPumpRepository:   NewFuelPumpRepositoryPostgres(db),
		ObserverRepository:   NewObserverRepositoryPostgres(db),
	}
}
