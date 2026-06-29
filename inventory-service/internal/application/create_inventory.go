package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/inventory-service/internal/domain"
	"GoCommerceX/inventory-service/internal/ports"
)

var (
	ErrCreateInventoryProductIDRequired = errors.New("create inventory: product id is required")
	ErrCreateInventoryQuantityInvalid   = errors.New("create inventory: quantity must be zero or greater")
	ErrInventoryAlreadyExists           = errors.New("create inventory: inventory already exists")
)

type CreateInventoryInput struct {
	ProductID         string
	AvailableQuantity int32
}

type CreateInventoryOutput struct {
	Inventory domain.Inventory
}

type CreateInventoryUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewCreateInventoryUseCase(inventoryRepo ports.InventoryRepository) *CreateInventoryUseCase {
	return &CreateInventoryUseCase{inventoryRepo: inventoryRepo}
}

func (uc *CreateInventoryUseCase) Execute(ctx context.Context, input CreateInventoryInput) (CreateInventoryOutput, error) {
	productID := strings.TrimSpace(input.ProductID)
	if productID == "" {
		return CreateInventoryOutput{}, ErrCreateInventoryProductIDRequired
	}
	if input.AvailableQuantity < 0 {
		return CreateInventoryOutput{}, ErrCreateInventoryQuantityInvalid
	}

	_, err := uc.inventoryRepo.FindByProductID(ctx, productID)
	if err == nil {
		return CreateInventoryOutput{}, ErrInventoryAlreadyExists
	}
	if !errors.Is(err, ports.ErrInventoryNotFound) {
		return CreateInventoryOutput{}, err
	}

	inventory := domain.NewInventory(productID, input.AvailableQuantity)
	if err := inventory.Validate(); err != nil {
		return CreateInventoryOutput{}, err
	}

	if err := uc.inventoryRepo.Save(ctx, inventory); err != nil {
		return CreateInventoryOutput{}, err
	}

	return CreateInventoryOutput{Inventory: inventory}, nil
}
