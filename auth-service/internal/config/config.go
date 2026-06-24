package config

import (
	"os"
)

type Config struct {
	GRPCPort   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
}

func Load() *Config {
	return &Config{
		GRPCPort:   getEnv("GRPC_PORT", "50051"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "gocommerce"),
		DBPassword: getEnv("DB_PASSWORD", "gocommerce_password"),
		DBName:     getEnv("DB_NAME", "gocommerce"),
		JWTSecret:  getEnv("JWT_SECRET", "change-me-in-production"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
