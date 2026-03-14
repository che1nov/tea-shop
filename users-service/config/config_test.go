package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadWithEnv(t *testing.T) {
	t.Setenv("DB_HOST", "db")
	t.Setenv("JWT_SECRET", "secret")
	t.Setenv("ADMIN_EMAIL", "admin@test.local")

	cfg := Load()
	require.Equal(t, "db", cfg.Database.Host)
	require.Equal(t, "secret", cfg.JWT.Secret)
	require.Equal(t, "admin@test.local", cfg.Admin.Email)
	require.Equal(t, 8001, cfg.Server.Port)
}
