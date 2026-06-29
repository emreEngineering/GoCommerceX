package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"GoCommerceX/cart-service/internal/adapters"
	"GoCommerceX/cart-service/internal/application"
	"GoCommerceX/cart-service/internal/config"
	grpchandler "GoCommerceX/cart-service/internal/transport/grpc"
	cartv1 "GoCommerceX/proto/cart/v1"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	fmt.Println("Cart Service starting...")
	fmt.Printf("gRPC port: %s\n", cfg.GRPCPort)
	fmt.Printf("Redis addr: %s\n", cfg.RedisAddr)

	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	defer client.Close()

	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("unable to ping redis: %v", err)
	}
	fmt.Println("Connected to Redis")

	cartRepo := adapters.NewRedisCartRepository(client, cfg.CartKeyPrefix)
	createCartUseCase := application.NewCreateCartUseCase(cartRepo)
	getCartUseCase := application.NewGetCartUseCase(cartRepo)
	addItemUseCase := application.NewAddItemUseCase(cartRepo)
	removeItemUseCase := application.NewRemoveItemUseCase(cartRepo)
	clearCartUseCase := application.NewClearCartUseCase(cartRepo)
	deleteCartUseCase := application.NewDeleteCartUseCase(cartRepo)

	cartHandler := grpchandler.NewCartHandler(
		createCartUseCase,
		getCartUseCase,
		addItemUseCase,
		removeItemUseCase,
		clearCartUseCase,
		deleteCartUseCase,
	)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	cartv1.RegisterCartServiceServer(grpcServer, cartHandler)
	reflection.Register(grpcServer)

	fmt.Println("gRPC server listening on port", cfg.GRPCPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
