package payment

import (
	"log"
	"strconv"

	"github.com/rtsoy/software-design-patterns-project/internal/model"
)

// CreditCard struct represents a credit card.
type CreditCard struct {
	creditCard *model.CreditCard
}

// NewCreditCard constructor creates a new CreditCard instance.
func NewCreditCard(creditCard *model.CreditCard) *CreditCard {
	return &CreditCard{
		creditCard: creditCard,
	}
}

// ProcessPayment processes the payment for the credit card.
// It logs the payment details to the console.
func (c *CreditCard) ProcessPayment(amount uint64) {
	str := strconv.FormatUint(c.creditCard.Number, 10)

	var number string

	for i := 0; i < len(str); i++ {
		if i%4 == 0 && i != 0 {
			number += " "
		}

		number += string(str[i])
	}

	log.Printf("got %d via credit card (%s)", amount, number)
}
