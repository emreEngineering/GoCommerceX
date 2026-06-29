package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/order-service/internal/adapters"
	"GoCommerceX/order-service/internal/application"
	"GoCommerceX/order-service/internal/config"
	"GoCommerceX/order-service/internal/infrastructure"
	grpchandler "GoCommerceX/order-service/internal/transport/grpc"
	notificationv1 "GoCommerceX/proto/notification/v1"
	orderv1 "GoCommerceX/proto/order/v1"
	paymentv1 "GoCommerceX/proto/payment/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Order Service starting...")
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

	paymentConn, err := grpc.NewClient(cfg.PaymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	notificationConn, err := grpc.NewClient(cfg.NotificationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("unable to connect to notification service: %v", err)
	}
	defer notificationConn.Close()

	orderRepo := adapters.NewPostgresOrderRepository(pool)
	paymentClient := adapters.NewPaymentServiceClient(paymentv1.NewPaymentServiceClient(paymentConn))
	notificationClient := adapters.NewNotificationServiceClient(notificationv1.NewNotificationServiceClient(notificationConn))

	createOrderUseCase := application.NewCreateOrderUseCase(orderRepo, paymentClient, notificationClient)
	getOrderUseCase := application.NewGetOrderUseCase(orderRepo)
	cancelOrderUseCase := application.NewCancelOrderUseCase(orderRepo)

	orderHandler := grpchandler.NewOrderHandler(createOrderUseCase, getOrderUseCase, cancelOrderUseCase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	orderv1.RegisterOrderServiceServer(grpcServer, orderHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
