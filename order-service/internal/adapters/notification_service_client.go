package adapters

import (
	"context"

	"GoCommerceX/order-service/internal/ports"
	notificationv1 "GoCommerceX/proto/notification/v1"
)

type NotificationServiceClient struct {
	client notificationv1.NotificationServiceClient
}

func NewNotificationServiceClient(client notificationv1.NotificationServiceClient) *NotificationServiceClient {
	return &NotificationServiceClient{client: client}
}

func (c *NotificationServiceClient) Send(ctx context.Context, req ports.NotificationRequest) error {
	_, err := c.client.SendNotification(ctx, &notificationv1.SendNotificationRequest{
		OrderId: req.OrderID,
		UserId:  req.UserID,
		Type:    req.Type,
		Message: req.Message,
	})
	return err
}
