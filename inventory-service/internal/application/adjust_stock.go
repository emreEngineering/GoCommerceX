package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/inventory-service/internal/domain"
	"GoCommerceX/inventory-service/internal/ports"
)

var (
	ErrAdjustStockIDRequired       = errors.New("adjust stock: id is required")
	ErrAdjustStockDeltaInvalid     = errors.New("adjust stock: delta cannot make quantity negative")
	ErrReserveStockIDRequired      = errors.New("reserve stock: id is required")
	ErrReserveStockQuantityInvalid = errors.New("reserve stock: quantity must be greater than zero")
	ErrReleaseStockIDRequired      = errors.New("release stock: id is required")
	ErrReleaseStockQuantityInvalid = errors.New("release stock: quantity must be greater than zero")
)

type AdjustStockInput struct {
	ID    string
	Delta int32
}

type AdjustStockOutput struct {
	Inventory domain.Inventory
}

type AdjustStockUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewAdjustStockUseCase(inventoryRepo ports.InventoryRepository) *AdjustStockUseCase {
	return &AdjustStockUseCase{inventoryRepo: inventoryRepo}
}

func (uc *AdjustStockUseCase) Execute(ctx context.Context, input AdjustStockInput) (AdjustStockOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return AdjustStockOutput{}, ErrAdjustStockIDRequired
	}

	inventory, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return AdjustStockOutput{}, ErrInventoryNotFound
		}
		return AdjustStockOutput{}, err
	}

	nextAvailable := inventory.AvailableQuantity + input.Delta
	if nextAvailable < 0 {
		return AdjustStockOutput{}, ErrAdjustStockDeltaInvalid
	}

	inventory.AvailableQuantity = nextAvailable
	if err := uc.inventoryRepo.Update(ctx, inventory); err != nil {
		return AdjustStockOutput{}, err
	}

	return AdjustStockOutput{Inventory: inventory}, nil
}

type ReserveStockInput struct {
	ID       string
	Quantity int32
}

type ReserveStockOutput struct {
	Inventory domain.Inventory
}

type ReserveStockUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewReserveStockUseCase(inventoryRepo ports.InventoryRepository) *ReserveStockUseCase {
	return &ReserveStockUseCase{inventoryRepo: inventoryRepo}
}

func (uc *ReserveStockUseCase) Execute(ctx context.Context, input ReserveStockInput) (ReserveStockOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return ReserveStockOutput{}, ErrReserveStockIDRequired
	}
	if input.Quantity <= 0 {
		return ReserveStockOutput{}, ErrReserveStockQuantityInvalid
	}

	inventory, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return ReserveStockOutput{}, ErrInventoryNotFound
		}
		return ReserveStockOutput{}, err
	}
	if inventory.AvailableQuantity < input.Quantity {
		return ReserveStockOutput{}, ErrAdjustStockDeltaInvalid
	}

	inventory.AvailableQuantity -= input.Quantity
	inventory.ReservedQuantity += input.Quantity
	if err := uc.inventoryRepo.Update(ctx, inventory); err != nil {
		return ReserveStockOutput{}, err
	}

	return ReserveStockOutput{Inventory: inventory}, nil
}

type ReleaseStockInput struct {
	ID       string
	Quantity int32
}

type ReleaseStockOutput struct {
	Inventory domain.Inventory
}

type ReleaseStockUseCase struct {
	inventoryRepo ports.InventoryRepository
}

func NewReleaseStockUseCase(inventoryRepo ports.InventoryRepository) *ReleaseStockUseCase {
	return &ReleaseStockUseCase{inventoryRepo: inventoryRepo}
}

func (uc *ReleaseStockUseCase) Execute(ctx context.Context, input ReleaseStockInput) (ReleaseStockOutput, error) {
	id := strings.TrimSpace(input.ID)
	if id == "" {
		return ReleaseStockOutput{}, ErrReleaseStockIDRequired
	}
	if input.Quantity <= 0 {
		return ReleaseStockOutput{}, ErrReleaseStockQuantityInvalid
	}

	inventory, err := uc.inventoryRepo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ports.ErrInventoryNotFound) {
			return ReleaseStockOutput{}, ErrInventoryNotFound
		}
		return ReleaseStockOutput{}, err
	}
	if inventory.ReservedQuantity < input.Quantity {
		return ReleaseStockOutput{}, ErrAdjustStockDeltaInvalid
	}

	inventory.AvailableQuantity += input.Quantity
	inventory.ReservedQuantity -= input.Quantity
	if err := uc.inventoryRepo.Update(ctx, inventory); err != nil {
		return ReleaseStockOutput{}, err
	}

	return ReleaseStockOutput{Inventory: inventory}, nil
}
