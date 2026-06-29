package ports

import (
	"GoCommerceX/payment-service/internal/domain"
	"context"
)

type PaymentRepository interface {
	Save(ctx context.Context, payment domain.Payment) error
	FindByID(ctx context.Context, id string) (domain.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (domain.Payment, error)
	Update(ctx context.Context, payment domain.Payment) error
}
