package model

type Observer struct {
	ID         uint `gorm:"primaryKey"`
	TelegramID uint
	FuelPumpID uint
	GasStation FuelPump `gorm:"foreignKey:FuelPumpID"`
}
