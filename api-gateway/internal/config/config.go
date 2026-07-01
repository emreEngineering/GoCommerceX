package config

import "os"

type Config struct {
	HTTPPort      string
	AuthAddr      string
	UserAddr      string
	ProductAddr   string
	InventoryAddr string
	CartAddr      string
	OrderAddr     string
	PaymentAddr   string
	NotifAddr     string
	JWTSecret     string
}

func Load() *Config {
	return &Config{
		HTTPPort:      getEnv("HTTP_PORT", "8080"),
		AuthAddr:      getEnv("AUTH_ADDR", "localhost:50051"),
		UserAddr:      getEnv("USER_ADDR", "localhost:50052"),
		ProductAddr:   getEnv("PRODUCT_ADDR", "localhost:50053"),
		InventoryAddr: getEnv("INVENTORY_ADDR", "localhost:50054"),
		CartAddr:      getEnv("CART_ADDR", "localhost:50055"),
		OrderAddr:     getEnv("ORDER_ADDR", "localhost:50056"),
		PaymentAddr:   getEnv("PAYMENT_ADDR", "localhost:50057"),
		NotifAddr:     getEnv("NOTIF_ADDR", "localhost:50058"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
