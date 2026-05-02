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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/che1nov/tea-shop/api-gateway/docs"
)

func main() {
	logger.Init()

	cfg := config.Load()

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

	router := gin.Default()

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	router.POST("/api/v1/auth/register", h.RegisterUser)
	router.POST("/api/v1/auth/login", h.Login)

	router.GET("/api/v1/goods", h.ListGoods)
	router.GET("/api/v1/goods/:id", h.GetGood)

	admin := router.Group("/api/v1/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/goods", h.CreateGood)
		admin.PUT("/goods/:id", h.UpdateGood)
		admin.DELETE("/goods/:id", h.DeleteGood)

		admin.GET("/deliveries", h.ListDeliveries)
		admin.PUT("/deliveries/:id/status", h.UpdateDeliveryStatus)
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
	{
		protected.GET("/users/me", h.GetUser)

		protected.POST("/orders", h.CreateOrder)
		protected.GET("/orders/:id", h.GetOrder)

		protected.GET("/payments/:id", h.GetPayment)

		protected.POST("/deliveries", h.CreateDelivery)
		protected.GET("/deliveries/:id", h.GetDelivery)
	}

	metricsPort := 9007
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", promhttp.Handler())
	metricsServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: metricsMux,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("API Gateway metrics server started", "port", metricsPort)
		if err := metricsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Metrics server error", "error", err)
		}
	}()

	go func() {
		logger.Info("API Gateway started", "port", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("HTTP server error", "error", err)
		}
	}()

	<-sigChan
	logger.Info("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := metricsServer.Close(); err != nil {
		logger.Error("Error closing metrics server", "error", err)
	}

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error shutting down server", "error", err)
	}

	if err := h.Close(); err != nil {
		logger.Error("Error closing gRPC connections", "error", err)
	}

	logger.Info("API Gateway stopped")
}
