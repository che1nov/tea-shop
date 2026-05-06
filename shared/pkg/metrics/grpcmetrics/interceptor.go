package grpcmetrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	grpcRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tea_shop_grpc_requests_total",
			Help: "Total number of gRPC requests.",
		},
		[]string{"service", "method", "code"},
	)

	grpcRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tea_shop_grpc_request_duration_seconds",
			Help:    "gRPC request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "code"},
	)
)

func init() {
	prometheus.MustRegister(grpcRequestsTotal, grpcRequestDuration)
}

func UnaryServerInterceptor(service string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		start := time.Now()

		resp, err := handler(ctx, req)
		code := codes.OK
		if err != nil {
			code = status.Code(err)
		}

		codeLabel := code.String()
		grpcRequestsTotal.WithLabelValues(service, info.FullMethod, codeLabel).Inc()
		grpcRequestDuration.WithLabelValues(service, info.FullMethod, codeLabel).Observe(time.Since(start).Seconds())

		return resp, err
	}
}
