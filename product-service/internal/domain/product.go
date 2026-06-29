package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProductIDRequired    = errors.New("product: id is required")
	ErrProductSKURequired   = errors.New("product: sku is required")
	ErrProductNameRequired  = errors.New("product: name is required")
	ErrProductPriceInvalid  = errors.New("product: price must be zero or greater")
	ErrProductStockInvalid  = errors.New("product: stock must be zero or greater")
)

type Product struct {
	ID          string
	SKU         string
	Name        string
	Description string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProduct(sku, name, description string, price float64, stock int32) Product {
	now := time.Now()
	return Product{
		ID:          uuid.NewString(),
		SKU:         strings.TrimSpace(sku),
		Name:        strings.TrimSpace(name),
		Description: strings.TrimSpace(description),
		Price:       price,
		Stock:       stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

func (p Product) Validate() error {
	if p.ID == "" {
		return ErrProductIDRequired
	}
	if p.SKU == "" {
		return ErrProductSKURequired
	}
	if p.Name == "" {
		return ErrProductNameRequired
	}
	if p.Price < 0 {
		return ErrProductPriceInvalid
	}
	if p.Stock < 0 {
		return ErrProductStockInvalid
	}
	return nil
}
