// Factory

package main

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

// GasStationFactory is an interface for creating GasStation objects.
type GasStationFactory interface {
	CreateGasStation(company, address string) (*model.GasStation, error)
}

// ConcreteGasStationFactory is an implementation of GasStationFactory that uses GORM to create GasStation records.
type ConcreteGasStationFactory struct {
	db *gorm.DB
}

// CreateGasStation creates a GasStation record with the given company and address.
func (g *ConcreteGasStationFactory) CreateGasStation(company, address string) (*model.GasStation, error) {
	station := &model.GasStation{
		Company: company,
		Address: address,
	}

	if err := g.db.Create(station).Error; err != nil {
		return nil, err
	}

	return station, nil
}

// FuelPumpFactory is an interface for creating FuelPump objects.
type FuelPumpFactory interface {
	CreateFuelPump(fuelType model.FuelType, price, gasStationID uint64) (*model.FuelPump, error)
}

// ConcreteFuelPumpFactory is an implementation of FuelPumpFactory that uses GORM to create FuelPump records.
type ConcreteFuelPumpFactory struct {
	db *gorm.DB
}

// CreateFuelPump creates a FuelPump record with the given fuel type, price, and gasStationID.
func (f *ConcreteFuelPumpFactory) CreateFuelPump(fuelType model.FuelType, price, gasStationID uint64) (*model.FuelPump, error) {
	pump := &model.FuelPump{
		FuelType:     fuelType,
		Price:        price,
		GasStationID: gasStationID,
	}

	if err := f.db.Create(pump).Error; err != nil {
		return nil, err
	}

	return pump, nil
}
