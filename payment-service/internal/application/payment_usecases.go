package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/payment-service/internal/domain"
	"GoCommerceX/payment-service/internal/ports"
)

var (
	ErrCreatePaymentOrderIDRequired = errors.New("create payment: order id is required")
	ErrCreatePaymentAmountInvalid   = errors.New("create payment: amount must be zero or greater")
	ErrPaymentAlreadyExists         = errors.New("create payment: payment already exists")
	ErrGetPaymentIDRequired         = errors.New("get payment: id is required")
	ErrPaymentNotFound              = errors.New("payment not found")
	ErrUpdatePaymentIDRequired      = errors.New("update payment: id is required")
	ErrUpdatePaymentStatusRequired  = errors.New("update payment: status is required")
)

type CreatePaymentInput struct {
	OrderID string
	Amount  float64
	Method  string
}

type CreatePaymentOutput struct {
	Payment domain.Payment
}

type CreatePaymentUseCase struct {
	paymentRepo ports.PaymentRepository
}

func NewCreatePaymentUseCase(paymentRepo ports.PaymentRepository) *CreatePaymentUseCase {
	return &CreatePaymentUseCase{paymentRepo: paymentRepo}
}

func (uc *CreatePaymentUseCase) Execute(ctx context.Context, input CreatePaymentInput) (CreatePaymentOutput, error) {
	orderID := strings.TrimSpace(input.OrderID)
	method := strings.TrimSpace(input.Method)
	if orderID == "" {
		return CreatePaymentOutput{}, ErrCreatePaymentOrderIDRequired
	}
	if input.Amount < 0 {
		return CreatePaymentOutput{}, ErrCreatePaymentAmountInvalid
	}

	_, err := uc.paymentRepo.FindByOrderID(ctx, orderID)
	if err == nil {
		return CreatePaymentOutput{}, ErrPaymentAlreadyExists
	}
	if !errors.Is(err, ports.ErrPaymentNotFound) {
		return CreatePaymentOutput{}, err
	}

	payment := domain.NewPayment(orderID, input.Amount, method)
	if err := payment.Validate(); err != nil {
		return CreatePaymentOutput{}, err
	}
	if err := uc.paymentRepo.Save(ctx, payment); err != nil {
		return CreatePaymentOutput{}, err
	}

	return CreatePaymentOutput{Payment: payment}, nil
}

type GetPaymentInput struct {
	ID string
}

type GetPaymentOutput struct {
	Payment domain.Payment
}

type GetPaymentUseCase struct {
	paymentRepo ports.PaymentRepository
}

func NewGetPaymentUseCase(paymentRepo ports.PaymentRepository) *GetPaymentUseCase {
	return &GetPaymentUseCase{paymentRepo: paymentRepo}
}

func (uc *GetPaymentUseCase) Execute(ctx context.Context, input GetPaymentInput) (GetPaymentOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetPaymentOutput{}, ErrGetPaymentIDRequired
	}

	payment, err := uc.paymentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrPaymentNotFound) {
			return GetPaymentOutput{}, ErrPaymentNotFound
		}
		return GetPaymentOutput{}, err
	}

	return GetPaymentOutput{Payment: payment}, nil
}

type UpdatePaymentStatusInput struct {
	ID     string
	Status string
}

type UpdatePaymentStatusOutput struct {
	Payment domain.Payment
}

type UpdatePaymentStatusUseCase struct {
	paymentRepo ports.PaymentRepository
}

func NewUpdatePaymentStatusUseCase(paymentRepo ports.PaymentRepository) *UpdatePaymentStatusUseCase {
	return &UpdatePaymentStatusUseCase{paymentRepo: paymentRepo}
}

func (uc *UpdatePaymentStatusUseCase) Execute(ctx context.Context, input UpdatePaymentStatusInput) (UpdatePaymentStatusOutput, error) {
	id := strings.TrimSpace(input.ID)
	status := strings.TrimSpace(input.Status)
	if id == "" {
		return UpdatePaymentStatusOutput{}, ErrUpdatePaymentIDRequired
	}
	if status == "" {
		return UpdatePaymentStatusOutput{}, ErrUpdatePaymentStatusRequired
	}

	payment, err := uc.paymentRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrPaymentNotFound) {
			return UpdatePaymentStatusOutput{}, ErrPaymentNotFound
		}
		return UpdatePaymentStatusOutput{}, err
	}

	payment.Status = status
	if err := uc.paymentRepo.Update(ctx, payment); err != nil {
		return UpdatePaymentStatusOutput{}, err
	}

	return UpdatePaymentStatusOutput{Payment: payment}, nil
}
