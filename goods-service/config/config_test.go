package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaultsAndEnv(t *testing.T) {
	t.Setenv("DB_NAME", "goods_test")
	cfg := Load()
	require.Equal(t, "goods_test", cfg.Database.Name)
	require.Equal(t, "5433", cfg.Database.Port)
	require.Equal(t, 8002, cfg.Server.Port)
}
