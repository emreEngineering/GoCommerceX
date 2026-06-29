package ports

import (
	"GoCommerceX/inventory-service/internal/domain"
	"context"
)

type InventoryRepository interface {
	Save(ctx context.Context, inventory domain.Inventory) error
	FindByID(ctx context.Context, id string) (domain.Inventory, error)
	FindByProductID(ctx context.Context, productID string) (domain.Inventory, error)
	Update(ctx context.Context, inventory domain.Inventory) error
	Delete(ctx context.Context, id string) error
}
