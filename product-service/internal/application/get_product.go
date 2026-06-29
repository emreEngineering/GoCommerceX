package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/product-service/internal/domain"
	"GoCommerceX/product-service/internal/ports"
)

var (
	ErrGetProductIDRequired  = errors.New("get product: id is required")
	ErrGetProductSKURequired = errors.New("get product by sku: sku is required")
	ErrProductNotFound       = errors.New("product not found")
)

type GetProductInput struct {
	ID string
}

type GetProductOutput struct {
	Product domain.Product
}

type GetProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewGetProductUseCase(productRepo ports.ProductRepository) *GetProductUseCase {
	return &GetProductUseCase{productRepo: productRepo}
}

func (uc *GetProductUseCase) Execute(ctx context.Context, input GetProductInput) (GetProductOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetProductOutput{}, ErrGetProductIDRequired
	}

	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrProductNotFound) {
			return GetProductOutput{}, ErrProductNotFound
		}
		return GetProductOutput{}, err
	}

	return GetProductOutput{Product: product}, nil
}

type GetProductBySKUInput struct {
	SKU string
}

type GetProductBySKUOutput struct {
	Product domain.Product
}

type GetProductBySKUUseCase struct {
	productRepo ports.ProductRepository
}

func NewGetProductBySKUUseCase(productRepo ports.ProductRepository) *GetProductBySKUUseCase {
	return &GetProductBySKUUseCase{productRepo: productRepo}
}

func (uc *GetProductBySKUUseCase) Execute(ctx context.Context, input GetProductBySKUInput) (GetProductBySKUOutput, error) {
	sku := strings.TrimSpace(input.SKU)
	if sku == "" {
		return GetProductBySKUOutput{}, ErrGetProductSKURequired
	}

	product, err := uc.productRepo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, ports.ErrProductNotFound) {
			return GetProductBySKUOutput{}, ErrProductNotFound
		}
		return GetProductBySKUOutput{}, err
	}

	return GetProductBySKUOutput{Product: product}, nil
}
