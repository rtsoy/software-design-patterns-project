package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type ObserverRepository interface {
	GetAll(pumpID uint) ([]uint, error)
	Add(telegramID, pumpID uint) error
	Delete(telegramID, pumpID uint) error
}

type FuelPumpRepository interface {
	GetAll(stationID uint) ([]model.FuelPump, error)
	TakeOrRelease(telegramID, pumpID uint) error
	ReleaseAll(telegramID uint) ([]uint, error)
}

type GasStationRepository interface {
	GetAll() ([]model.GasStation, error)
}

type CreditCardRepository interface {
	GetAll(telegramID uint) ([]model.CreditCard, error)
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
