package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type FuelPumpRepository interface {
	GetAll(stationID uint) ([]model.FuelPump, error)
	TakeOrRelease(pumpID uint) error
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
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		CreditCardRepository: NewCreditCardRepository(db),
		GasStationRepository: NewGasStationRepository(db),
		FuelPumpRepository:   NewFuelPumpRepositoryPostgres(db),
	}
}
