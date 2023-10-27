package model

import "time"

type CreditCard struct {
	ID             uint64 `gorm:"primaryKey"`
	TelegramID     int64
	Number         uint64 `gorm:"unique"`
	Cardholder     string
	ExpirationDate time.Time
	CVV            uint64
}
