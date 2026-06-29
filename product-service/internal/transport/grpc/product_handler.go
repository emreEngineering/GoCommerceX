package grpc

import (
	"context"
	"errors"
	"log"
	"time"

	"GoCommerceX/product-service/internal/application"
	"GoCommerceX/product-service/internal/domain"
	productv1 "GoCommerceX/proto/product/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ProductHandler struct {
	productv1.UnimplementedProductServiceServer
	createProductUseCase   *application.CreateProductUseCase
	getProductUseCase      *application.GetProductUseCase
	getProductBySKUUseCase *application.GetProductBySKUUseCase
	updateProductUseCase   *application.UpdateProductUseCase
	deleteProductUseCase   *application.DeleteProductUseCase
}

func NewProductHandler(
	createProductUseCase *application.CreateProductUseCase,
	getProductUseCase *application.GetProductUseCase,
	getProductBySKUUseCase *application.GetProductBySKUUseCase,
	updateProductUseCase *application.UpdateProductUseCase,
	deleteProductUseCase *application.DeleteProductUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProductUseCase:   createProductUseCase,
		getProductUseCase:      getProductUseCase,
		getProductBySKUUseCase: getProductBySKUUseCase,
		updateProductUseCase:   updateProductUseCase,
		deleteProductUseCase:   deleteProductUseCase,
	}
}

func (h *ProductHandler) CreateProduct(ctx context.Context, req *productv1.CreateProductRequest) (*productv1.CreateProductResponse, error) {
	input := application.CreateProductInput{
		SKU:         req.GetSku(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       req.GetStockQuantity(),
	}

	output, err := h.createProductUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrCreateProductSKURequired),
			errors.Is(err, application.ErrCreateProductNameRequired),
			errors.Is(err, application.ErrCreateProductPriceInvalid),
			errors.Is(err, application.ErrCreateProductStockInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrProductAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			log.Printf("CreateProduct error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &productv1.CreateProductResponse{Product: toProductProto(output.Product)}, nil
}

func (h *ProductHandler) GetProduct(ctx context.Context, req *productv1.GetProductRequest) (*productv1.GetProductResponse, error) {
	output, err := h.getProductUseCase.Execute(ctx, application.GetProductInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetProductIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetProduct error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &productv1.GetProductResponse{Product: toProductProto(output.Product)}, nil
}

func (h *ProductHandler) GetProductBySKU(ctx context.Context, req *productv1.GetProductBySKURequest) (*productv1.GetProductBySKUResponse, error) {
	output, err := h.getProductBySKUUseCase.Execute(ctx, application.GetProductBySKUInput{SKU: req.GetSku()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrGetProductSKURequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("GetProductBySKU error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &productv1.GetProductBySKUResponse{Product: toProductProto(output.Product)}, nil
}

func (h *ProductHandler) UpdateProduct(ctx context.Context, req *productv1.UpdateProductRequest) (*productv1.UpdateProductResponse, error) {
	input := application.UpdateProductInput{
		ID:          req.GetId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		Price:       req.GetPrice(),
		Stock:       req.GetStockQuantity(),
	}

	output, err := h.updateProductUseCase.Execute(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrUpdateProductIDRequired),
			errors.Is(err, application.ErrUpdateProductNameRequired),
			errors.Is(err, application.ErrUpdateProductPriceInvalid),
			errors.Is(err, application.ErrUpdateProductStockInvalid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("UpdateProduct error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &productv1.UpdateProductResponse{Product: toProductProto(output.Product)}, nil
}

func (h *ProductHandler) DeleteProduct(ctx context.Context, req *productv1.DeleteProductRequest) (*productv1.DeleteProductResponse, error) {
	output, err := h.deleteProductUseCase.Execute(ctx, application.DeleteProductInput{ID: req.GetId()})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrDeleteProductIDRequired):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, application.ErrProductNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			log.Printf("DeleteProduct error: %v", err)
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &productv1.DeleteProductResponse{Success: output.Success}, nil
}

func toProductProto(product domain.Product) *productv1.Product {
	return &productv1.Product{
		Id:            product.ID,
		Sku:           product.SKU,
		Name:          product.Name,
		Description:   product.Description,
		Price:         product.Price,
		StockQuantity: product.Stock,
		CreatedAt:     product.CreatedAt.Format(time.RFC3339),
		UpdatedAt:     product.UpdatedAt.Format(time.RFC3339),
	}
}
