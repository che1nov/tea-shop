package config

import "os"

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
	cfg.Kafka.Brokers = []string{"localhost:9092"}
	cfg.Services.GoodsService = "localhost:8002"
	cfg.Services.PaymentService = "localhost:8004"
	cfg.Services.DeliveryService = "localhost:8005"
	cfg.Services.UserService = "localhost:8001"

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
