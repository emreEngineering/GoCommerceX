package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/payment-service/internal/application"
	"GoCommerceX/payment-service/internal/domain"
	paymentv1 "GoCommerceX/proto/payment/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentv1.UnimplementedPaymentServiceServer
	createPaymentUseCase       *application.CreatePaymentUseCase
	getPaymentUseCase          *application.GetPaymentUseCase
	updatePaymentStatusUseCase *application.UpdatePaymentStatusUseCase
}

func NewPaymentHandler(createPaymentUseCase *application.CreatePaymentUseCase, getPaymentUseCase *application.GetPaymentUseCase, updatePaymentStatusUseCase *application.UpdatePaymentStatusUseCase) *PaymentHandler {
	return &PaymentHandler{
		createPaymentUseCase:       createPaymentUseCase,
		getPaymentUseCase:          getPaymentUseCase,
		updatePaymentStatusUseCase: updatePaymentStatusUseCase,
	}
}

func (h *PaymentHandler) CreatePayment(ctx context.Context, req *paymentv1.CreatePaymentRequest) (*paymentv1.CreatePaymentResponse, error) {
	output, err := h.createPaymentUseCase.Execute(ctx, application.CreatePaymentInput{
		OrderID: req.GetOrderId(),
		Amount:  req.GetAmount(),
		Method:  req.GetMethod(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreatePaymentOrderIDRequired),
			errors.Is(err, application.ErrCreatePaymentAmountInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrPaymentAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("CreatePayment error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &paymentv1.CreatePaymentResponse{Payment: toPaymentProto(output.Payment)}, nil
}

func (h *PaymentHandler) GetPayment(ctx context.Context, req *paymentv1.GetPaymentRequest) (*paymentv1.GetPaymentResponse, error) {
	output, err := h.getPaymentUseCase.Execute(ctx, application.GetPaymentInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetPaymentIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrPaymentNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetPayment error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &paymentv1.GetPaymentResponse{Payment: toPaymentProto(output.Payment)}, nil
}

func (h *PaymentHandler) UpdatePaymentStatus(ctx context.Context, req *paymentv1.UpdatePaymentStatusRequest) (*paymentv1.UpdatePaymentStatusResponse, error) {
	output, err := h.updatePaymentStatusUseCase.Execute(ctx, application.UpdatePaymentStatusInput{
		ID:     req.GetId(),
		Status: req.GetStatus(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrUpdatePaymentIDRequired),
			errors.Is(err, application.ErrUpdatePaymentStatusRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrPaymentNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("UpdatePaymentStatus error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &paymentv1.UpdatePaymentStatusResponse{Payment: toPaymentProto(output.Payment)}, nil
}

func toPaymentProto(payment domain.Payment) *paymentv1.Payment {
	return &paymentv1.Payment{
		Id:        payment.ID,
		OrderId:   payment.OrderID,
		Amount:    payment.Amount,
		Method:    payment.Method,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Format(time.RFC3339),
		UpdatedAt: payment.UpdatedAt.Format(time.RFC3339),
	}
}
