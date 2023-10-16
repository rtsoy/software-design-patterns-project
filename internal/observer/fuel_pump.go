package observer

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
)

type FuelPump struct {
	repo   repository.ObserverRepository
	bot    *tgbotapi.BotAPI
	pumpID uint
}

func NewFuelPump(repo repository.ObserverRepository, bot *tgbotapi.BotAPI, pumpID uint) *FuelPump {
	return &FuelPump{
		repo:   repo,
		bot:    bot,
		pumpID: pumpID,
	}
}

func (f *FuelPump) Subscribe(telegramID uint) error {
	return f.repo.Add(telegramID, f.pumpID)
}

func (f *FuelPump) Unsubscribe(telegramID uint) error {
	return f.repo.Delete(telegramID, f.pumpID)
}

func (f *FuelPump) NotifyAll() ([]uint, error) {
	ids, err := f.repo.GetAll(f.pumpID)
	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		text := fmt.Sprintf("❗️ Колонка #%d освободилась! Скорее займите её!", f.pumpID)

		msg := tgbotapi.NewMessage(int64(id), text)

		_, err = f.bot.Send(msg)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
