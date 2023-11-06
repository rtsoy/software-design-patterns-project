// Factory + Decorator

package main

import (
	"log"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
)

// LoggerGasStationFactory is a decorator for GasStationFactory, which logs the creation of GasStation objects.
type LoggerGasStationFactory struct {
	factory GasStationFactory
}

// CreateGasStation creates a GasStation object while logging the process and forwarding the request to the underlying factory.
func (g *LoggerGasStationFactory) CreateGasStation(company, address string) (*model.GasStation, error) {
	log.Printf("Creating gas station: Company=%s, Address=%s", company, address)

	station, err := g.factory.CreateGasStation(company, address)

	if err != nil {
		log.Printf("Error creating gas station: %v", err)

		return nil, err
	}

	log.Printf("Gas station created successfully: ID=%d", station.ID)

	return station, nil
}

// LoggerFuelPumpFactory is a decorator for FuelPumpFactory, which logs the creation of FuelPump objects.
type LoggerFuelPumpFactory struct {
	factory FuelPumpFactory
}

// CreateFuelPump creates a FuelPump object while logging the process and forwarding the request to the underlying factory.
func (f *LoggerFuelPumpFactory) CreateFuelPump(fuelType model.FuelType, price, gasStationID uint64) (*model.FuelPump, error) {
	log.Printf("Creating fuel pump: FuelType=%s, Price=%d, GasStationID=%d", fuelType, price, gasStationID)

	pump, err := f.factory.CreateFuelPump(fuelType, price, gasStationID)

	if err != nil {
		log.Printf("Error creating fuel pump: %v", err)

		return nil, err
	}

	log.Printf("Fuel pump created successfully: ID=%d", pump.ID)

	return pump, nil
}
