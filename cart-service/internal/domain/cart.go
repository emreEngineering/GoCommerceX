package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrCartIDRequired      = errors.New("cart: id is required")
	ErrCartUserIDRequired  = errors.New("cart: user id is required")
	ErrCartItemRequired    = errors.New("cart: item is required")
	ErrCartQuantityInvalid = errors.New("cart: quantity must be greater than zero")
)

type CartItem struct {
	ProductID string
	Quantity  int32
}

type Cart struct {
	ID        string
	UserID    string
	Items     []CartItem
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCart(userID string) Cart {
	now := time.Now()
	return Cart{
		ID:        uuid.NewString(),
		UserID:    strings.TrimSpace(userID),
		Items:     []CartItem{},
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (c Cart) Validate() error {
	if strings.TrimSpace(c.ID) == "" {
		return ErrCartIDRequired
	}
	if strings.TrimSpace(c.UserID) == "" {
		return ErrCartUserIDRequired
	}
	return nil
}

func (c Cart) TotalQuantity() int32 {
	var total int32
	for _, item := range c.Items {
		total += item.Quantity
	}
	return total
}

func (c Cart) FindItemIndex(productID string) int {
	for i, item := range c.Items {
		if item.ProductID == productID {
			return i
		}
	}
	return -1
}

func (c *Cart) AddItem(productID string, quantity int32) error {
	productID = strings.TrimSpace(productID)
	if productID == "" {
		return ErrCartItemRequired
	}
	if quantity <= 0 {
		return ErrCartQuantityInvalid
	}

	idx := c.FindItemIndex(productID)
	if idx >= 0 {
		c.Items[idx].Quantity += quantity
		return nil
	}

	c.Items = append(c.Items, CartItem{ProductID: productID, Quantity: quantity})
	return nil
}

func (c *Cart) RemoveItem(productID string) bool {
	idx := c.FindItemIndex(strings.TrimSpace(productID))
	if idx < 0 {
		return false
	}
	c.Items = append(c.Items[:idx], c.Items[idx+1:]...)
	return true
}

func (c *Cart) Clear() {
	c.Items = []CartItem{}
}
