package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tea_shop_http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"service", "method", "route", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "tea_shop_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "route", "status"},
	)

	businessEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tea_shop_business_events_total",
			Help: "Total number of business events.",
		},
		[]string{"service", "event", "status"},
	)

	kafkaMessagesTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tea_shop_kafka_messages_total",
			Help: "Total number of Kafka messages.",
		},
		[]string{"service", "topic", "event", "direction", "status"},
	)
)

func init() {
	prometheus.MustRegister(
		httpRequestsTotal,
		httpRequestDuration,
		businessEventsTotal,
		kafkaMessagesTotal,
	)
}

func HTTPMiddleware(service string, routeName func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			recorder := &statusRecorder{
				ResponseWriter: w,
				statusCode:     http.StatusOK,
			}

			next.ServeHTTP(recorder, r)

			route := r.URL.Path
			if routeName != nil {
				if name := routeName(r); name != "" {
					route = name
				}
			}

			statusLabel := strconv.Itoa(recorder.statusCode)
			httpRequestsTotal.WithLabelValues(service, r.Method, route, statusLabel).Inc()
			httpRequestDuration.WithLabelValues(service, r.Method, route, statusLabel).Observe(time.Since(start).Seconds())
		})
	}
}

func ObserveBusinessEvent(service, event, status string) {
	businessEventsTotal.WithLabelValues(service, event, status).Inc()
}

func ObserveKafkaMessage(service, topic, event, direction, status string) {
	kafkaMessagesTotal.WithLabelValues(service, topic, event, direction, status).Inc()
}

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (r *statusRecorder) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}
