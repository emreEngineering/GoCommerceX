package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/inventory-service/internal/ports"
)

var ErrDeleteInventoryIDRequired = errors.New("delete inventory: id is required")

type DeleteInventoryInput struct {
	ID string
}

type DeleteInventoryOutput struct {
	Success bool
}

type DeleteInventoryUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewDeleteInventoryUseCase(inventoryRepo ports.InventoryRepository) *DeleteInventoryUseCase {
	return &DeleteInventoryUseCase{inventoryRepo: inventoryRepo}
}

func (uc *DeleteInventoryUseCase) Execute(ctx context.Context, input DeleteInventoryInput) (DeleteInventoryOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return DeleteInventoryOutput{}, ErrDeleteInventoryIDRequired
	}

	if err := uc.inventoryRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return DeleteInventoryOutput{}, ErrInventoryNotFound
		}
		return DeleteInventoryOutput{}, err
	}

	return DeleteInventoryOutput{Success: true}, nil
}
