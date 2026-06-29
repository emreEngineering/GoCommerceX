package ports

import "context"

type PaymentRequest struct {
	OrderID string
	Amount  float64
	Method  string
}

type PaymentResult struct {
	PaymentID string
	Status    string
}

type PaymentProcessor interface {
	CreatePayment(ctx context.Context, req PaymentRequest) (PaymentResult, error)
}
