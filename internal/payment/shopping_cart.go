package payment

// A ShoppingCart struct represents a shopping cart. It has a paymentStrategy field,
// which is used to specify the strategy for processing payments.
type ShoppingCart struct {
	paymentStrategy Strategy
}

// The SetPaymentStrategy method sets the payment strategy for the shopping cart.
func (s *ShoppingCart) SetPaymentStrategy(strategy Strategy) {
	s.paymentStrategy = strategy
}

// The Checkout method checks out the shopping cart and processes the payment.
func (s *ShoppingCart) Checkout(amount uint64) {
	// Call the ProcessPayment method on the payment strategy.
	s.paymentStrategy.ProcessPayment(amount)
}
