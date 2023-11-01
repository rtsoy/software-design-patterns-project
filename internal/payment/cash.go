package payment

import (
	"log"
)

// Cash struct represents a cash payment strategy.
type Cash struct {
}

// NewCash constructor creates a new Cash instance.
func NewCash() *Cash {
	return &Cash{}
}

// ProcessPayment processes the payment for the specified amount using cash.
// It logs the payment details to the console.
func (c *Cash) ProcessPayment(amount uint64) {
	log.Printf("got %d in cash", amount)
}
