package telebot

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
	"github.com/rtsoy/software-design-patterns-project/internal/observer"
	"github.com/rtsoy/software-design-patterns-project/internal/payment"
	"github.com/rtsoy/software-design-patterns-project/internal/repository"
	tele "gopkg.in/telebot.v3"
)

const (
	onStart  = "/start"
	onCancel = "⛔️ Отмена!"
)

type order struct {
	price  uint64
	amount uint64
}

var (
	cards  = map[int64]model.CreditCard{}
	orders = map[int64]order{}
)

// initHandlers initializes the command handlers for the bot.
func (b *Bot) initHandlers() {
	b.bot.Handle(onStart, b.wrap(b.handleStart))
	b.bot.Handle(onCancel, b.wrap(b.handleCancel))
	b.bot.Handle(tele.OnText, b.wrap(b.handleMessage))
}

// handleStart handles the "/start" command.
func (b *Bot) handleStart(c tele.Context) (*tele.Message, error) {
	// Retrieve markup for available stations and send a start message.
	markup, err := b.getStationsMarkup()
	if err != nil {
		return b.bot.Send(c.Sender(), errorMessage)
	}

	return b.bot.Send(c.Sender(), startMessage, markup)
}

// handleCancel handles the cancellation command.
func (b *Bot) handleCancel(c tele.Context) (*tele.Message, error) {
	// Release all fuel pumps reserved by the user.
	ids, err := b.repo.ReleaseAll(c.Sender().ID)
	if err != nil {
		log.Printf("failed to release all users' fuel pumps: %v", err)
	}

	delete(orders, c.Sender().ID)

	// Notify subscribed users and unsubscribe them from fuel pump updates.
	for _, id := range ids {
		var obs observer.Subject = observer.NewFuelPump(b.repo.ObserverRepository, b.bot, id)

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
func (b *Bot) handleMessage(c tele.Context) (*tele.Message, error) { //nolint
	switch {
	// ex. "1 | ⛽️ Gas Energy (просп. Кабанбай Батыра 45B)"
	case strings.Contains(c.Text(), "⛽️"):
		return b.handleFuelPumps(c)

	// ex. "1 | 🛢 Газ - 76₸ ✅"
	case strings.Contains(c.Text(), "🛢"):
		return b.handleOrder(c)

	// ex. "💼 Привязать карту."
	case strings.Contains(c.Text(), "💼"):
		return b.handleCardLink(c)

	// Check if the last sent message was for entering card number
	case b.lastSentMessage.get(c.Sender().ID) == enterCardNumberMessage:
		return b.handleCardNumber(c)

	// Check if the last sent message was for entering card expiration date
	case b.lastSentMessage.get(c.Sender().ID) == enterCardExpDateMessage:
		return b.handleCardExpDate(c)

	// Check if the last sent message was for entering cardholder information
	case b.lastSentMessage.get(c.Sender().ID) == enterCardHolderMessage:
		return b.handleCardHolder(c)

	// Check if the last sent message was for entering card CVV
	case b.lastSentMessage.get(c.Sender().ID) == enterCardCVVMessage:
		return b.handleCVV(c)

	// Check if the last sent message was for inserting fuel amount
	case b.lastSentMessage.get(c.Sender().ID) == insertFuelAmountMessage:
		return b.handleFuelAmount(c)

	// ex. "💵 Наличный рассчет"
	case strings.Contains(c.Text(), "💵"):
		return b.handleCashPayment(c)

	// ex. "☑️ Готово!"
	case strings.Contains(c.Text(), "☑️"):
		return b.handlePaymentToTheCashier(c)

	// ex. "💳 Безналичный рассчет"
	case strings.Contains(c.Text(), "💳"):
		return b.handleCreditCardPayment(c)

	// Check if the last sent message was for choosing credit card
	case b.lastSentMessage.get(c.Sender().ID) == chooseCreditCard:
		return b.handleBankPayment(c)

	// ex. "▶️ Начать!"
	case strings.Contains(c.Text(), "▶️"):
		return b.handleRefueling(c)

	// ex. "🛑 Закончить заправку!"
	case strings.Contains(c.Text(), "🛑"):
		return b.handleCancel(c)

	// Unexpected response (command/message)
	default:
		return b.bot.Send(c.Sender(), unknownCommandMessage)
	}
}

// handleRefueling handles the user's request to stop refueling.
func (b *Bot) handleRefueling(c tele.Context) (*tele.Message, error) {
	return b.bot.Send(c.Sender(), stopRefueling, b.getStopMarkup())
}

// handleBankPayment handles the user's request to make a payment using a bank card.
func (b *Bot) handleBankPayment(c tele.Context) (*tele.Message, error) {
	// Get the user's credit cards.
	userCards, err := b.repo.CreditCardRepository.GetAll(c.Sender().ID)
	if err != nil {
		if _, err := b.bot.Send(c.Sender(), somethingWentWrongMessage); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.handleCancel(c)
	}

	// Extract the card ID from the user's input.
	re := regexp.MustCompile(`(\d+).`)
	matches := re.FindAllStringSubmatch(c.Text(), -1)
	cardID, _ := strconv.Atoi(matches[0][1])

	// Get the selected card.
	card := userCards[cardID-1]

	// Create a payment strategy using the selected card.
	payStrategy := payment.NewCreditCard(&card)

	// Get the user's order and calculate the total amount to be paid.
	userOrder := orders[c.Sender().ID]
	total := userOrder.price * userOrder.amount

	// Process the payment using the selected card.
	payStrategy.ProcessPayment(total)

	// Send a confirmation message to the user with the deducted amount.
	if _, err := b.bot.Send(c.Sender(), fmt.Sprintf("-%d₸", total)); err != nil {
		log.Printf("failed to send a message to the user: %v", err)
	}

	return b.bot.Send(c.Sender(), startRefueling, b.getStartMarkup())
}

// handleCreditCardPayment handles the user's request to pay using a credit card.
func (b *Bot) handleCreditCardPayment(c tele.Context) (*tele.Message, error) {
	// Get the user's credit cards.
	userCards, err := b.repo.CreditCardRepository.GetAll(c.Sender().ID)
	if err != nil {
		if _, err := b.bot.Send(c.Sender(), somethingWentWrongMessage); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.handleCancel(c)
	}

	// Check if the user has no attached credit cards and guide them to add one.
	if len(userCards) == 0 {
		if _, err := b.bot.Send(c.Sender(), noAttachedCreditCards); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		// Provide options for payment methods.
		return b.bot.Send(c.Sender(), choosePaymentMethodMessage, b.getPaymentMarkup())
	}

	// Provide a list of the user's credit cards for selection.
	return b.bot.Send(c.Sender(), chooseCreditCard, b.getCreditCardsMarkup(userCards))
}

// handlePaymentToTheCashier handles the user's request to make a payment to the cashier using cash.
func (b *Bot) handlePaymentToTheCashier(c tele.Context) (*tele.Message, error) {
	// Create a payment strategy for cash.
	payStrategy := payment.NewCash()

	// Get the user's order and calculate the total amount to be paid.
	userOrder := orders[c.Sender().ID]
	total := userOrder.price * userOrder.amount

	// Process the payment using cash.
	payStrategy.ProcessPayment(total)

	// Send a confirmation message to the user and return to the refueling process.
	return b.bot.Send(c.Sender(), startRefueling, b.getStartMarkup())
}

// handleCashPayment handles the user's request to pay with cash.
func (b *Bot) handleCashPayment(c tele.Context) (*tele.Message, error) {
	// Get the user's order and calculate the total amount to be paid.
	userOrder := orders[c.Sender().ID]
	total := userOrder.price * userOrder.amount

	// Provide the total amount and options to complete or cancel the payment.
	return c.Bot().Send(c.Sender(), b.getCashPaymentMessage(total), b.getDoneOrCancelMarkup())
}

// handleFuelAmount handle the user's input of the fuel amount.
func (b *Bot) handleFuelAmount(c tele.Context) (*tele.Message, error) {
	amount, err := strconv.ParseUint(c.Text(), 10, 64)
	if err != nil {
		if _, err := b.bot.Send(c.Sender(), somethingWentWrongMessage); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.bot.Send(c.Sender(), insertFuelAmountMessage)
	}

	// Update the user's order with the selected fuel amount.
	ord := orders[c.Sender().ID]
	ord.amount = amount
	orders[c.Sender().ID] = ord

	return b.bot.Send(c.Sender(), choosePaymentMethodMessage, b.getPaymentMarkup())
}

// handleCVV handles the user's entry of the CVV (Card Verification Value).
func (b *Bot) handleCVV(c tele.Context) (*tele.Message, error) {
	// Verify the CVV entered by the user.
	cvv, err := b.verifyCVV(c.Text())
	if err != nil {
		if _, err := b.bot.Send(c.Sender(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.bot.Send(c.Sender(), enterCardCVVMessage)
	}

	// Update the CVV in the user's credit card information.
	card := cards[c.Sender().ID]
	card.CVV = cvv
	cards[c.Sender().ID] = card

	// Store the card information in the repository.
	if err = b.repo.Create(card); err != nil {
		if _, err = b.bot.Send(c.Sender(), somethingWentWrongMessage); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.handleStart(c)
	}

	// Delete the stored card information and send a success message.
	delete(cards, c.Sender().ID)

	if _, err := b.bot.Send(c.Sender(), cardAttachmentSuccessMessage); err != nil {
		log.Printf("failed to send a message to the user: %v", err)
	}

	return b.handleStart(c)
}

// handleCardHolder handles the user's entry of the cardholder's name.
func (b *Bot) handleCardHolder(c tele.Context) (*tele.Message, error) {
	// Verify the cardholder's name entered by the user.
	holder, err := b.verifyCardHolder(c.Text())
	if err != nil {
		if _, err := b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.bot.Send(c.Sender(), enterCardHolderMessage)
	}

	// Update the cardholder's name in the user's credit card information.
	card := cards[c.Sender().ID]
	card.Cardholder = holder
	cards[c.Sender().ID] = card

	// Prompt the user to enter the CVV.
	return b.bot.Send(c.Sender(), enterCardCVVMessage)
}

// handleCardExpDate handles the user's entry of the card expiration date.
func (b *Bot) handleCardExpDate(c tele.Context) (*tele.Message, error) {
	// Verify the card expiration date entered by the user.
	exp, err := b.verifyExpDate(c.Text())
	if err != nil {
		if _, err := b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.bot.Send(c.Sender(), enterCardExpDateMessage)
	}

	// Update the expiration date in the user's credit card information.
	card := cards[c.Sender().ID]
	card.ExpirationDate = exp
	cards[c.Sender().ID] = card

	// Prompt the user to enter the cardholder's name.
	return b.bot.Send(c.Sender(), enterCardHolderMessage)
}

// handleCardNumber handles the user's entry of the card number.
func (b *Bot) handleCardNumber(c tele.Context) (*tele.Message, error) {
	// Verify the card number entered by the user.
	num, err := b.verifyCardNumber(c.Text())
	if err != nil {
		if _, err := b.bot.Send(c.Recipient(), err.Error()); err != nil {
			log.Printf("failed to send a message to the user: %v", err)
		}

		return b.bot.Send(c.Sender(), enterCardNumberMessage)
	}

	// Store the user's credit card information with the card number.
	cards[c.Sender().ID] = model.CreditCard{Number: num, TelegramID: c.Sender().ID}

	// Prompt the user to enter the card expiration date.
	return b.bot.Send(c.Sender(), enterCardExpDateMessage)
}

// handleCardLink handles the card linking process.
func (b *Bot) handleCardLink(c tele.Context) (*tele.Message, error) {
	// Prompt the user to enter the card number.
	return b.bot.Send(c.Sender(), enterCardNumberMessage, b.getCancelMarkup())
}

// handleOrder handles the order process.
func (b *Bot) handleOrder(c tele.Context) (*tele.Message, error) {
	// Extract fuel pump ID from the user's message.
	id, err := b.extractIDFromMessage(c.Text())
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return nil, err
	}

	// Create an OrderHandler to manage the order states.
	handler := &OrderHandler{}

	// Attempt to take or release a fuel pump for the user.
	err = b.repo.FuelPumpRepository.TakeOrRelease(c.Sender().ID, id)
	if err != nil {
		if errors.Is(err, repository.ErrForbidden) {
			// Set the current state to NotAvailableOrder when fuel pump is unavailable.
			handler.SetCurrentState(&NotAvailableOrder{pumpID: id})

			return handler.HandleOrder(b, c)
		}

		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return nil, err
	}

	// Set the current state to AvailableOrder when the pump is available.
	handler.SetCurrentState(&AvailableOrder{})

	return handler.HandleOrder(b, c)
}

// handleFuelPumps handles the fuel pump selection process.
func (b *Bot) handleFuelPumps(c tele.Context) (*tele.Message, error) {
	// Extract fuel pump ID from the user's message.
	id, err := b.extractIDFromMessage(c.Text())
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // ???

		return nil, err
	}

	// Retrieve markup for available fuel pumps at the selected station and send it to the user.
	markup, err := b.getFuelPumpsMarkup(id)
	if err != nil {
		log.Printf("msg: %s | %v", c.Text(), err) // / ???

		return nil, err
	}

	return b.bot.Send(c.Sender(), chooseFuelPumpMessage, markup)
}
