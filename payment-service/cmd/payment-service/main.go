package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/payment-service/internal/adapters"
	"GoCommerceX/payment-service/internal/application"
	"GoCommerceX/payment-service/internal/config"
	"GoCommerceX/payment-service/internal/infrastructure"
	grpchandler "GoCommerceX/payment-service/internal/transport/grpc"
	paymentv1 "GoCommerceX/proto/payment/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Payment Service starting...")
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

	paymentRepo := adapters.NewPostgresPaymentRepository(pool)
	createPaymentUseCase := application.NewCreatePaymentUseCase(paymentRepo)
	getPaymentUseCase := application.NewGetPaymentUseCase(paymentRepo)
	updatePaymentStatusUseCase := application.NewUpdatePaymentStatusUseCase(paymentRepo)

	paymentHandler := grpchandler.NewPaymentHandler(createPaymentUseCase, getPaymentUseCase, updatePaymentStatusUseCase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	paymentv1.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
