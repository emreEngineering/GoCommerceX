package config

import "os"

type Config struct {
	GRPCPort                string
	DBHost                  string
	DBPort                  string
	DBUser                  string
	DBPassword              string
	DBName                  string
	PaymentServiceAddr      string
	NotificationServiceAddr string
}

func Load() *Config {
	return &Config{
		GRPCPort:                getEnv("GRPC_PORT", "50056"),
		DBHost:                  getEnv("DB_HOST", "localhost"),
		DBPort:                  getEnv("DB_PORT", "5432"),
		DBUser:                  getEnv("DB_USER", "gocommerce"),
		DBPassword:              getEnv("DB_PASSWORD", "gocommerce_password"),
		DBName:                  getEnv("DB_NAME", "gocommerce"),
		PaymentServiceAddr:      getEnv("PAYMENT_SERVICE_ADDR", "localhost:50057"),
		NotificationServiceAddr: getEnv("NOTIFICATION_SERVICE_ADDR", "localhost:50058"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
