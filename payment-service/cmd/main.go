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

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pb "github.com/che1nov/tea-shop/shared/pb"
	"github.com/che1nov/tea-shop/shared/pkg/logger"

	"github.com/che1nov/tea-shop/payment-service/config"
	"github.com/che1nov/tea-shop/payment-service/internal/handler"
	"github.com/che1nov/tea-shop/payment-service/internal/repository"
	"github.com/che1nov/tea-shop/payment-service/internal/service"
)

func main() {
	logger.Init()

	cfg := config.Load()

	dbConnStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.Host,
		cfg.Database.Port,
	)

	db, err := sql.Open("pgx", dbConnStr)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		panic(err)
	}

	if err := db.Ping(); err != nil {
		logger.Error("Failed to connect to database", "error", err)
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	logger.Info("Database connection established")

	createTablesSQL := `
		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL UNIQUE,
			amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			method VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);

		CREATE INDEX IF NOT EXISTS idx_payments_order ON payments(order_id);
	`
	if _, err := db.Exec(createTablesSQL); err != nil {
		panic(err)
	}

	repo := repository.New(db)
	svc := service.New(repo)
	hdlr := handler.New(svc)

	metricsPort := 9004
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: metricsMux,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Payment Service metrics server starting", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", "error", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Server.Port))
	if err != nil {
		logger.Error("Failed to create gRPC listener", "error", err, "port", cfg.Server.Port)
		logger.Warn("gRPC server will not start, but metrics server is running")
		<-sigChan
		logger.Info("Shutting down Payment Service...")
		if err := metricsServer.Close(); err != nil {
			logger.Error("Error closing metrics server", "error", err)
		}
		if err := db.Close(); err != nil {
			logger.Error("Error closing database", "error", err)
		}
		logger.Info("Payment Service stopped")
		return
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentsServiceServer(grpcServer, hdlr)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("pb.PaymentsService", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(grpcServer)

	go func() {
		logger.Info("Payment Service gRPC server started", "port", cfg.Server.Port)
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server error", "error", err)
		}
	}()

	<-sigChan
	logger.Info("Shutting down Payment Service...")

	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	grpcServer.GracefulStop()

	if err := db.Close(); err != nil {
		logger.Error("Error closing database", "error", err)
	}

	logger.Info("Payment Service stopped")
}
