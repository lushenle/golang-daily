package main

import (
	"errors"
	"fmt"
)

// PaymentGatewayType defines the type of payment gateway
type PaymentGatewayType int

const (
	PayPalGatewayType PaymentGatewayType = iota
	StripeGatewayType
)

// PaymentGateway represents the common interface for payment gateways
type PaymentGateway interface {
	ProcessPayment(amount float64) error
}

// PayPalGateway is a concrete payment gateway
type PayPalGateway struct{}

func (pg *PayPalGateway) ProcessPayment(amount float64) error {
	fmt.Printf("Processing PayPal payment of $%.2f\n", amount)
	// Simulate PayPal payment processing logic
	return nil
}

// StripeGateway is another concrete payment gateway
type StripeGateway struct{}

func (sg *StripeGateway) ProcessPayment(amount float64) error {
	fmt.Printf("Processing Stripe payment of $%.2f\n", amount)
	// Simulate Stripe payment processing logic
	return nil
}

// NewPaymentGateway creates a payment gateway based on the provided type
func NewPaymentGateway(gwType PaymentGatewayType) (PaymentGateway, error) {
	switch gwType {
	case PayPalGatewayType:
		return &PayPalGateway{}, nil
	case StripeGatewayType:
		return &StripeGateway{}, nil
	default:
		return nil, errors.New("unsupported payment gateway type")
	}
}

func main() {
	payPalGateway, _ := NewPaymentGateway(PayPalGatewayType)
	payPalGateway.ProcessPayment(100.00)

	stripeGateway, _ := NewPaymentGateway(StripeGatewayType)
	stripeGateway.ProcessPayment(150.50)
}
