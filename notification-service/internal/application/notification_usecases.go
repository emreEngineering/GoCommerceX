package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/notification-service/internal/domain"
	"GoCommerceX/notification-service/internal/ports"
)

var (
	ErrSendNotificationOrderIDRequired = errors.New("send notification: order id is required")
	ErrSendNotificationUserIDRequired  = errors.New("send notification: user id is required")
	ErrSendNotificationTypeRequired    = errors.New("send notification: type is required")
	ErrSendNotificationMessageRequired = errors.New("send notification: message is required")
	ErrNotificationAlreadyExists       = errors.New("send notification: notification already exists")
	ErrGetNotificationIDRequired       = errors.New("get notification: id is required")
	ErrNotificationNotFound            = errors.New("notification not found")
)

type SendNotificationInput struct {
	OrderID string
	UserID  string
	Type    string
	Message string
}

type SendNotificationOutput struct {
	Notification domain.Notification
}

type SendNotificationUseCase struct {
	notificationRepo ports.NotificationRepository
}

func NewSendNotificationUseCase(notificationRepo ports.NotificationRepository) *SendNotificationUseCase {
	return &SendNotificationUseCase{notificationRepo: notificationRepo}
}

func (uc *SendNotificationUseCase) Execute(ctx context.Context, input SendNotificationInput) (SendNotificationOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	userID := strings.TrimSpace(input.UserID)
	typ := strings.TrimSpace(input.Type)
	message := strings.TrimSpace(input.Message)
	if orderID == "" {
		return SendNotificationOutput{}, ErrSendNotificationOrderIDRequired
	}
	if userID == "" {
		return SendNotificationOutput{}, ErrSendNotificationUserIDRequired
	}
	if typ == "" {
		return SendNotificationOutput{}, ErrSendNotificationTypeRequired
	}
	if message == "" {
		return SendNotificationOutput{}, ErrSendNotificationMessageRequired
	}

	_, err := uc.notificationRepo.FindByOrderID(ctx, orderID)
	if err == nil {
		return SendNotificationOutput{}, ErrNotificationAlreadyExists
	}
	if !errors.Is(err, ports.ErrNotificationNotFound) {
		return SendNotificationOutput{}, err
	}

	notification := domain.NewNotification(orderID, userID, typ, message)
	if err := notification.Validate(); err != nil {
		return SendNotificationOutput{}, err
	}
	if err := uc.notificationRepo.Save(ctx, notification); err != nil {
		return SendNotificationOutput{}, err
	}

	return SendNotificationOutput{Notification: notification}, nil
}

type GetNotificationInput struct {
	ID string
}

type GetNotificationOutput struct {
	Notification domain.Notification
}

type GetNotificationUseCase struct {
	notificationRepo ports.NotificationRepository
}

func NewGetNotificationUseCase(notificationRepo ports.NotificationRepository) *GetNotificationUseCase {
	return &GetNotificationUseCase{notificationRepo: notificationRepo}
}

func (uc *GetNotificationUseCase) Execute(ctx context.Context, input GetNotificationInput) (GetNotificationOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetNotificationOutput{}, ErrGetNotificationIDRequired
	}

	notification, err := uc.notificationRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrNotificationNotFound) {
			return GetNotificationOutput{}, ErrNotificationNotFound
		}
		return GetNotificationOutput{}, err
	}

	return GetNotificationOutput{Notification: notification}, nil
}
