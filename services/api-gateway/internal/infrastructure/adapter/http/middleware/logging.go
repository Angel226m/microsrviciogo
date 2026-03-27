// ═══════════════════════════════════════════════════════════════
// Request Logger Middleware – Structured logging for all requests
// ═══════════════════════════════════════════════════════════════
package middleware

import (
	"net/http"
	"time"

	"github.com/cloudmart/api-gateway/pkg/logger"
	chimw "github.com/go-chi/chi/v5/middleware"
)

// RequestLogger returns middleware that logs each HTTP request with structured fields.
func RequestLogger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			defer func() {
				log.Info("request",
					"method", r.Method,
					"path", r.URL.Path,
					"status", ww.Status(),
					"bytes", ww.BytesWritten(),
					"duration_ms", time.Since(start).Milliseconds(),
					"remote", r.RemoteAddr,
					"request_id", chimw.GetReqID(r.Context()),
					"user_agent", r.UserAgent(),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
