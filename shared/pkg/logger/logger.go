package logger

import (
	"context"
	"log/slog"
	"os"
)

var defaultLogger *slog.Logger

// Init инициализирует логер с нужными параметрами
func Init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
		// Используй LevelDebug для разработки
		// Level: slog.LevelDebug,
	}

	// JSON handler для продакшна
	handler := slog.NewJSONHandler(os.Stdout, opts)

	// Или текстовый handler для разработки
	// handler := slog.NewTextHandler(os.Stdout, opts)

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// GetLogger возвращает дефолтный логер
func GetLogger() *slog.Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}

// Удобные функции для логирования
func Info(msg string, args ...any) {
	GetLogger().Info(msg, args...)
}

func Error(msg string, args ...any) {
	GetLogger().Error(msg, args...)
}

func Debug(msg string, args ...any) {
	GetLogger().Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	GetLogger().Warn(msg, args...)
}

// LogRequest логирует HTTP/gRPC запрос
func LogRequest(ctx context.Context, method string, params ...any) {
	GetLogger().InfoContext(ctx, "request", append([]any{"method", method}, params...)...)
}

// LogResponse логирует HTTP/gRPC ответ
func LogResponse(ctx context.Context, method string, duration float64, params ...any) {
	GetLogger().InfoContext(ctx, "response", append([]any{"method", method, "duration_ms", duration}, params...)...)
}

// LogError логирует ошибку с контекстом
func LogError(ctx context.Context, msg string, err error, params ...any) {
	args := append([]any{"error", err}, params...)
	GetLogger().ErrorContext(ctx, msg, args...)
}
