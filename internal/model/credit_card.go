package model

import "time"

type CreditCard struct {
	ID             uint `gorm:"primaryKey"`
	TelegramID     uint
	Number         uint `gorm:"unique"`
	Cardholder     string
	ExpirationDate time.Time
	CVV            uint
}
