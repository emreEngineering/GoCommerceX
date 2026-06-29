package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInventoryIDRequired        = errors.New("inventory: id is required")
	ErrInventoryProductIDRequired = errors.New("inventory: product id is required")
	ErrInventoryQuantityInvalid   = errors.New("inventory: quantity must be zero or greater")
	ErrInventoryReservedInvalid   = errors.New("inventory: reserved quantity must be zero or greater")
)

type Inventory struct {
	ID                string
	ProductID         string
	AvailableQuantity int32
	ReservedQuantity  int32
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func NewInventory(productID string, availableQuantity int32) Inventory {
	now := time.Now()
	return Inventory{
		ID:                uuid.NewString(),
		ProductID:         strings.TrimSpace(productID),
		AvailableQuantity: availableQuantity,
		ReservedQuantity:  0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}
}

func (i Inventory) Validate() error {
	if strings.TrimSpace(i.ID) == "" {
		return ErrInventoryIDRequired
	}
	if strings.TrimSpace(i.ProductID) == "" {
		return ErrInventoryProductIDRequired
	}
	if i.AvailableQuantity < 0 {
		return ErrInventoryQuantityInvalid
	}
	if i.ReservedQuantity < 0 {
		return ErrInventoryReservedInvalid
	}
	return nil
}
