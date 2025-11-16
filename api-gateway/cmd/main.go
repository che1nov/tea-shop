// @title           E-commerce Tea Shop API
// @version         1.0
// @description     Микросервисная платформа для интернет-магазина чая. API Gateway для управления товарами, заказами, платежами и доставкой.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description JWT токен. Формат: Bearer {token}

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/che1nov/tea-shop/api-gateway/config"
	"github.com/che1nov/tea-shop/api-gateway/internal/handler"
	"github.com/che1nov/tea-shop/api-gateway/internal/middleware"
	"github.com/che1nov/tea-shop/shared/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	
	_ "github.com/che1nov/tea-shop/api-gateway/docs"
)

func main() {
	// Инициализируем logger
	logger.Init()

	cfg := config.Load()

	// Инициализируем handler
	h, err := handler.New(
		cfg.Services.UsersService,
		cfg.Services.GoodsService,
		cfg.Services.OrdersService,
		cfg.Services.PaymentsService,
		cfg.Services.DeliveryService,
	)
	if err != nil {
		logger.Error("Failed to initialize handler", "error", err)
		panic(err)
	}

	// Инициализируем Gin
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public endpoints
	router.POST("/api/v1/auth/register", h.RegisterUser)
	router.POST("/api/v1/auth/login", h.Login)

	// Goods endpoints (публичные - доступны всем)
	router.GET("/api/v1/goods", h.ListGoods)
	router.GET("/api/v1/goods/:id", h.GetGood)
	
	// Admin endpoints (требуют аутентификацию и роль администратора)
	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/goods", h.CreateGood)
		admin.PUT("/goods/:id", h.UpdateGood)
		admin.DELETE("/goods/:id", h.DeleteGood)
		
		// Deliveries endpoints (только для админа)
		admin.GET("/deliveries", h.ListDeliveries)
		admin.PUT("/deliveries/:id/status", h.UpdateDeliveryStatus)
	}

	// Protected endpoints (требуют аутентификацию)
	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		// User endpoints
		protected.GET("/users/me", h.GetUser)

		// Orders endpoints
		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)

		// Payments endpoints
		protected.GET("/payments/:id", h.GetPayment)

		// Delivery endpoints
		protected.POST("/deliveries", h.CreateDelivery)
		protected.GET("/deliveries/:id", h.GetDelivery)
	}

	// Запускаем HTTP сервер для метрик Prometheus
	metricsPort := 9007
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: metricsMux,
	}

	// Создаем HTTP сервер для API
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// Канал для получения сигналов ОС
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Запускаем HTTP сервер для метрик в отдельной горутине
	go func() {
		logger.Info("API Gateway metrics server started", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", "error", err)
		}
	}()

	// Запускаем API сервер в отдельной горутине
	go func() {
		logger.Info("API Gateway started", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	// Ожидаем сигнал для graceful shutdown
	<-sigChan
	logger.Info("Shutting down API Gateway...")

	// Создаем контекст с таймаутом для graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown HTTP сервера метрик
	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	// Graceful shutdown HTTP сервера
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down server", "error", err)
	}

	// Закрываем gRPC соединения
	if err := h.Close(); err != nil {
		logger.Error("Error closing gRPC connections", "error", err)
	}

	logger.Info("API Gateway stopped")
}
