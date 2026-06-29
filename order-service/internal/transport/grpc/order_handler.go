package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/order-service/internal/application"
	"GoCommerceX/order-service/internal/domain"
	orderv1 "GoCommerceX/proto/order/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrderHandler struct {
	orderv1.UnimplementedOrderServiceServer
	createOrderUseCase *application.CreateOrderUseCase
	getOrderUseCase    *application.GetOrderUseCase
	cancelOrderUseCase *application.CancelOrderUseCase
}

func NewOrderHandler(createOrderUseCase *application.CreateOrderUseCase, getOrderUseCase *application.GetOrderUseCase, cancelOrderUseCase *application.CancelOrderUseCase) *OrderHandler {
	return &OrderHandler{
		createOrderUseCase: createOrderUseCase,
		getOrderUseCase:    getOrderUseCase,
		cancelOrderUseCase: cancelOrderUseCase,
	}
}

func (h *OrderHandler) CreateOrder(ctx context.Context, req *orderv1.CreateOrderRequest) (*orderv1.CreateOrderResponse, error) {
	output, err := h.createOrderUseCase.Execute(ctx, application.CreateOrderInput{
		UserID:      req.GetUserId(),
		TotalAmount: req.GetTotalAmount(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreateOrderUserIDRequired),
			errors.Is(err, application.ErrCreateOrderTotalInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			log.Printf("CreateOrder error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &orderv1.CreateOrderResponse{Order: toOrderProto(output.Order)}, nil
}

func (h *OrderHandler) GetOrder(ctx context.Context, req *orderv1.GetOrderRequest) (*orderv1.GetOrderResponse, error) {
	output, err := h.getOrderUseCase.Execute(ctx, application.GetOrderInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetOrderIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrOrderNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetOrder error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &orderv1.GetOrderResponse{Order: toOrderProto(output.Order)}, nil
}

func (h *OrderHandler) CancelOrder(ctx context.Context, req *orderv1.CancelOrderRequest) (*orderv1.CancelOrderResponse, error) {
	output, err := h.cancelOrderUseCase.Execute(ctx, application.CancelOrderInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCancelOrderIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrOrderNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("CancelOrder error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &orderv1.CancelOrderResponse{Order: toOrderProto(output.Order)}, nil
}

func toOrderProto(order domain.Order) *orderv1.Order {
	return &orderv1.Order{
		Id:          order.ID,
		UserId:      order.UserID,
		PaymentId:   order.PaymentID,
		Status:      order.Status,
		TotalAmount: order.TotalAmount,
		CreatedAt:   order.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   order.UpdatedAt.Format(time.RFC3339),
	}
}
