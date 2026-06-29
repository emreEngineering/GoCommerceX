package ports

import (
	"GoCommerceX/notification-service/internal/domain"
	"context"
)

type NotificationRepository interface {
	Save(ctx context.Context, notification domain.Notification) error
	FindByID(ctx context.Context, id string) (domain.Notification, error)
	FindByOrderID(ctx context.Context, orderID string) (domain.Notification, error)
}
