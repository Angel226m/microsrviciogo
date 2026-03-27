// ═══════════════════════════════════════════════════════════════
// HTTP Router – API Gateway routing with Chi
// Reverse proxy + middleware chain (auth, rate limit, circuit breaker)
// ═══════════════════════════════════════════════════════════════
package http

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/cloudmart/api-gateway/internal/infrastructure/config"
	"github.com/cloudmart/api-gateway/internal/infrastructure/adapter/http/middleware"
	"github.com/cloudmart/api-gateway/pkg/logger"
)

// NewRouter creates the main API Gateway router with all middleware and proxy routes.
func NewRouter(cfg *config.Config, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// ── Global Middleware ────────────────────────────────────
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.RequestLogger(log))
	r.Use(middleware.MetricsMiddleware())
	r.Use(chimw.Recoverer)
	r.Use(chimw.Compress(5))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://cloudmart.dev"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID", "X-Trace-ID"},
		ExposedHeaders:   []string{"X-Request-ID", "X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ── Rate Limiting ───────────────────────────────────────
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// ── Health & Metrics ────────────────────────────────────
	r.Get("/health", healthCheck)
	r.Get("/ready", readinessCheck)
	r.Handle("/metrics", promhttp.Handler())

	// ── API v1 Routes ───────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes (no auth required)
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", proxyTo(cfg.UserServiceURL, "/api/v1/auth/login"))
			r.Post("/auth/register", proxyTo(cfg.UserServiceURL, "/api/v1/auth/register"))
			r.Get("/products", proxyTo(cfg.ProductServiceURL, "/api/v1/products"))
			r.Get("/products/{slug}", proxyTo(cfg.ProductServiceURL, "/api/v1/products/{slug}"))
			r.Get("/categories", proxyTo(cfg.ProductServiceURL, "/api/v1/categories"))
			r.Get("/products/{id}/reviews", proxyTo(cfg.ProductServiceURL, "/api/v1/products/{id}/reviews"))
		})

		// Protected routes (auth required)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWTSecret))

			// User
			r.Get("/users/me", proxyTo(cfg.UserServiceURL, "/api/v1/users/me"))
			r.Put("/users/me", proxyTo(cfg.UserServiceURL, "/api/v1/users/me"))
			r.Get("/users/me/addresses", proxyTo(cfg.UserServiceURL, "/api/v1/users/me/addresses"))
			r.Post("/users/me/addresses", proxyTo(cfg.UserServiceURL, "/api/v1/users/me/addresses"))

			// Products (write)
			r.Post("/products", proxyTo(cfg.ProductServiceURL, "/api/v1/products"))
			r.Put("/products/{id}", proxyTo(cfg.ProductServiceURL, "/api/v1/products/{id}"))
			r.Delete("/products/{id}", proxyTo(cfg.ProductServiceURL, "/api/v1/products/{id}"))
			r.Post("/products/{id}/reviews", proxyTo(cfg.ProductServiceURL, "/api/v1/products/{id}/reviews"))

			// Orders
			r.Post("/orders", proxyTo(cfg.OrderServiceURL, "/api/v1/orders"))
			r.Get("/orders", proxyTo(cfg.OrderServiceURL, "/api/v1/orders"))
			r.Get("/orders/{id}", proxyTo(cfg.OrderServiceURL, "/api/v1/orders/{id}"))
			r.Put("/orders/{id}/cancel", proxyTo(cfg.OrderServiceURL, "/api/v1/orders/{id}/cancel"))

			// Payments
			r.Post("/payments", proxyTo(cfg.PaymentServiceURL, "/api/v1/payments"))
			r.Get("/payments/{id}", proxyTo(cfg.PaymentServiceURL, "/api/v1/payments/{id}"))

			// Inventory (admin)
			r.Get("/inventory", proxyTo(cfg.InventoryServiceURL, "/api/v1/inventory"))
			r.Get("/inventory/{productId}", proxyTo(cfg.InventoryServiceURL, "/api/v1/inventory/{productId}"))
			r.Put("/inventory/{productId}", proxyTo(cfg.InventoryServiceURL, "/api/v1/inventory/{productId}"))
		})
	})

	return r
}

// ── Reverse Proxy Helpers ──────────────────────────────────────────

func proxyTo(targetBase, path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		target, err := url.Parse(targetBase)
		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, `{"error":"service_unavailable","message":"upstream service is not responding"}`, http.StatusServiceUnavailable)
		}

		r.URL.Path = path
		r.URL.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = target.Host

		proxy.ServeHTTP(w, r)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"api-gateway","version":"1.0.0"}`))
}

func readinessCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}
