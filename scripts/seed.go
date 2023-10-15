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

	// Gas Energy

	gasEnergy := model.GasStation{
		Company: "Gas Energy",
		Address: "просп. Кабанбай Батыра 45B",
	}

	err = db.Create(&gasEnergy).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NEW RECORD: (GasStation) %+v\n", gasEnergy)

	gasEnergyPumps := []model.FuelPump{
		{
			FuelType:     model.Gas,
			Price:        76,
			IsAvailable:  true,
			GasStationID: gasEnergy.ID,
		},
		{
			FuelType:     model.AI92,
			Price:        205,
			IsAvailable:  true,
			GasStationID: gasEnergy.ID,
		},
		{
			FuelType:     model.AI95,
			Price:        265,
			IsAvailable:  true,
			GasStationID: gasEnergy.ID,
		},
		{
			FuelType:     model.AI98,
			Price:        310,
			IsAvailable:  true,
			GasStationID: gasEnergy.ID,
		},
		{
			FuelType:     model.Diesel,
			Price:        295,
			IsAvailable:  true,
			GasStationID: gasEnergy.ID,
		},
	}

	for _, pump := range gasEnergyPumps {
		err = db.Create(&pump).Error
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("NEW RECORD: (FuelPump) %+v\n", pump)
	}

	// Qazaq Oil

	qazaqOil := model.GasStation{
		Company: "Qazaq Oil",
		Address: "шоссе Каркаралы 1/1",
	}

	err = db.Create(&qazaqOil).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NEW RECORD: (GasStation) %+v\n", qazaqOil)

	qazaqOilPumps := []model.FuelPump{
		{
			FuelType:     model.AI92,
			Price:        205,
			IsAvailable:  true,
			GasStationID: qazaqOil.ID,
		},
		{
			FuelType:     model.AI95,
			Price:        260,
			IsAvailable:  true,
			GasStationID: qazaqOil.ID,
		},
		{
			FuelType:     model.Diesel,
			Price:        295,
			IsAvailable:  true,
			GasStationID: qazaqOil.ID,
		},
	}

	for _, pump := range qazaqOilPumps {
		err = db.Create(&pump).Error
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("NEW RECORD: (FuelPump) %+v\n", pump)
	}

	// Helios

	helios := model.GasStation{
		Company: "Helios",
		Address: "ул. Енбекшилер 18/1",
	}

	err = db.Create(&helios).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NEW RECORD: (GasStation) %+v\n", helios)

	heliosPumps := []model.FuelPump{
		{
			FuelType:     model.AI92,
			Price:        205,
			IsAvailable:  true,
			GasStationID: helios.ID,
		},
		{
			FuelType:     model.AI95,
			Price:        249,
			IsAvailable:  true,
			GasStationID: helios.ID,
		},
		{
			FuelType:     model.AI98,
			Price:        248,
			IsAvailable:  true,
			GasStationID: helios.ID,
		},
	}

	for _, pump := range heliosPumps {
		err = db.Create(&pump).Error
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("NEW RECORD: (FuelPump) %+v\n", pump)
	}

	// Газпромнефть

	gaspromneft := model.GasStation{
		Company: "Газпромнефть",
		Address: "просп. Богенбай батыра 24",
	}

	err = db.Create(&gaspromneft).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("NEW RECORD: (GasStation) %+v\n", gaspromneft)

	gaspromneftPumps := []model.FuelPump{
		{
			FuelType:     model.AI92,
			Price:        205,
			IsAvailable:  true,
			GasStationID: gaspromneft.ID,
		},
		{
			FuelType:     model.AI95,
			Price:        255,
			IsAvailable:  true,
			GasStationID: gaspromneft.ID,
		},
		{
			FuelType:     model.AI98,
			Price:        280,
			IsAvailable:  true,
			GasStationID: gaspromneft.ID,
		},
		{
			FuelType:     model.Diesel,
			Price:        295,
			IsAvailable:  true,
			GasStationID: gaspromneft.ID,
		},
	}

	for _, pump := range gaspromneftPumps {
		err = db.Create(&pump).Error
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("NEW RECORD: (FuelPump) %+v\n", pump)
	}
}
