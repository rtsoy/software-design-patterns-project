package payment

// Strategy interface defines the behavior for processing payments.
type Strategy interface {
	// ProcessPayment processes the payment for the specified amount.
	ProcessPayment(amount uint64)
}
