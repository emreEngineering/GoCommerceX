package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/order-service/internal/domain"
	"GoCommerceX/order-service/internal/ports"
)

var (
	ErrGetOrderIDRequired    = errors.New("get order: id is required")
	ErrCancelOrderIDRequired = errors.New("cancel order: id is required")
	ErrOrderNotFound         = errors.New("order not found")
)

type GetOrderInput struct {
	ID string
}

type GetOrderOutput struct {
	Order domain.Order
}

type GetOrderUseCase struct {
	orderRepo ports.OrderRepository
}

func NewGetOrderUseCase(orderRepo ports.OrderRepository) *GetOrderUseCase {
	return &GetOrderUseCase{orderRepo: orderRepo}
}

func (uc *GetOrderUseCase) Execute(ctx context.Context, input GetOrderInput) (GetOrderOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetOrderOutput{}, ErrGetOrderIDRequired
	}

	order, err := uc.orderRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrOrderNotFound) {
			return GetOrderOutput{}, ErrOrderNotFound
		}
		return GetOrderOutput{}, err
	}

	return GetOrderOutput{Order: order}, nil
}

type CancelOrderInput struct {
	ID string
}

type CancelOrderOutput struct {
	Order domain.Order
}

type CancelOrderUseCase struct {
	orderRepo ports.OrderRepository
}

func NewCancelOrderUseCase(orderRepo ports.OrderRepository) *CancelOrderUseCase {
	return &CancelOrderUseCase{orderRepo: orderRepo}
}

func (uc *CancelOrderUseCase) Execute(ctx context.Context, input CancelOrderInput) (CancelOrderOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return CancelOrderOutput{}, ErrCancelOrderIDRequired
	}

	order, err := uc.orderRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrOrderNotFound) {
			return CancelOrderOutput{}, ErrOrderNotFound
		}
		return CancelOrderOutput{}, err
	}

	order.Status = "cancelled"
	if err := uc.orderRepo.Update(ctx, order); err != nil {
		return CancelOrderOutput{}, err
	}

	return CancelOrderOutput{Order: order}, nil
}
