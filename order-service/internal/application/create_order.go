package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/order-service/internal/domain"
	"GoCommerceX/order-service/internal/ports"
)

var (
	ErrCreateOrderUserIDRequired = errors.New("create order: user id is required")
	ErrCreateOrderTotalInvalid   = errors.New("create order: total amount must be zero or greater")
	ErrOrderCreateFailed         = errors.New("create order: failed")
)

type CreateOrderInput struct {
	UserID      string
	TotalAmount float64
}

type CreateOrderOutput struct {
	Order domain.Order
}

type CreateOrderUseCase struct {
	orderRepo          ports.OrderRepository
	paymentProcessor   ports.PaymentProcessor
	notificationSender ports.NotificationSender
}

func NewCreateOrderUseCase(orderRepo ports.OrderRepository, paymentProcessor ports.PaymentProcessor, notificationSender ports.NotificationSender) *CreateOrderUseCase {
	return &CreateOrderUseCase{
		orderRepo:          orderRepo,
		paymentProcessor:   paymentProcessor,
		notificationSender: notificationSender,
	}
}

func (uc *CreateOrderUseCase) Execute(ctx context.Context, input CreateOrderInput) (CreateOrderOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return CreateOrderOutput{}, ErrCreateOrderUserIDRequired
	}
	if input.TotalAmount < 0 {
		return CreateOrderOutput{}, ErrCreateOrderTotalInvalid
	}

	order := domain.NewOrder(userID, input.TotalAmount)
	if err := uc.orderRepo.Save(ctx, order); err != nil {
		return CreateOrderOutput{}, err
	}

	paymentResult, err := uc.paymentProcessor.CreatePayment(ctx, ports.PaymentRequest{
		OrderID: order.ID,
		Amount:  order.TotalAmount,
		Method:  "card",
	})
	if err != nil {
		_ = uc.orderRepo.Delete(ctx, order.ID)
		return CreateOrderOutput{}, err
	}

	order.PaymentID = paymentResult.PaymentID
	order.Status = paymentResult.Status
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return CreateOrderOutput{}, err
	}

	if err := uc.notificationSender.Send(ctx, ports.NotificationRequest{
		OrderID: order.ID,
		UserID:  order.UserID,
		Type:    "order_created",
		Message: "Your order has been created successfully.",
	}); err != nil {
		return CreateOrderOutput{}, err
	}

	return CreateOrderOutput{Order: order}, nil
}
