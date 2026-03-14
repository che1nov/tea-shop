package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadDefaults(t *testing.T) {
	cfg := Load()
	require.Equal(t, "5434", cfg.Database.Port)
	require.Equal(t, 8003, cfg.Server.Port)
	require.Equal(t, "localhost:8002", cfg.Services.GoodsService)
	require.Equal(t, "localhost:9092", cfg.Kafka.Brokers[0])
}
