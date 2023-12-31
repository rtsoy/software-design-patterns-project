package main

import (
	"fmt"
	"log"

	"github.com/rtsoy/software-design-patterns-project/config"
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var (
		// Get application configuration settings.
		cfg = config.Get()

		// Build a Data Source Name (DSN) string for the PostgreSQL database connection.
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
			cfg.PostgresUser, cfg.PostgresUser, cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresDBName)
	)

	// Initialize a GORM database connection.
	db, err := gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("connection to db failure: %v", err)
	}

	// Check if the database is empty by attempting to find GasStation records.
	tx := db.Find(&[]model.GasStation{})
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	if tx.RowsAffected > 0 {
		// If the database is not empty, print a message and exit the program.
		fmt.Println("db is not empty, exiting...")
		return
	}

	var (
		// Define station factories for creating GasStation objects.
		stationFactory       GasStationFactory = &ConcreteGasStationFactory{db: db}
		loggerStationFactory GasStationFactory = &LoggerGasStationFactory{factory: stationFactory}

		// Define pump factories for creating FuelPump objects.
		pumpFactory       FuelPumpFactory = &ConcreteFuelPumpFactory{db: db}
		loggerPumpFactory FuelPumpFactory = &LoggerFuelPumpFactory{factory: pumpFactory}
	)

	// Define data for creating GasStation and FuelPump objects.
	data := []struct {
		station model.GasStation
		pumps   []model.FuelPump
	}{
		{
			station: model.GasStation{Company: "Gas Energy", Address: "просп. Кабанбай Батыра 45B"},
			pumps: []model.FuelPump{
				{
					FuelType:    model.Gas,
					Price:       76,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI92,
					Price:       205,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI95,
					Price:       265,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI98,
					Price:       310,
					IsAvailable: true,
				},
				{
					FuelType:    model.Diesel,
					Price:       295,
					IsAvailable: true,
				},
			},
		},
		{
			station: model.GasStation{Company: "Qazaq Oil", Address: "шоссе Каркаралы 1/1"},
			pumps: []model.FuelPump{
				{
					FuelType:    model.AI92,
					Price:       205,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI95,
					Price:       260,
					IsAvailable: true,
				},
				{
					FuelType:    model.Diesel,
					Price:       295,
					IsAvailable: true,
				},
			},
		},
		{
			station: model.GasStation{Company: "Helios", Address: "ул. Енбекшилер 18/1"},
			pumps: []model.FuelPump{
				{
					FuelType:    model.AI92,
					Price:       205,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI95,
					Price:       249,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI98,
					Price:       248,
					IsAvailable: true,
				},
			},
		},
		{
			station: model.GasStation{Company: "Газпромнефть", Address: "просп. Богенбай батыра 24"},
			pumps: []model.FuelPump{
				{
					FuelType:    model.AI92,
					Price:       205,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI95,
					Price:       255,
					IsAvailable: true,
				},
				{
					FuelType:    model.AI98,
					Price:       280,
					IsAvailable: true,
				},
				{
					FuelType:    model.Diesel,
					Price:       295,
					IsAvailable: true,
				},
			},
		},
	}

	// Loop through the data and create GasStation and FuelPump objects.
	for _, elem := range data {
		station, _ := loggerStationFactory.CreateGasStation(elem.station.Address, elem.station.Company)

		for _, pump := range elem.pumps {
			loggerPumpFactory.CreateFuelPump(pump.FuelType, pump.Price, station.ID)
		}
	}
}
