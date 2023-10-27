package model

type Observer struct {
	ID         uint64 `gorm:"primaryKey"`
	TelegramID int64
	FuelPumpID uint64
	GasStation FuelPump `gorm:"foreignKey:FuelPumpID"`
}
