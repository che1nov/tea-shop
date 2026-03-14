package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	require.Equal(t, "5435", cfg.Database.Port)
	require.Equal(t, 8004, cfg.Server.Port)
}
