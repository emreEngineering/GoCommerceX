package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/cart-service/internal/application"
	"GoCommerceX/cart-service/internal/domain"
	cartv1 "GoCommerceX/proto/cart/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CartHandler struct {
	cartv1.UnimplementedCartServiceServer
	createCartUseCase *application.CreateCartUseCase
	getCartUseCase    *application.GetCartUseCase
	addItemUseCase    *application.AddItemUseCase
	removeItemUseCase *application.RemoveItemUseCase
	clearCartUseCase  *application.ClearCartUseCase
	deleteCartUseCase *application.DeleteCartUseCase
}

func NewCartHandler(
	createCartUseCase *application.CreateCartUseCase,
	getCartUseCase *application.GetCartUseCase,
	addItemUseCase *application.AddItemUseCase,
	removeItemUseCase *application.RemoveItemUseCase,
	clearCartUseCase *application.ClearCartUseCase,
	deleteCartUseCase *application.DeleteCartUseCase,
) *CartHandler {
	return &CartHandler{
		createCartUseCase: createCartUseCase,
		getCartUseCase:    getCartUseCase,
		addItemUseCase:    addItemUseCase,
		removeItemUseCase: removeItemUseCase,
		clearCartUseCase:  clearCartUseCase,
		deleteCartUseCase: deleteCartUseCase,
	}
}

func (h *CartHandler) CreateCart(ctx context.Context, req *cartv1.CreateCartRequest) (*cartv1.CreateCartResponse, error) {
	output, err := h.createCartUseCase.Execute(ctx, application.CreateCartInput{UserID: req.GetUserId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreateCartUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("CreateCart error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.CreateCartResponse{Cart: toCartProto(output.Cart)}, nil
}

func (h *CartHandler) GetCart(ctx context.Context, req *cartv1.GetCartRequest) (*cartv1.GetCartResponse, error) {
	output, err := h.getCartUseCase.Execute(ctx, application.GetCartInput{UserID: req.GetUserId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetCartUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetCart error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.GetCartResponse{Cart: toCartProto(output.Cart)}, nil
}

func (h *CartHandler) AddItem(ctx context.Context, req *cartv1.AddItemRequest) (*cartv1.AddItemResponse, error) {
	output, err := h.addItemUseCase.Execute(ctx, application.AddItemInput{
		UserID:    req.GetUserId(),
		ProductID: req.GetProductId(),
		Quantity:  req.GetQuantity(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrAddItemUserIDRequired),
			errors.Is(err, application.ErrAddItemProductIDRequired),
			errors.Is(err, application.ErrAddItemQuantityInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("AddItem error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.AddItemResponse{Cart: toCartProto(output.Cart)}, nil
}

func (h *CartHandler) RemoveItem(ctx context.Context, req *cartv1.RemoveItemRequest) (*cartv1.RemoveItemResponse, error) {
	output, err := h.removeItemUseCase.Execute(ctx, application.RemoveItemInput{
		UserID:    req.GetUserId(),
		ProductID: req.GetProductId(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrRemoveItemUserIDRequired),
			errors.Is(err, application.ErrRemoveItemProductIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("RemoveItem error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.RemoveItemResponse{Cart: toCartProto(output.Cart)}, nil
}

func (h *CartHandler) ClearCart(ctx context.Context, req *cartv1.ClearCartRequest) (*cartv1.ClearCartResponse, error) {
	output, err := h.clearCartUseCase.Execute(ctx, application.ClearCartInput{UserID: req.GetUserId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrClearCartUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("ClearCart error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.ClearCartResponse{Cart: toCartProto(output.Cart)}, nil
}

func (h *CartHandler) DeleteCart(ctx context.Context, req *cartv1.DeleteCartRequest) (*cartv1.DeleteCartResponse, error) {
	output, err := h.deleteCartUseCase.Execute(ctx, application.DeleteCartInput{UserID: req.GetUserId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrDeleteCartUserIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrCartNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("DeleteCart error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &cartv1.DeleteCartResponse{Success: output.Success}, nil
}

func toCartProto(cart domain.Cart) *cartv1.Cart {
	items := make([]*cartv1.CartItem, 0, len(cart.Items))
	for _, item := range cart.Items {
		items = append(items, &cartv1.CartItem{
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	return &cartv1.Cart{
		Id:            cart.ID,
		UserId:        cart.UserID,
		Items:         items,
		TotalQuantity: cart.TotalQuantity(),
		CreatedAt:     cart.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     cart.UpdatedAt.Format(time.RFC3339),
	}
}
