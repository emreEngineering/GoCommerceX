package ports

import (
	"GoCommerceX/cart-service/internal/domain"
	"context"
)

type CartRepository interface {
	Save(ctx context.Context, cart domain.Cart) error
	FindByUserID(ctx context.Context, userID string) (domain.Cart, error)
	Update(ctx context.Context, cart domain.Cart) error
	Delete(ctx context.Context, userID string) error
}
