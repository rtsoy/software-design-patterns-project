package model

type FuelType string

const (
	Gas    = "Газ"
	AI92   = "АИ-92"
	AI95   = "АИ-95"
	AI98   = "АИ-98"
	Diesel = "Дизель"
)

type FuelPump struct {
	ID           uint `gorm:"primaryKey"`
	FuelType     FuelType
	Price        uint
	IsAvailable  bool `gorm:"default:true"`
	GasStationID uint
	GasStation   GasStation `gorm:"foreignKey:GasStationID"`
}