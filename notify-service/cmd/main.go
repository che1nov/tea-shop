package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/che1nov/tea-shop/notify-service/config"
	"github.com/che1nov/tea-shop/notify-service/internal/kafka"
	"github.com/che1nov/tea-shop/notify-service/internal/service"
	"github.com/che1nov/tea-shop/shared/pkg/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Инициализируем logger
	logger.Init()

	cfg := config.Load()

	// Инициализируем Kafka consumer
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Group)

	// Инициализируем сервис
	svc := service.New(cfg.Email.From)

	// Создаем контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем HTTP сервер для метрик Prometheus
	metricsPort := 9006
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: metricsMux,
	}

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запускаем HTTP сервер для метрик в отдельной горутине
	go func() {
		logger.Info("Notify Service metrics server started", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", "error", err)
		}
	}()

	// Запускаем consumer в отдельной горутине
	errChan := make(chan error, 1)
	go func() {
		logger.Info("Notify Service started, listening for events...")
		if err := consumer.Start(ctx, svc.HandleEvent); err != nil {
			errChan <- err
		}
	}()

	// Ожидаем сигнал или ошибку
	select {
	case <-sigChan:
		logger.Info("Shutting down Notify Service...")
		cancel() // Отменяем контекст, чтобы consumer остановился
	case err := <-errChan:
		logger.Error("Consumer error", "error", err)
		cancel()
	}

	// Graceful shutdown HTTP сервера метрик
	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	// Закрываем consumer
	if err := consumer.Close(); err != nil {
		logger.Error("Error closing Kafka consumer", "error", err)
	}

	logger.Info("Notify Service stopped")
}
