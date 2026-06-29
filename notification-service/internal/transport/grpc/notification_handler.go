package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/notification-service/internal/application"
	"GoCommerceX/notification-service/internal/domain"
	notificationv1 "GoCommerceX/proto/notification/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationHandler struct {
	notificationv1.UnimplementedNotificationServiceServer
	sendNotificationUseCase *application.SendNotificationUseCase
	getNotificationUseCase  *application.GetNotificationUseCase
}

func NewNotificationHandler(sendNotificationUseCase *application.SendNotificationUseCase, getNotificationUseCase *application.GetNotificationUseCase) *NotificationHandler {
	return &NotificationHandler{
		sendNotificationUseCase: sendNotificationUseCase,
		getNotificationUseCase:  getNotificationUseCase,
	}
}

func (h *NotificationHandler) SendNotification(ctx context.Context, req *notificationv1.SendNotificationRequest) (*notificationv1.SendNotificationResponse, error) {
	output, err := h.sendNotificationUseCase.Execute(ctx, application.SendNotificationInput{
		OrderID: req.GetOrderId(),
		UserID:  req.GetUserId(),
		Type:    req.GetType(),
		Message: req.GetMessage(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrSendNotificationOrderIDRequired),
			errors.Is(err, application.ErrSendNotificationUserIDRequired),
			errors.Is(err, application.ErrSendNotificationTypeRequired),
			errors.Is(err, application.ErrSendNotificationMessageRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrNotificationAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("SendNotification error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &notificationv1.SendNotificationResponse{Notification: toNotificationProto(output.Notification)}, nil
}

func (h *NotificationHandler) GetNotification(ctx context.Context, req *notificationv1.GetNotificationRequest) (*notificationv1.GetNotificationResponse, error) {
	output, err := h.getNotificationUseCase.Execute(ctx, application.GetNotificationInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetNotificationIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrNotificationNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetNotification error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &notificationv1.GetNotificationResponse{Notification: toNotificationProto(output.Notification)}, nil
}

func toNotificationProto(notification domain.Notification) *notificationv1.Notification {
	return &notificationv1.Notification{
		Id:        notification.ID,
		OrderId:   notification.OrderID,
		UserId:    notification.UserID,
		Type:      notification.Type,
		Message:   notification.Message,
		Status:    notification.Status,
		CreatedAt: notification.CreatedAt.Format(time.RFC3339),
		UpdatedAt: notification.UpdatedAt.Format(time.RFC3339),
	}
}
