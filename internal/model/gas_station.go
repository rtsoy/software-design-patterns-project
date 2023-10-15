package model

type GasStation struct {
	ID      uint `gorm:"primaryKey"`
	Company string
	Address string
}
