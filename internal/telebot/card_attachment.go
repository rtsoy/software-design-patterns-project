package telebot

import (
	"log"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
	tele "gopkg.in/telebot.v3"
)

var (
	cards = map[int64]model.CreditCard{}
)

// CardHandler is an interface for handling different steps in the credit card processing flow.
type CardHandler interface {
	Handle(c tele.Context) (*tele.Message, error)
	SetNext(handler CardHandler)
	GetNext() CardHandler
}

// ConcreteCardHandler handles the step of entering the card number.
type ConcreteCardHandler struct {
	next CardHandler
	b    *Bot
}

func (h *ConcreteCardHandler) Handle(c tele.Context) (*tele.Message, error) {
	return h.b.bot.Send(c.Sender(), enterCardNumberMessage, h.b.getCancelMarkup())
}

func (h *ConcreteCardHandler) SetNext(handler CardHandler) {
	h.next = handler
}

func (h *ConcreteCardHandler) GetNext() CardHandler {
	return h.next
}

// CardNumberHandler handles the step of verifying and storing the card number.
type CardNumberHandler struct {
	next CardHandler
	b    *Bot
}

func (h *CardNumberHandler) Handle(c tele.Context) (*tele.Message, error) {
	// Verify the card number entered by the user.
	num, err := h.b.verifyCardNumber(c.Text())
	if err != nil {
		if _, err = h.b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return h.b.bot.Send(c.Sender(), enterCardNumberMessage)
	}

	// Store the user's credit card information with the card number.
	cards[c.Sender().ID] = model.CreditCard{Number: num, TelegramID: c.Sender().ID}

	// Prompt the user to enter the card expiration date.
	return h.b.bot.Send(c.Sender(), enterCardExpDateMessage)
}

func (h *CardNumberHandler) SetNext(handler CardHandler) {
	h.next = handler
}

func (h *CardNumberHandler) GetNext() CardHandler {
	return h.next
}

// CardExpDateHandler handles the step of verifying and storing the card expiration date.
type CardExpDateHandler struct {
	next CardHandler
	b    *Bot
}

func (h *CardExpDateHandler) Handle(c tele.Context) (*tele.Message, error) {
	// Verify the card expiration date entered by the user.
	exp, err := h.b.verifyExpDate(c.Text())
	if err != nil {
		if _, err = h.b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return h.b.bot.Send(c.Sender(), enterCardExpDateMessage)
	}

	// Update the expiration date in the user's credit card information.
	card := cards[c.Sender().ID]
	card.ExpirationDate = exp
	cards[c.Sender().ID] = card

	// Prompt the user to enter the cardholder's name.
	return h.b.bot.Send(c.Sender(), enterCardHolderMessage)
}

func (h *CardExpDateHandler) SetNext(handler CardHandler) {
	h.next = handler
}

func (h *CardExpDateHandler) GetNext() CardHandler {
	return h.next
}

// CardHolderHandler handles the step of verifying and storing the cardholder's name.
type CardHolderHandler struct {
	next CardHandler
	b    *Bot
}

func (h *CardHolderHandler) Handle(c tele.Context) (*tele.Message, error) {
	// Verify the cardholder's name entered by the user.
	holder, err := h.b.verifyCardHolder(c.Text())
	if err != nil {
		if _, err := h.b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return h.b.bot.Send(c.Sender(), enterCardHolderMessage)
	}

	// Update the cardholder's name in the user's credit card information.
	card := cards[c.Sender().ID]
	card.Cardholder = holder
	cards[c.Sender().ID] = card

	// Prompt the user to enter the CVV.
	return h.b.bot.Send(c.Sender(), enterCardCVVMessage)
}

func (h *CardHolderHandler) SetNext(handler CardHandler) {
	h.next = handler
}

func (h *CardHolderHandler) GetNext() CardHandler {
	return h.next
}

// CardCVVHandler handles the step of verifying and storing the CVV.
type CardCVVHandler struct {
	next CardHandler
	b    *Bot
}

func (h *CardCVVHandler) Handle(c tele.Context) (*tele.Message, error) {
	// Verify the CVV entered by the user.
	cvv, err := h.b.verifyCVV(c.Text())
	if err != nil {
		if _, err = h.b.bot.Send(c.Sender(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return h.b.bot.Send(c.Sender(), enterCardCVVMessage)
	}

	// Update the CVV in the user's credit card information.
	card := cards[c.Sender().ID]
	card.CVV = cvv
	cards[c.Sender().ID] = card

	// Store the card information in the repository.
	if err = h.b.repo.Create(card); err != nil {
		if _, err = h.b.bot.Send(c.Sender(), somethingWentWrongMessage); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return h.b.handleStart(c)
	}

	// Delete the stored card information and send a success message.
	delete(cards, c.Sender().ID)

	if _, err = h.b.bot.Send(c.Sender(), cardAttachmentSuccessMessage); err != nil {
		log.Printf("failed to send a message to the user: %v", err)
	}

	return h.b.handleStart(c)
}

func (h *CardCVVHandler) SetNext(handler CardHandler) {
	h.next = handler
}

func (h *CardCVVHandler) GetNext() CardHandler {
	return h.next
}
