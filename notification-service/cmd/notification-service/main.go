package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/notification-service/internal/adapters"
	"GoCommerceX/notification-service/internal/application"
	"GoCommerceX/notification-service/internal/config"
	"GoCommerceX/notification-service/internal/infrastructure"
	grpchandler "GoCommerceX/notification-service/internal/transport/grpc"
	notificationv1 "GoCommerceX/proto/notification/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Notification Service starting...")
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

	notificationRepo := adapters.NewPostgresNotificationRepository(pool)
	sendNotificationUseCase := application.NewSendNotificationUseCase(notificationRepo)
	getNotificationUseCase := application.NewGetNotificationUseCase(notificationRepo)

	notificationHandler := grpchandler.NewNotificationHandler(sendNotificationUseCase, getNotificationUseCase)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	notificationv1.RegisterNotificationServiceServer(grpcServer, notificationHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
