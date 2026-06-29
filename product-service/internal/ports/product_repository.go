package ports

import (
	"GoCommerceX/product-service/internal/domain"
	"context"
)

type ProductRepository interface {
	Save(ctx context.Context, product domain.Product) error
	FindByID(ctx context.Context, id string) (domain.Product, error)
	FindBySKU(ctx context.Context, sku string) (domain.Product, error)
	Update(ctx context.Context, product domain.Product) error
	Delete(ctx context.Context, id string) error
}
