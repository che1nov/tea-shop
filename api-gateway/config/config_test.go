package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadWithJWTEnv(t *testing.T) {
	t.Setenv("JWT_SECRET", "gateway-secret")
	cfg := Load()
	require.Equal(t, 8080, cfg.Server.Port)
	require.Equal(t, "gateway-secret", cfg.JWT.Secret)
	require.Equal(t, "localhost:8001", cfg.Services.UsersService)
}
