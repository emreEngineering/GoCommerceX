package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/product-service/internal/domain"
	"GoCommerceX/product-service/internal/ports"
)

var (
	ErrUpdateProductIDRequired   = errors.New("update product: id is required")
	ErrUpdateProductNameRequired = errors.New("update product: name is required")
	ErrUpdateProductPriceInvalid = errors.New("update product: price must be zero or greater")
	ErrUpdateProductStockInvalid = errors.New("update product: stock must be zero or greater")
	ErrDeleteProductIDRequired   = errors.New("delete product: id is required")
)

type UpdateProductInput struct {
	ID          string
	Name        string
	Description string
	Price       float64
	Stock       int32
}

type UpdateProductOutput struct {
	Product domain.Product
}

type UpdateProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewUpdateProductUseCase(productRepo ports.ProductRepository) *UpdateProductUseCase {
	return &UpdateProductUseCase{productRepo: productRepo}
}

func (uc *UpdateProductUseCase) Execute(ctx context.Context, input UpdateProductInput) (UpdateProductOutput, error) {
	id := strings.TrimSpace(input.ID)
	name := strings.TrimSpace(input.Name)
	description := strings.TrimSpace(input.Description)

	if id == "" {
		return UpdateProductOutput{}, ErrUpdateProductIDRequired
	}
	if name == "" {
		return UpdateProductOutput{}, ErrUpdateProductNameRequired
	}
	if input.Price < 0 {
		return UpdateProductOutput{}, ErrUpdateProductPriceInvalid
	}
	if input.Stock < 0 {
		return UpdateProductOutput{}, ErrUpdateProductStockInvalid
	}

	existingProduct, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrProductNotFound) {
			return UpdateProductOutput{}, ErrProductNotFound
		}
		return UpdateProductOutput{}, err
	}

	existingProduct.Name = name
	existingProduct.Description = description
	existingProduct.Price = input.Price
	existingProduct.Stock = input.Stock

	if err := uc.productRepo.Update(ctx, existingProduct); err != nil {
		return UpdateProductOutput{}, err
	}

	return UpdateProductOutput{Product: existingProduct}, nil
}

type DeleteProductInput struct {
	ID string
}

type DeleteProductOutput struct {
	Success bool
}

type DeleteProductUseCase struct {
	productRepo ports.ProductRepository
}

func NewDeleteProductUseCase(productRepo ports.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{productRepo: productRepo}
}

func (uc *DeleteProductUseCase) Execute(ctx context.Context, input DeleteProductInput) (DeleteProductOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return DeleteProductOutput{}, ErrDeleteProductIDRequired
	}

	if err := uc.productRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, ports.ErrProductNotFound) {
			return DeleteProductOutput{}, ErrProductNotFound
		}
		return DeleteProductOutput{}, err
	}

	return DeleteProductOutput{Success: true}, nil
}
