package ports

import "context"

type NotificationRequest struct {
	OrderID string
	UserID  string
	Type    string
	Message string
}

type NotificationSender interface {
	Send(ctx context.Context, req NotificationRequest) error
}
