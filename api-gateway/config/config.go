package config

import (
	"os"
	"strings"
)

type Config struct {
	Server struct {
		Port int
	}
	Services struct {
		UsersService    string
		GoodsService    string
		OrdersService   string
		PaymentsService string
		DeliveryService string
	}
	JWT struct {
		Secret string
	}
	CORS struct {
		AllowedOrigins []string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Server.Port = 8080
	cfg.Services.UsersService = getEnv("USERS_SERVICE", "localhost:8001")
	cfg.Services.GoodsService = getEnv("GOODS_SERVICE", "localhost:8002")
	cfg.Services.OrdersService = getEnv("ORDERS_SERVICE", "localhost:8003")
	cfg.Services.PaymentsService = getEnv("PAYMENTS_SERVICE", "localhost:8004")
	cfg.Services.DeliveryService = getEnv("DELIVERY_SERVICE", "localhost:8005")
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")
	cfg.CORS.AllowedOrigins = getEnvList("CORS_ALLOWED_ORIGINS", []string{"http://localhost:5173"})

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvList(key string, defaultValue []string) []string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	if len(result) == 0 {
		return defaultValue
	}
	return result
}
