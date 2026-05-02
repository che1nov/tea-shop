package config

import (
	"os"
	"strings"
)

type Config struct {
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
	}
	Server struct {
		Port int
	}
	Kafka struct {
		Brokers []string
	}
	Services struct {
		GoodsService    string
		PaymentService  string
		DeliveryService string
		UserService     string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5434")
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "orders_db")
	cfg.Server.Port = 8003
	cfg.Kafka.Brokers = getEnvList("KAFKA_BROKERS", []string{"localhost:9092"})
	cfg.Services.GoodsService = getEnv("GOODS_SERVICE", "localhost:8002")
	cfg.Services.PaymentService = getEnv("PAYMENT_SERVICE", "localhost:8004")
	cfg.Services.DeliveryService = getEnv("DELIVERY_SERVICE", "localhost:8005")
	cfg.Services.UserService = getEnv("USER_SERVICE", "localhost:8001")

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
