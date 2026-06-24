package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/auth-service/internal/adapters"
	"GoCommerceX/auth-service/internal/application"
	"GoCommerceX/auth-service/internal/config"
	"GoCommerceX/auth-service/internal/infrastructure"
	grpchandler "GoCommerceX/auth-service/internal/transport/grpc"
	"GoCommerceX/proto/auth/v1"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Auth Service starting...")
	fmt.Printf("gRPC port: %s\n", cfg.GRPCPort)
	fmt.Printf("DB host: %s\n", cfg.DBHost)

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

	// Migration'ları çalıştır
	if err := infrastructure.RunMigrations(pool); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Bağımlılıkları oluştur
	userRepo := adapters.NewPostgresUserRepository(pool)
	passwordHasher := adapters.NewBcryptPasswordHasher()
	tokenGenerator := adapters.NewJWTTokenGenerator(cfg.JWTSecret)

	registerUseCase := application.NewRegisterUserUseCase(userRepo, passwordHasher)
	loginUseCase := application.NewLoginUserUseCase(userRepo, passwordHasher, tokenGenerator)

	// gRPC handler'ı oluştur
	authHandler := grpchandler.NewAuthHandler(registerUseCase, loginUseCase)

	// gRPC sunucusunu başlat
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
