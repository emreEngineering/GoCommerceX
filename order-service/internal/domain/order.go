package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrOrderIDRequired     = errors.New("order: id is required")
	ErrOrderUserIDRequired = errors.New("order: user id is required")
	ErrOrderTotalInvalid   = errors.New("order: total amount must be zero or greater")
	ErrOrderStatusRequired = errors.New("order: status is required")
)

type Order struct {
	ID          string
	UserID      string
	PaymentID   string
	Status      string
	TotalAmount float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewOrder(userID string, totalAmount float64) Order {
	now := time.Now()
	return Order{
		ID:          uuid.NewString(),
		UserID:      strings.TrimSpace(userID),
		Status:      "pending",
		TotalAmount: totalAmount,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (o Order) Validate() error {
	if strings.TrimSpace(o.ID) == "" {
		return ErrOrderIDRequired
	}
	if strings.TrimSpace(o.UserID) == "" {
		return ErrOrderUserIDRequired
	}
	if strings.TrimSpace(o.Status) == "" {
		return ErrOrderStatusRequired
	}
	if o.TotalAmount < 0 {
		return ErrOrderTotalInvalid
	}
	return nil
}
