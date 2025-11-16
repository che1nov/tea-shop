package config

import "os"

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
}

func Load() *Config {
	cfg := &Config{}

	cfg.Server.Port = 8080
	cfg.Services.UsersService = "localhost:8001"
	cfg.Services.GoodsService = "localhost:8002"
	cfg.Services.OrdersService = "localhost:8003"
	cfg.Services.PaymentsService = "localhost:8004"
	cfg.Services.DeliveryService = "localhost:8005"
	cfg.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
