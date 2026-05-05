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
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/che1nov/tea-shop/api-gateway/docs"
)

// @title           E-commerce Tea Shop API
// @version         1.0
// @description     Микросервисная платформа для интернет-магазина чая. API Gateway для управления товарами, заказами, платежами и доставкой.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @contact.email   support@example.com
// @host            localhost:8080
// @BasePath        /api/v1
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
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

	router := chi.NewRouter()
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)
	router.Use(middleware.CORSMiddleware(cfg.CORS.AllowedOrigins))

	router.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("/swagger/doc.json")))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	router.Post("/api/v1/auth/register", h.RegisterUser)
	router.Post("/api/v1/auth/login", h.Login)

	router.Get("/api/v1/goods", h.ListGoods)
	router.Get("/api/v1/goods/{id}", h.GetGood)

	router.Route("/api/v1/admin", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		r.Use(middleware.AdminMiddleware())

		r.Post("/goods", h.CreateGood)
		r.Put("/goods/{id}", h.UpdateGood)
		r.Delete("/goods/{id}", h.DeleteGood)

		r.Get("/deliveries", h.ListDeliveries)
		r.Put("/deliveries/{id}/status", h.UpdateDeliveryStatus)
	})

	router.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(cfg.JWT.Secret))

		r.Get("/users/me", h.GetUser)

		r.Post("/orders", h.CreateOrder)
		r.Get("/orders/{id}", h.GetOrder)

		r.Get("/payments/{id}", h.GetPayment)

		r.Post("/deliveries", h.CreateDelivery)
		r.Get("/deliveries/{id}", h.GetDelivery)
	})

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
