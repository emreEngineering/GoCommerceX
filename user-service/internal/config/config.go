package config

import "os"

type Config struct {
	GRPCPort string
	DBHost   string
	DBPort   string
	DBUser   string
	DBPassword string
	DBName   string
}

func Load() *Config {
	return &Config{
		GRPCPort:   getEnv("GRPC_PORT", "50052"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "gocommerce"),
		DBPassword: getEnv("DB_PASSWORD", "gocommerce_password"),
		DBName:     getEnv("DB_NAME", "gocommerce"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
