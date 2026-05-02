package config

import (
	"os"
	"strings"
)

type Config struct {
	Server struct {
		Port int
	}
	Kafka struct {
		Brokers []string
		Group   string
	}
	Email struct {
		SMTPHost string
		SMTPPort string
		From     string
		Password string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Server.Port = 8006
	cfg.Kafka.Brokers = getEnvList("KAFKA_BROKERS", []string{"localhost:9092"})
	cfg.Kafka.Group = getEnv("KAFKA_GROUP", "notify-service")
	cfg.Email.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")
	cfg.Email.SMTPPort = getEnv("SMTP_PORT", "587")
	cfg.Email.From = getEnv("EMAIL_FROM", "noreply@ecommerce.com")
	cfg.Email.Password = getEnv("EMAIL_PASSWORD", "")

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
