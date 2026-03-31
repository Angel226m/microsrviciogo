// ═══════════════════════════════════════════════════════════════
// Enrutador HTTP – API Gateway con Chi
// Proxy inverso + cadena de middlewares (auth, rate limit, seguridad OWASP)
// Cumple con controles de seguridad ISO 27001 Anexo A.14
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

	"github.com/cloudmart/api-gateway/internal/infrastructure/adapter/http/middleware"
	"github.com/cloudmart/api-gateway/internal/infrastructure/config"
	"github.com/cloudmart/api-gateway/pkg/logger"
)

// NuevoEnrutador crea el enrutador principal del API Gateway con todos los middlewares y rutas proxy.
func NewRouter(cfg *config.Config, log *logger.Logger) http.Handler {
	r := chi.NewRouter()

	// ── Middlewares globales ─────────────────────────────────
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.CabecerasSeguridad())       // Cabeceras OWASP
	r.Use(middleware.ValidarContentType())       // Validación Content-Type
	r.Use(middleware.LimitarTamanoBody(1048576)) // Límite 1MB por petición
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

	// ── Limitación de tasa (OWASP – protección contra fuerza bruta) ──
	r.Use(httprate.LimitByIP(100, 1*time.Minute))

	// ── Salud y métricas ────────────────────────────────────
	r.Get("/health", verificarSalud)
	r.Get("/ready", verificarDisponibilidad)
	r.Handle("/metrics", promhttp.Handler())

	// ── Rutas API v1 ────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		// Rutas públicas (sin autenticación requerida)
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", proxyHacia(cfg.UserServiceURL, "/api/v1/auth/login"))
			r.Post("/auth/register", proxyHacia(cfg.UserServiceURL, "/api/v1/auth/register"))
			r.Get("/products", proxyHacia(cfg.ProductServiceURL, "/api/v1/products"))
			r.Get("/products/{slug}", proxyHacia(cfg.ProductServiceURL, "/api/v1/products/{slug}"))
			r.Get("/categories", proxyHacia(cfg.ProductServiceURL, "/api/v1/categories"))
			r.Get("/products/{id}/reviews", proxyHacia(cfg.ProductServiceURL, "/api/v1/products/{id}/reviews"))
		})

		// Rutas protegidas (autenticación JWT requerida)
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTAuth(cfg.JWTSecret))

			// Usuarios
			r.Get("/users/me", proxyHacia(cfg.UserServiceURL, "/api/v1/users/me"))
			r.Put("/users/me", proxyHacia(cfg.UserServiceURL, "/api/v1/users/me"))
			r.Get("/users/me/addresses", proxyHacia(cfg.UserServiceURL, "/api/v1/users/me/addresses"))
			r.Post("/users/me/addresses", proxyHacia(cfg.UserServiceURL, "/api/v1/users/me/addresses"))

			// Productos (escritura)
			r.Post("/products", proxyHacia(cfg.ProductServiceURL, "/api/v1/products"))
			r.Put("/products/{id}", proxyHacia(cfg.ProductServiceURL, "/api/v1/products/{id}"))
			r.Delete("/products/{id}", proxyHacia(cfg.ProductServiceURL, "/api/v1/products/{id}"))
			r.Post("/products/{id}/reviews", proxyHacia(cfg.ProductServiceURL, "/api/v1/products/{id}/reviews"))

			// Pedidos
			r.Post("/orders", proxyHacia(cfg.OrderServiceURL, "/api/v1/orders"))
			r.Get("/orders", proxyHacia(cfg.OrderServiceURL, "/api/v1/orders"))
			r.Get("/orders/{id}", proxyHacia(cfg.OrderServiceURL, "/api/v1/orders/{id}"))
			r.Put("/orders/{id}/cancel", proxyHacia(cfg.OrderServiceURL, "/api/v1/orders/{id}/cancel"))

			// Pagos
			r.Post("/payments", proxyHacia(cfg.PaymentServiceURL, "/api/v1/payments"))
			r.Get("/payments/{id}", proxyHacia(cfg.PaymentServiceURL, "/api/v1/payments/{id}"))

			// Inventario (administración)
			r.Get("/inventory", proxyHacia(cfg.InventoryServiceURL, "/api/v1/inventory"))
			r.Get("/inventory/{productId}", proxyHacia(cfg.InventoryServiceURL, "/api/v1/inventory/{productId}"))
			r.Put("/inventory/{productId}", proxyHacia(cfg.InventoryServiceURL, "/api/v1/inventory/{productId}"))
		})
	})

	return r
}

// ── Helpers de Proxy Inverso ───────────────────────────────────────

// proxyHacia crea un handler que redirige la petición al servicio destino.
func proxyHacia(urlBase, ruta string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		destino, err := url.Parse(urlBase)
		if err != nil {
			http.Error(w, "Bad gateway", http.StatusBadGateway)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(destino)
		proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, `{"error":"servicio_no_disponible","mensaje":"el servicio upstream no responde"}`, http.StatusServiceUnavailable)
		}

		r.URL.Path = ruta
		r.URL.Host = destino.Host
		r.URL.Scheme = destino.Scheme
		r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))
		r.Host = destino.Host

		proxy.ServeHTTP(w, r)
	}
}

// verificarSalud devuelve el estado de salud del servicio.
func verificarSalud(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"estado":"saludable","servicio":"api-gateway","version":"1.0.0"}`))
}

// verificarDisponibilidad devuelve si el servicio está listo para recibir tráfico.
func verificarDisponibilidad(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"estado":"listo"}`))
}
