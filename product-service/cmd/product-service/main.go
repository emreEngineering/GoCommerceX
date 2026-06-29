package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/product-service/internal/adapters"
	"GoCommerceX/product-service/internal/application"
	"GoCommerceX/product-service/internal/config"
	"GoCommerceX/product-service/internal/infrastructure"
	grpchandler "GoCommerceX/product-service/internal/transport/grpc"
	productv1 "GoCommerceX/proto/product/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Product Service starting...")
	fmt.Printf("gRPC port: %s\n", cfg.GRPCPort)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to ping database: %v", err)
	}
	fmt.Println("Connected to PostgreSQL")

	if err := infrastructure.RunMigrations(pool); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	productRepo := adapters.NewPostgresProductRepository(pool)
	createProductUseCase := application.NewCreateProductUseCase(productRepo)
	getProductUseCase := application.NewGetProductUseCase(productRepo)
	getProductBySKUUseCase := application.NewGetProductBySKUUseCase(productRepo)
	updateProductUseCase := application.NewUpdateProductUseCase(productRepo)
	deleteProductUseCase := application.NewDeleteProductUseCase(productRepo)

	productHandler := grpchandler.NewProductHandler(
		createProductUseCase,
		getProductUseCase,
		getProductBySKUUseCase,
		updateProductUseCase,
		deleteProductUseCase,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	productv1.RegisterProductServiceServer(grpcServer, productHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
