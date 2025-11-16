package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/che1nov/tea-shop/goods-service/config"
	"github.com/che1nov/tea-shop/goods-service/internal/handler"
	"github.com/che1nov/tea-shop/goods-service/internal/repository"
	"github.com/che1nov/tea-shop/goods-service/internal/service"
	pb "github.com/che1nov/tea-shop/shared/pb"
	"github.com/che1nov/tea-shop/shared/pkg/logger"
)

func main() {
	// Инициализируем logger
	logger.Init()

	cfg := config.Load()

	// Подключение к БД
	dbConnStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		panic(err)
	}

	// Проверяем подключение к БД
	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	logger.Info("Database connection established")

	// Создаём таблицы
	createTablesSQL := `
		CREATE TABLE IF NOT EXISTS goods (
			id SERIAL PRIMARY KEY,
			sku VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE TABLE IF NOT EXISTS stock_reservations (
			id SERIAL PRIMARY KEY,
			good_id INT NOT NULL REFERENCES goods(id),
			order_id INT NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_goods_name ON goods(name);
		CREATE INDEX IF NOT EXISTS idx_goods_sku ON goods(sku);
		CREATE INDEX IF NOT EXISTS idx_reservations_order ON stock_reservations(order_id);

		-- Миграция: добавляем колонку sku для существующих товаров, если её нет
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='goods' AND column_name='sku') THEN
				ALTER TABLE goods ADD COLUMN sku VARCHAR(50);
				UPDATE goods SET sku = 'GOOD-' || LPAD(id::text, 6, '0') WHERE sku IS NULL;
				ALTER TABLE goods ALTER COLUMN sku SET NOT NULL;
				ALTER TABLE goods ADD CONSTRAINT goods_sku_unique UNIQUE (sku);
				CREATE INDEX IF NOT EXISTS idx_goods_sku ON goods(sku);
			END IF;
		END $$;
	`
	if _, err := db.Exec(createTablesSQL); err != nil {
		panic(err)
	}

	// Инициализируем слои
	repo := repository.New(db)
	svc := service.New(repo)
	hdlr := handler.New(svc)

	// Запускаем HTTP сервер для метрик Prometheus ПЕРВЫМ
	metricsPort := 9002
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: metricsMux,
	}

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запускаем HTTP сервер для метрик в отдельной горутине (ПЕРВЫМ)
	go func() {
		logger.Info("Goods Service metrics server starting", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", "error", err)
		}
	}()

	// Небольшая задержка для запуска HTTP сервера метрик
	time.Sleep(100 * time.Millisecond)

	// Запускаем gRPC сервер
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		logger.Error("Failed to create gRPC listener", "error", err, "port", cfg.Server.Port)
		logger.Warn("gRPC server will not start, but metrics server is running")
		<-sigChan
		logger.Info("Shutting down Goods Service...")
		if err := metricsServer.Close(); err != nil {
			logger.Error("Error closing metrics server", "error", err)
		}
		logger.Info("Goods Service stopped")
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoodsServiceServer(grpcServer, hdlr)

	// Health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("pb.GoodsService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable gRPC reflection for easier testing
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		logger.Info("Goods Service gRPC server started", "port", cfg.Server.Port)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	<-sigChan
	logger.Info("Shutting down Goods Service...")

	// Graceful shutdown HTTP сервера метрик
	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	// Graceful shutdown gRPC сервера
	grpcServer.GracefulStop()

	// Закрываем соединение с БД
	if err := db.Close(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Goods Service stopped")
}
