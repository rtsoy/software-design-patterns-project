package telebot

import (
	"log"
	"regexp"
	"strconv"

	"github.com/rtsoy/software-design-patterns-project/internal/observer"
	tele "gopkg.in/telebot.v3"
)

// State is an interface that represents the state of order handling.
type State interface {
	HandleOrder(b *Bot, c tele.Context) (*tele.Message, error)
}

// OrderHandler handles orders and delegates the handling to the current state.
type OrderHandler struct {
	currentState State
}

// HandleOrder delegates the order handling to the current state.
func (h *OrderHandler) HandleOrder(b *Bot, c tele.Context) (*tele.Message, error) {
	return h.currentState.HandleOrder(b, c)
}

// SetCurrentState sets the current state of the OrderHandler.
func (h *OrderHandler) SetCurrentState(state State) {
	h.currentState = state
}

// AvailableOrder represents the state when an order is available.
type AvailableOrder struct{}

// HandleOrder handles an available order request.
func (h *AvailableOrder) HandleOrder(b *Bot, c tele.Context) (*tele.Message, error) {
	// Parse the pump price from the user's message.
	re := regexp.MustCompile(`(\d+)₸`)
	matches := re.FindAllStringSubmatch(c.Text(), -1)
	pumpPrice, _ := strconv.ParseUint(matches[0][1], 10, 64)

	orders[c.Sender().ID] = order{price: pumpPrice}

	return b.bot.Send(c.Sender(), insertFuelAmountMessage, b.getCancelMarkup())
}

// NotAvailableOrder represents the state when an order is not available.
type NotAvailableOrder struct {
	pumpID uint64
}

// HandleOrder handles a not available order request.
func (h *NotAvailableOrder) HandleOrder(b *Bot, c tele.Context) (*tele.Message, error) {
	// If the fuel pump is not available, notify the user and subscribe to updates.
	err := c.Send(fuelPumpIsNotAvailableMessage)
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return nil, err
	}

	// Create a subject for observing fuel pump updates.
	var obs observer.Subject = observer.NewFuelPump(b.repo.ObserverRepository, b.bot, h.pumpID)

	// Subscribe the user to receive updates from the fuel pump.
	err = obs.Subscribe(c.Sender().ID)

	if err != nil {
		log.Printf("failed to subcribe a user to fuel pump: %v", err)
	}

	// Redirect the user to the start command handler.
	return b.handleStart(c)
}
