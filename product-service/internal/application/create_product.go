package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/product-service/internal/domain"
	"GoCommerceX/product-service/internal/ports"
)

var (
	ErrCreateProductSKURequired  = errors.New("create product: sku is required")
	ErrCreateProductNameRequired = errors.New("create product: name is required")
	ErrCreateProductPriceInvalid = errors.New("create product: price must be zero or greater")
	ErrCreateProductStockInvalid = errors.New("create product: stock must be zero or greater")
	ErrProductAlreadyExists      = errors.New("create product: product already exists")
)

type CreateProductInput struct {
	SKU         string
	Name        string
	Description string
	Price       float64
	Stock       int32
}

type CreateProductOutput struct {
	Product domain.Product
}

type CreateProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewCreateProductUseCase(productRepo ports.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{productRepo: productRepo}
}

func (uc *CreateProductUseCase) Execute(ctx context.Context, input CreateProductInput) (CreateProductOutput, error) {
	sku := strings.TrimSpace(input.SKU)
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	if sku == "" {
		return CreateProductOutput{}, ErrCreateProductSKURequired
	}
	if name == "" {
		return CreateProductOutput{}, ErrCreateProductNameRequired
	}
	if input.Price < 0 {
		return CreateProductOutput{}, ErrCreateProductPriceInvalid
	}
	if input.Stock < 0 {
		return CreateProductOutput{}, ErrCreateProductStockInvalid
	}

	_, err := uc.productRepo.FindBySKU(ctx, sku)
	if err == nil {
		return CreateProductOutput{}, ErrProductAlreadyExists
	}
	if !errors.Is(err, ports.ErrProductNotFound) {
		return CreateProductOutput{}, err
	}

	product := domain.NewProduct(sku, name, description, input.Price, input.Stock)
	if err := product.Validate(); err != nil {
		return CreateProductOutput{}, err
	}

	if err := uc.productRepo.Save(ctx, product); err != nil {
		return CreateProductOutput{}, err
	}

	return CreateProductOutput{Product: product}, nil
}
