package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/inventory-service/internal/application"
	"GoCommerceX/inventory-service/internal/domain"
	inventoryv1 "GoCommerceX/proto/inventory/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InventoryHandler struct {
	inventoryv1.UnimplementedInventoryServiceServer
	createInventoryUseCase         *application.CreateInventoryUseCase
	getInventoryUseCase            *application.GetInventoryUseCase
	getInventoryByProductIDUseCase *application.GetInventoryByProductIDUseCase
	adjustStockUseCase             *application.AdjustStockUseCase
	reserveStockUseCase            *application.ReserveStockUseCase
	releaseStockUseCase            *application.ReleaseStockUseCase
	deleteInventoryUseCase         *application.DeleteInventoryUseCase
}

func NewInventoryHandler(
	createInventoryUseCase *application.CreateInventoryUseCase,
	getInventoryUseCase *application.GetInventoryUseCase,
	getInventoryByProductIDUseCase *application.GetInventoryByProductIDUseCase,
	adjustStockUseCase *application.AdjustStockUseCase,
	reserveStockUseCase *application.ReserveStockUseCase,
	releaseStockUseCase *application.ReleaseStockUseCase,
	deleteInventoryUseCase *application.DeleteInventoryUseCase,
) *InventoryHandler {
	return &InventoryHandler{
		createInventoryUseCase:         createInventoryUseCase,
		getInventoryUseCase:            getInventoryUseCase,
		getInventoryByProductIDUseCase: getInventoryByProductIDUseCase,
		adjustStockUseCase:             adjustStockUseCase,
		reserveStockUseCase:            reserveStockUseCase,
		releaseStockUseCase:            releaseStockUseCase,
		deleteInventoryUseCase:         deleteInventoryUseCase,
	}
}

func (h *InventoryHandler) CreateInventory(ctx context.Context, req *inventoryv1.CreateInventoryRequest) (*inventoryv1.CreateInventoryResponse, error) {
	output, err := h.createInventoryUseCase.Execute(ctx, application.CreateInventoryInput{
		ProductID:         req.GetProductId(),
		AvailableQuantity: req.GetAvailableQuantity(),
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreateInventoryProductIDRequired),
			errors.Is(err, application.ErrCreateInventoryQuantityInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("CreateInventory error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.CreateInventoryResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) GetInventory(ctx context.Context, req *inventoryv1.GetInventoryRequest) (*inventoryv1.GetInventoryResponse, error) {
	output, err := h.getInventoryUseCase.Execute(ctx, application.GetInventoryInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetInventoryIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetInventory error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.GetInventoryResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) GetInventoryByProductID(ctx context.Context, req *inventoryv1.GetInventoryByProductIDRequest) (*inventoryv1.GetInventoryByProductIDResponse, error) {
	output, err := h.getInventoryByProductIDUseCase.Execute(ctx, application.GetInventoryByProductIDInput{ProductID: req.GetProductId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetInventoryProductIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetInventoryByProductID error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.GetInventoryByProductIDResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) AdjustStock(ctx context.Context, req *inventoryv1.AdjustStockRequest) (*inventoryv1.AdjustStockResponse, error) {
	output, err := h.adjustStockUseCase.Execute(ctx, application.AdjustStockInput{ID: req.GetId(), Delta: req.GetDelta()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrAdjustStockIDRequired),
			errors.Is(err, application.ErrAdjustStockDeltaInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("AdjustStock error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.AdjustStockResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) ReserveStock(ctx context.Context, req *inventoryv1.ReserveStockRequest) (*inventoryv1.ReserveStockResponse, error) {
	output, err := h.reserveStockUseCase.Execute(ctx, application.ReserveStockInput{ID: req.GetId(), Quantity: req.GetQuantity()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrReserveStockIDRequired),
			errors.Is(err, application.ErrReserveStockQuantityInvalid),
			errors.Is(err, application.ErrAdjustStockDeltaInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("ReserveStock error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.ReserveStockResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) ReleaseStock(ctx context.Context, req *inventoryv1.ReleaseStockRequest) (*inventoryv1.ReleaseStockResponse, error) {
	output, err := h.releaseStockUseCase.Execute(ctx, application.ReleaseStockInput{ID: req.GetId(), Quantity: req.GetQuantity()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrReleaseStockIDRequired),
			errors.Is(err, application.ErrReleaseStockQuantityInvalid),
			errors.Is(err, application.ErrAdjustStockDeltaInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("ReleaseStock error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.ReleaseStockResponse{Inventory: toInventoryProto(output.Inventory)}, nil
}

func (h *InventoryHandler) DeleteInventory(ctx context.Context, req *inventoryv1.DeleteInventoryRequest) (*inventoryv1.DeleteInventoryResponse, error) {
	output, err := h.deleteInventoryUseCase.Execute(ctx, application.DeleteInventoryInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrDeleteInventoryIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrInventoryNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("DeleteInventory error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}
	return &inventoryv1.DeleteInventoryResponse{Success: output.Success}, nil
}

func toInventoryProto(inventory domain.Inventory) *inventoryv1.Inventory {
	return &inventoryv1.Inventory{
		Id:                inventory.ID,
		ProductId:         inventory.ProductID,
		AvailableQuantity: inventory.AvailableQuantity,
		ReservedQuantity:  inventory.ReservedQuantity,
		CreatedAt:         inventory.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         inventory.UpdatedAt.Format(time.RFC3339),
	}
}
