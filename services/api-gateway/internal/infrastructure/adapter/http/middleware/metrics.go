// ═══════════════════════════════════════════════════════════════
// Prometheus Metrics Middleware
// Records HTTP request duration, count, and response size
// ═══════════════════════════════════════════════════════════════
package middleware

import (
	"net/http"
	"strconv"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "cloudmart_gateway_http_requests_total",
			Help: "Total number of HTTP requests processed by the API Gateway",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudmart_gateway_http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5},
		},
		[]string{"method", "path"},
	)

	httpResponseSize = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "cloudmart_gateway_http_response_size_bytes",
			Help:    "Size of HTTP responses in bytes",
			Buckets: prometheus.ExponentialBuckets(100, 10, 7),
		},
		[]string{"method", "path"},
	)

	activeConnections = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "cloudmart_gateway_active_connections",
			Help: "Number of active connections to the API Gateway",
		},
	)
)

// MetricsMiddleware records Prometheus metrics for every request.
func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			activeConnections.Inc()
			defer activeConnections.Dec()

			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()
			status := strconv.Itoa(ww.Status())
			path := r.URL.Path

			httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
			httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
			httpResponseSize.WithLabelValues(r.Method, path).Observe(float64(ww.BytesWritten()))
		})
	}
}
