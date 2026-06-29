package ports

import (
	"GoCommerceX/order-service/internal/domain"
	"context"
)

type OrderRepository interface {
	Save(ctx context.Context, order domain.Order) error
	FindByID(ctx context.Context, id string) (domain.Order, error)
	Update(ctx context.Context, order domain.Order) error
	Delete(ctx context.Context, id string) error
}
