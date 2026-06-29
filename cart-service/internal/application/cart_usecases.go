package application

import (
	"context"
	"errors"
	"strings"

	"GoCommerceX/cart-service/internal/domain"
	"GoCommerceX/cart-service/internal/ports"
)

var (
	ErrCreateCartUserIDRequired    = errors.New("create cart: user id is required")
	ErrCartAlreadyExists           = errors.New("create cart: cart already exists")
	ErrGetCartUserIDRequired       = errors.New("get cart: user id is required")
	ErrCartNotFound                = errors.New("cart not found")
	ErrAddItemUserIDRequired       = errors.New("add item: user id is required")
	ErrAddItemProductIDRequired    = errors.New("add item: product id is required")
	ErrAddItemQuantityInvalid      = errors.New("add item: quantity must be greater than zero")
	ErrRemoveItemUserIDRequired    = errors.New("remove item: user id is required")
	ErrRemoveItemProductIDRequired = errors.New("remove item: product id is required")
	ErrClearCartUserIDRequired     = errors.New("clear cart: user id is required")
	ErrDeleteCartUserIDRequired    = errors.New("delete cart: user id is required")
)

type CreateCartInput struct {
	UserID string
}

type CreateCartOutput struct {
	Cart domain.Cart
}

type CreateCartUseCase struct {
	cartRepo ports.CartRepository
}

func NewCreateCartUseCase(cartRepo ports.CartRepository) *CreateCartUseCase {
	return &CreateCartUseCase{cartRepo: cartRepo}
}

func (uc *CreateCartUseCase) Execute(ctx context.Context, input CreateCartInput) (CreateCartOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return CreateCartOutput{}, ErrCreateCartUserIDRequired
	}

	_, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err == nil {
		return CreateCartOutput{}, ErrCartAlreadyExists
	}
	if !errors.Is(err, ports.ErrCartNotFound) {
		return CreateCartOutput{}, err
	}

	cart := domain.NewCart(userID)
	if err := cart.Validate(); err != nil {
		return CreateCartOutput{}, err
	}
	if err := uc.cartRepo.Save(ctx, cart); err != nil {
		return CreateCartOutput{}, err
	}

	return CreateCartOutput{Cart: cart}, nil
}

type GetCartInput struct {
	UserID string
}

type GetCartOutput struct {
	Cart domain.Cart
}

type GetCartUseCase struct {
	cartRepo ports.CartRepository
}

func NewGetCartUseCase(cartRepo ports.CartRepository) *GetCartUseCase {
	return &GetCartUseCase{cartRepo: cartRepo}
}

func (uc *GetCartUseCase) Execute(ctx context.Context, input GetCartInput) (GetCartOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return GetCartOutput{}, ErrGetCartUserIDRequired
	}

	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			return GetCartOutput{}, ErrCartNotFound
		}
		return GetCartOutput{}, err
	}

	return GetCartOutput{Cart: cart}, nil
}

type AddItemInput struct {
	UserID    string
	ProductID string
	Quantity  int32
}

type AddItemOutput struct {
	Cart domain.Cart
}

type AddItemUseCase struct {
	cartRepo ports.CartRepository
}

func NewAddItemUseCase(cartRepo ports.CartRepository) *AddItemUseCase {
	return &AddItemUseCase{cartRepo: cartRepo}
}

func (uc *AddItemUseCase) Execute(ctx context.Context, input AddItemInput) (AddItemOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	productID := strings.TrimSpace(input.ProductID)
	if userID == "" {
		return AddItemOutput{}, ErrAddItemUserIDRequired
	}
	if productID == "" {
		return AddItemOutput{}, ErrAddItemProductIDRequired
	}
	if input.Quantity <= 0 {
		return AddItemOutput{}, ErrAddItemQuantityInvalid
	}

	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			cart = domain.NewCart(userID)
		} else {
			return AddItemOutput{}, err
		}
	}

	if err := cart.AddItem(productID, input.Quantity); err != nil {
		return AddItemOutput{}, err
	}
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			if err := uc.cartRepo.Save(ctx, cart); err != nil {
				return AddItemOutput{}, err
			}
			return AddItemOutput{Cart: cart}, nil
		}
		return AddItemOutput{}, err
	}

	return AddItemOutput{Cart: cart}, nil
}

type RemoveItemInput struct {
	UserID    string
	ProductID string
}

type RemoveItemOutput struct {
	Cart domain.Cart
}

type RemoveItemUseCase struct {
	cartRepo ports.CartRepository
}

func NewRemoveItemUseCase(cartRepo ports.CartRepository) *RemoveItemUseCase {
	return &RemoveItemUseCase{cartRepo: cartRepo}
}

func (uc *RemoveItemUseCase) Execute(ctx context.Context, input RemoveItemInput) (RemoveItemOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	productID := strings.TrimSpace(input.ProductID)
	if userID == "" {
		return RemoveItemOutput{}, ErrRemoveItemUserIDRequired
	}
	if productID == "" {
		return RemoveItemOutput{}, ErrRemoveItemProductIDRequired
	}

	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			return RemoveItemOutput{}, ErrCartNotFound
		}
		return RemoveItemOutput{}, err
	}

	if !cart.RemoveItem(productID) {
		return RemoveItemOutput{}, ErrCartNotFound
	}
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return RemoveItemOutput{}, err
	}

	return RemoveItemOutput{Cart: cart}, nil
}

type ClearCartInput struct {
	UserID string
}

type ClearCartOutput struct {
	Cart domain.Cart
}

type ClearCartUseCase struct {
	cartRepo ports.CartRepository
}

func NewClearCartUseCase(cartRepo ports.CartRepository) *ClearCartUseCase {
	return &ClearCartUseCase{cartRepo: cartRepo}
}

func (uc *ClearCartUseCase) Execute(ctx context.Context, input ClearCartInput) (ClearCartOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return ClearCartOutput{}, ErrClearCartUserIDRequired
	}

	cart, err := uc.cartRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			return ClearCartOutput{}, ErrCartNotFound
		}
		return ClearCartOutput{}, err
	}

	cart.Clear()
	if err := uc.cartRepo.Update(ctx, cart); err != nil {
		return ClearCartOutput{}, err
	}

	return ClearCartOutput{Cart: cart}, nil
}

type DeleteCartInput struct {
	UserID string
}

type DeleteCartOutput struct {
	Success bool
}

type DeleteCartUseCase struct {
	cartRepo ports.CartRepository
}

func NewDeleteCartUseCase(cartRepo ports.CartRepository) *DeleteCartUseCase {
	return &DeleteCartUseCase{cartRepo: cartRepo}
}

func (uc *DeleteCartUseCase) Execute(ctx context.Context, input DeleteCartInput) (DeleteCartOutput, error) {
	userID := strings.TrimSpace(input.UserID)
	if userID == "" {
		return DeleteCartOutput{}, ErrDeleteCartUserIDRequired
	}

	if err := uc.cartRepo.Delete(ctx, userID); err != nil {
		if errors.Is(err, ports.ErrCartNotFound) {
			return DeleteCartOutput{}, ErrCartNotFound
		}
		return DeleteCartOutput{}, err
	}

	return DeleteCartOutput{Success: true}, nil
}
