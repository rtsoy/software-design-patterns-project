package repository

import (
	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"gorm.io/gorm"
)

type CreditCardRepositoryPostgres struct {
	db *gorm.DB
}

func NewCreditCardRepository(db *gorm.DB) *CreditCardRepositoryPostgres {
	return &CreditCardRepositoryPostgres{
		db: db,
	}
}

func (r *CreditCardRepositoryPostgres) GetAll(telegramID int64) ([]model.CreditCard, error) {
	var result []model.CreditCard

	err := r.db.Where("telegram_id = ?", telegramID).Find(&result).Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *CreditCardRepositoryPostgres) Create(card model.CreditCard) error {
	return r.db.Create(&card).Error
}
