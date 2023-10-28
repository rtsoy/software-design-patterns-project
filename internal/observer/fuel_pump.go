package observer

import (
	"fmt"

	"github.com/rtsoy/software-design-patterns-project/internal/repository"
	tele "gopkg.in/telebot.v3"
)

// FuelPump represents a fuel pump observer that allows users to subscribe, unsubscribe,
// and receive notifications when the pump becomes available.
type FuelPump struct {
	repo   repository.ObserverRepository
	bot    *tele.Bot
	pumpID uint64
}

// NewFuelPump creates a new FuelPump observer.
func NewFuelPump(repo repository.ObserverRepository, bot *tele.Bot, pumpID uint64) *FuelPump {
	return &FuelPump{
		repo:   repo,
		bot:    bot,
		pumpID: pumpID,
	}
}

// Subscribe allows a user to subscribe to notifications for the fuel pump.
func (f *FuelPump) Subscribe(telegramID int64) error {
	return f.repo.Add(telegramID, f.pumpID)
}

// Unsubscribe allows a user to unsubscribe from notifications for the fuel pump.
func (f *FuelPump) Unsubscribe(telegramID int64) error {
	return f.repo.Delete(telegramID, f.pumpID)
}

// NotifyAll sends notifications to all subscribed users when the pump becomes available.
func (f *FuelPump) NotifyAll() ([]int64, error) {
	// Retrieve the list of subscribed users.
	ids, err := f.repo.GetAll(f.pumpID)
	if err != nil {
		return nil, err
	}

	// Send notifications to each subscribed user.
	for _, id := range ids {
		text := fmt.Sprintf("❗️ Колонка #%d освободилась! Скорее займите её!", f.pumpID)

		_, err = f.bot.Send(&tele.User{ID: id}, text)
		if err != nil {
			return nil, err
		}
	}

	return ids, nil
}
