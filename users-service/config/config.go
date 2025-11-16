package config

import (
	"os"
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
	JWT struct {
		Secret string
	}
	Admin struct {
		Email    string
		Password string
	}
}

func Load() *Config {
	cfg := &Config{}

	cfg.Database.Host = getEnv("DB_HOST", "localhost")
	cfg.Database.Port = getEnv("DB_PORT", "5432")
	cfg.Database.User = getEnv("DB_USER", "user")
	cfg.Database.Password = getEnv("DB_PASSWORD", "password")
	cfg.Database.Name = getEnv("DB_NAME", "users_db")
	cfg.Server.Port = 8001
	cfg.JWT.Secret = getEnv("JWT_SECRET", "sfndknjs4234njndfgk")
	cfg.Admin.Email = getEnv("ADMIN_EMAIL", "admin@example.com")
	cfg.Admin.Password = getEnv("ADMIN_PASSWORD", "admin123")

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
