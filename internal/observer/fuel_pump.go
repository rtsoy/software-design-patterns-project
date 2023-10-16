package observer

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
)

// FuelPump represents a fuel pump observer that allows users to subscribe, unsubscribe,
// and receive notifications when the pump becomes available.
type FuelPump struct {
	repo   repository.ObserverRepository
	bot    *tgbotapi.BotAPI
	pumpID uint
}

// NewFuelPump creates a new FuelPump observer.
func NewFuelPump(repo repository.ObserverRepository, bot *tgbotapi.BotAPI, pumpID uint) *FuelPump {
	return &FuelPump{
		repo:   repo,
		bot:    bot,
		pumpID: pumpID,
	}
}

// Subscribe allows a user to subscribe to notifications for the fuel pump.
func (f *FuelPump) Subscribe(telegramID uint) error {
	return f.repo.Add(telegramID, f.pumpID)
}

// Unsubscribe allows a user to unsubscribe from notifications for the fuel pump.
func (f *FuelPump) Unsubscribe(telegramID uint) error {
	return f.repo.Delete(telegramID, f.pumpID)
}

// NotifyAll sends notifications to all subscribed users when the pump becomes available.
func (f *FuelPump) NotifyAll() ([]uint, error) {
	// Retrieve the list of subscribed users.
	ids, err := f.repo.GetAll(f.pumpID)
	if err != nil {
		return nil, err
	}

	// Send notifications to each subscribed user.
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
