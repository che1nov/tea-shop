package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	require.Equal(t, "5436", cfg.Database.Port)
	require.Equal(t, 8005, cfg.Server.Port)
	require.Equal(t, "localhost:8004", cfg.Services.PaymentService)
}
