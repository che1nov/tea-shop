package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadWithEnv(t *testing.T) {
	t.Setenv("SMTP_HOST", "smtp.test")
	t.Setenv("EMAIL_FROM", "robot@test.local")

	cfg := Load()
	require.Equal(t, "smtp.test", cfg.Email.SMTPHost)
	require.Equal(t, "robot@test.local", cfg.Email.From)
	require.Equal(t, "notify-service", cfg.Kafka.Group)
	require.Equal(t, 8006, cfg.Server.Port)
}
