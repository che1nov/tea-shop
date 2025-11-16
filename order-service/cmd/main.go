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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/che1nov/tea-shop/order-service/config"
	"github.com/che1nov/tea-shop/order-service/internal/handler"
	"github.com/che1nov/tea-shop/order-service/internal/kafka"
	"github.com/che1nov/tea-shop/order-service/internal/repository"
	"github.com/che1nov/tea-shop/order-service/internal/service"
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
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			items JSONB NOT NULL,
			status VARCHAR(50) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			address TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id);
		CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);

		-- Миграция: добавляем колонку address для существующих заказов, если её нет
		DO $$
		BEGIN
			IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='orders' AND column_name='address') THEN
				ALTER TABLE orders ADD COLUMN address TEXT;
			END IF;
		END $$;
	`
	if _, err := db.Exec(createTablesSQL); err != nil {
		panic(err)
	}

	// Инициализируем Kafka producer
	producer := kafka.NewProducer(cfg.Kafka.Brokers)

	// Подключаемся к другим сервисам через gRPC
	goodsConn, err := grpc.Dial(cfg.Services.GoodsService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	paymentConn, err := grpc.Dial(cfg.Services.PaymentService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	deliveryConn, err := grpc.Dial(cfg.Services.DeliveryService, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	// Инициализируем слои
	repo := repository.New(db)
	goodsClient := pb.NewGoodsServiceClient(goodsConn)
	paymentClient := pb.NewPaymentsServiceClient(paymentConn)
	deliveryClient := pb.NewDeliveryServiceClient(deliveryConn)
	svc := service.New(repo, producer, goodsClient, paymentClient, deliveryClient)
	hdlr := handler.New(svc)

	// Запускаем HTTP сервер для метрик Prometheus ПЕРВЫМ
	metricsPort := 9003
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
		logger.Info("Order Service metrics server starting", "port", metricsPort)
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
		// Закрываем соединения
		if err := producer.Close(); err != nil {
			logger.Error("Error closing Kafka producer", "error", err)
		}
		<-sigChan
		logger.Info("Shutting down Order Service...")
		if err := metricsServer.Close(); err != nil {
			logger.Error("Error closing metrics server", "error", err)
		}
		logger.Info("Order Service stopped")
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterOrdersServiceServer(grpcServer, hdlr)

	// Health check
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("pb.OrdersService", grpc_health_v1.HealthCheckResponse_SERVING)

	// Enable gRPC reflection for easier testing
	reflection.Register(grpcServer)

	// Запускаем gRPC сервер в отдельной горутине
	go func() {
		logger.Info("Order Service gRPC server started", "port", cfg.Server.Port)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	<-sigChan
	logger.Info("Shutting down Order Service...")

	// Graceful shutdown HTTP сервера метрик
	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	// Graceful shutdown gRPC сервера
	grpcServer.GracefulStop()

	// Закрываем соединения
	if err := producer.Close(); err != nil {
		logger.Error("Error closing Kafka producer", "error", err)
	}

	if err := goodsConn.Close(); err != nil {
		logger.Error("Error closing goods service connection", "error", err)
	}

	if err := paymentConn.Close(); err != nil {
		logger.Error("Error closing payment service connection", "error", err)
	}

	if err := db.Close(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Order Service stopped")
}
