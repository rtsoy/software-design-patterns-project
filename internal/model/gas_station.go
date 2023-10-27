package model

type GasStation struct {
	ID      uint64 `gorm:"primaryKey"`
	Company string
	Address string
}
