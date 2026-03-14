package logger

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLoggerAndHelpers(t *testing.T) {
	Init()
	require.NotNil(t, GetLogger())

	Info("info")
	Warn("warn")
	Debug("debug")
	Error("error")
	LogRequest(context.Background(), "GET /health", "x", 1)
	LogResponse(context.Background(), "GET /health", 1.2, "status", 200)
	LogError(context.Background(), "boom", context.Canceled)
}
