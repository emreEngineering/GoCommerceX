package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/inventory-service/internal/adapters"
	"GoCommerceX/inventory-service/internal/application"
	"GoCommerceX/inventory-service/internal/config"
	"GoCommerceX/inventory-service/internal/infrastructure"
	grpchandler "GoCommerceX/inventory-service/internal/transport/grpc"
	inventoryv1 "GoCommerceX/proto/inventory/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Inventory Service starting...")
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

	inventoryRepo := adapters.NewPostgresInventoryRepository(pool)
	createInventoryUseCase := application.NewCreateInventoryUseCase(inventoryRepo)
	getInventoryUseCase := application.NewGetInventoryUseCase(inventoryRepo)
	getInventoryByProductIDUseCase := application.NewGetInventoryByProductIDUseCase(inventoryRepo)
	adjustStockUseCase := application.NewAdjustStockUseCase(inventoryRepo)
	reserveStockUseCase := application.NewReserveStockUseCase(inventoryRepo)
	releaseStockUseCase := application.NewReleaseStockUseCase(inventoryRepo)
	deleteInventoryUseCase := application.NewDeleteInventoryUseCase(inventoryRepo)

	inventoryHandler := grpchandler.NewInventoryHandler(
		createInventoryUseCase,
		getInventoryUseCase,
		getInventoryByProductIDUseCase,
		adjustStockUseCase,
		reserveStockUseCase,
		releaseStockUseCase,
		deleteInventoryUseCase,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	inventoryv1.RegisterInventoryServiceServer(grpcServer, inventoryHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
