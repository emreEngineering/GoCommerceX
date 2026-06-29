package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/proto/user/v1"
	"GoCommerceX/user-service/internal/adapters"
	"GoCommerceX/user-service/internal/application"
	"GoCommerceX/user-service/internal/config"
	"GoCommerceX/user-service/internal/infrastructure"
	grpchandler "GoCommerceX/user-service/internal/transport/grpc"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("User Service starting...")
	fmt.Printf("gRPC port: %s\n", cfg.GRPCPort)

	// Veritabanı bağlantısı
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}
	fmt.Println("Connected to PostgreSQL")

	// Migration
	if err := infrastructure.RunMigrations(pool); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Bağımlılıklar
	userRepo := adapters.NewPostgresUserRepository(pool)

	createUserUseCase := application.NewCreateUserUseCase(userRepo)
	getUserUseCase := application.NewGetUserUseCase(userRepo)
	getUserByEmailUseCase := application.NewGetUserByEmailUseCase(userRepo)
	updateUserUseCase := application.NewUpdateUserUseCase(userRepo)
	deleteUserUseCase := application.NewDeleteUserUseCase(userRepo)

	// gRPC handler
	userHandler := grpchandler.NewUserHandler(
		createUserUseCase, getUserUseCase, getUserByEmailUseCase,
		updateUserUseCase, deleteUserUseCase,
	)

	// gRPC sunucusu
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userv1.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
