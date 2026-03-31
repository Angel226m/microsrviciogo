// ═══════════════════════════════════════════════════════════════
// Middleware de Seguridad – Cabeceras OWASP Top 10 e ISO 27001
// Protege contra XSS, Clickjacking, MIME sniffing y más
// ═══════════════════════════════════════════════════════════════
package middleware

import "net/http"

// CabecerasSeguridad agrega cabeceras de seguridad HTTP según OWASP.
// Cumple con los controles A.14.1.2 y A.14.1.3 de ISO 27001.
func CabecerasSeguridad() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Prevenir clickjacking (OWASP A5)
			w.Header().Set("X-Frame-Options", "DENY")

			// Prevenir MIME type sniffing
			w.Header().Set("X-Content-Type-Options", "nosniff")

			// Protección XSS del navegador
			w.Header().Set("X-XSS-Protection", "1; mode=block")

			// Política de referencia estricta
			w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")

			// Política de permisos restrictiva
			w.Header().Set("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=(self)")

			// Forzar HTTPS en producción
			w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")

			// Política de seguridad de contenido
			w.Header().Set("Content-Security-Policy", "default-src 'self'; frame-ancestors 'none'")

			// Prevenir cache de datos sensibles
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			w.Header().Set("Pragma", "no-cache")

			// Ocultar tecnología del servidor
			w.Header().Del("X-Powered-By")
			w.Header().Del("Server")

			next.ServeHTTP(w, r)
		})
	}
}

// ValidarContentType verifica que las peticiones POST/PUT tengan Content-Type válido.
// Previene ataques de inyección (OWASP A1/A3).
func ValidarContentType() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
				ct := r.Header.Get("Content-Type")
				if ct == "" || (ct != "application/json" && ct != "application/json; charset=utf-8") {
					writeError(w, http.StatusUnsupportedMediaType, "Content-Type debe ser application/json")
					return
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// LimitarTamanoBody restringe el tamaño máximo del cuerpo de la petición.
// Previene ataques de DoS por payloads grandes (OWASP A6).
func LimitarTamanoBody(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body != nil {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
