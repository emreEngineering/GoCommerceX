package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrPaymentIDRequired      = errors.New("payment: id is required")
	ErrPaymentOrderIDRequired = errors.New("payment: order id is required")
	ErrPaymentAmountInvalid   = errors.New("payment: amount must be zero or greater")
	ErrPaymentStatusRequired  = errors.New("payment: status is required")
)

type Payment struct {
	ID        string
	OrderID   string
	Amount    float64
	Method    string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewPayment(orderID string, amount float64, method string) Payment {
	now := time.Now()
	return Payment{
		ID:        uuid.NewString(),
		OrderID:   strings.TrimSpace(orderID),
		Amount:    amount,
		Method:    strings.TrimSpace(method),
		Status:    "succeeded",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (p Payment) Validate() error {
	if strings.TrimSpace(p.ID) == "" {
		return ErrPaymentIDRequired
	}
	if strings.TrimSpace(p.OrderID) == "" {
		return ErrPaymentOrderIDRequired
	}
	if strings.TrimSpace(p.Status) == "" {
		return ErrPaymentStatusRequired
	}
	if p.Amount < 0 {
		return ErrPaymentAmountInvalid
	}
	return nil
}
