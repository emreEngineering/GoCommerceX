package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/inventory-service/internal/domain"
	"GoCommerceX/inventory-service/internal/ports"
)

var (
	ErrGetInventoryIDRequired        = errors.New("get inventory: id is required")
	ErrGetInventoryProductIDRequired = errors.New("get inventory by product id: product id is required")
	ErrInventoryNotFound             = errors.New("inventory not found")
)

type GetInventoryInput struct {
	ID string
}

type GetInventoryOutput struct {
	Inventory domain.Inventory
}

type GetInventoryUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewGetInventoryUseCase(inventoryRepo ports.InventoryRepository) *GetInventoryUseCase {
	return &GetInventoryUseCase{inventoryRepo: inventoryRepo}
}

func (uc *GetInventoryUseCase) Execute(ctx context.Context, input GetInventoryInput) (GetInventoryOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return GetInventoryOutput{}, ErrGetInventoryIDRequired
	}

	inventory, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return GetInventoryOutput{}, ErrInventoryNotFound
		}
		return GetInventoryOutput{}, err
	}

	return GetInventoryOutput{Inventory: inventory}, nil
}

type GetInventoryByProductIDInput struct {
	ProductID string
}

type GetInventoryByProductIDOutput struct {
	Inventory domain.Inventory
}

type GetInventoryByProductIDUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewGetInventoryByProductIDUseCase(inventoryRepo ports.InventoryRepository) *GetInventoryByProductIDUseCase {
	return &GetInventoryByProductIDUseCase{inventoryRepo: inventoryRepo}
}

func (uc *GetInventoryByProductIDUseCase) Execute(ctx context.Context, input GetInventoryByProductIDInput) (GetInventoryByProductIDOutput, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return GetInventoryByProductIDOutput{}, ErrGetInventoryProductIDRequired
	}

	inventory, err := uc.inventoryRepo.FindByProductID(ctx, productID)
	if err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return GetInventoryByProductIDOutput{}, ErrInventoryNotFound
		}
		return GetInventoryByProductIDOutput{}, err
	}

	return GetInventoryByProductIDOutput{Inventory: inventory}, nil
}
