package telebot

import (
	"errors"
	"log"
	"strings"

	"github.com/rtsoy/software-design-patterns-project/internal/observer"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
	tele "gopkg.in/telebot.v3"
)

const (
	onStart  = "/start"
	onCancel = "⛔️ Отмена!"
)

// initHandlers initializes the command handlers for the bot.
func (b *Bot) initHandlers() {
	b.bot.Handle(onStart, b.handleStart)
	b.bot.Handle(onCancel, b.handleCancel)
	b.bot.Handle(tele.OnText, b.handleMessage)
}

// handleStart handles the "/start" command.
func (b *Bot) handleStart(c tele.Context) error {
	// Retrieve markup for available stations and send a start message.
	markup, err := b.getStationsMarkup()
	if err != nil {
		return c.Send(errorMessage)
	}

	return c.Send(startMessage, markup)
}

// handleCancel handles the cancellation command.
func (b *Bot) handleCancel(c tele.Context) error {
	// Release all fuel pumps reserved by the user.
	ids, err := b.repo.ReleaseAll(c.Sender().ID)
	if err != nil {
		log.Printf("failed to release all users' fuel pumps: %v", err)
	}

	// Notify subscribed users and unsubscribe them from fuel pump updates.
	for _, id := range ids {
		obs := observer.NewFuelPump(b.repo.ObserverRepository, b.bot, id)

		users, err := obs.NotifyAll()
		if err != nil {
			log.Printf("failed to notify users: %v", err)
		}

		for _, user := range users {
			err = obs.Unsubscribe(user)
			if err != nil {
				log.Printf("failed to unsubscribe a user from fuel pump: %v", err)
			}
		}
	}

	// Redirect the user to the start command handler.
	return b.handleStart(c)
}

// handleMessage handles user messages.
func (b *Bot) handleMessage(c tele.Context) error {
	switch {
	// ex. "1 | ⛽️ Gas Energy (просп. Кабанбай Батыра 45B)"
	case strings.Contains(c.Text(), "⛽️"):
		return b.handleFuelPumps(c)

	// ex. "1 | 🛢 Газ - 76₸ ✅"
	case strings.Contains(c.Text(), "🛢"):
		return b.handleOrder(c)

	// Unexpected response (command/message)
	default:
		return c.Send(unknownCommandMessage)
	}
}

// handleOrder handles the order process.
func (b *Bot) handleOrder(c tele.Context) error {
	// Extract fuel pump ID from the user's message.
	id, err := b.extractIDFromMessage(c.Text())
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return err
	}

	// Attempt to take or release a fuel pump for the user.
	err = b.repo.FuelPumpRepository.TakeOrRelease(c.Sender().ID, id)
	if err != nil { //nolint
		if errors.Is(err, repository.ErrForbidden) {
			// If the fuel pump is not available, notify the user and subscribe to updates.
			err = c.Send(fuelPumpIsNotAvailableMessage)
			if err != nil {
				log.Printf("msg: %s | %v", c.Text(), err) // ???

				return err
			}

			obs := observer.NewFuelPump(b.repo.ObserverRepository, b.bot, id)
			err = obs.Subscribe(c.Sender().ID)

			if err != nil {
				log.Printf("failed to subcribe a user to fuel pump: %v", err)
			}

			// Redirect the user to the start command handler.
			return b.handleStart(c)
		}

		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return err
	}

	// Provide the user with payment method options.
	markup := b.getPaymentMarkup()

	// TODO: Payment process
	return c.Send(choosePaymentMethodMessage, markup)
}

// handleFuelPumps handles the fuel pump selection process.
func (b *Bot) handleFuelPumps(c tele.Context) error {
	// Extract fuel pump ID from the user's message.
	id, err := b.extractIDFromMessage(c.Text())
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return err
	}

	// Retrieve markup for available fuel pumps at the selected station and send it to the user.
	markup, err := b.getFuelPumpsMarkup(id)
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // / ???

		return err
	}

	return c.Send(chooseFuelPumpMessage, markup)
}
