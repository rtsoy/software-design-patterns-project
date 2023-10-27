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
	ID           uint64 `gorm:"primaryKey"`
	FuelType     FuelType
	Price        uint64
	IsAvailable  bool `gorm:"default:true"`
	TelegramID   int64
	GasStationID uint64
	GasStation   GasStation `gorm:"foreignKey:GasStationID"`
}
