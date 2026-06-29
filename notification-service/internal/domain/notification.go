package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotificationIDRequired      = errors.New("notification: id is required")
	ErrNotificationOrderIDRequired = errors.New("notification: order id is required")
	ErrNotificationUserIDRequired  = errors.New("notification: user id is required")
	ErrNotificationTypeRequired    = errors.New("notification: type is required")
	ErrNotificationMessageRequired = errors.New("notification: message is required")
	ErrNotificationStatusRequired  = errors.New("notification: status is required")
)

type Notification struct {
	ID        string
	OrderID   string
	UserID    string
	Type      string
	Message   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewNotification(orderID, userID, typ, message string) Notification {
	now := time.Now()
	return Notification{
		ID:        uuid.NewString(),
		OrderID:   strings.TrimSpace(orderID),
		UserID:    strings.TrimSpace(userID),
		Type:      strings.TrimSpace(typ),
		Message:   strings.TrimSpace(message),
		Status:    "sent",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (n Notification) Validate() error {
	if strings.TrimSpace(n.ID) == "" {
		return ErrNotificationIDRequired
	}
	if strings.TrimSpace(n.OrderID) == "" {
		return ErrNotificationOrderIDRequired
	}
	if strings.TrimSpace(n.UserID) == "" {
		return ErrNotificationUserIDRequired
	}
	if strings.TrimSpace(n.Type) == "" {
		return ErrNotificationTypeRequired
	}
	if strings.TrimSpace(n.Message) == "" {
		return ErrNotificationMessageRequired
	}
	if strings.TrimSpace(n.Status) == "" {
		return ErrNotificationStatusRequired
	}
	return nil
}
