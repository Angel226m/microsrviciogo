// ═══════════════════════════════════════════════════════════════
// JWT Authentication Middleware
// Validates Bearer tokens and injects user context
// ═══════════════════════════════════════════════════════════════
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "user_role"
	EmailKey  contextKey = "user_email"
)

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// JWTAuth returns middleware that validates JWT tokens from the Authorization header.
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeError(w, http.StatusUnauthorized, "missing authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
				writeError(w, http.StatusUnauthorized, "invalid authorization format")
				return
			}

			tokenString := parts[1]

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(secret), nil
			})

			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "invalid or expired token")
				return
			}

			// Inject user info into request context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			// Forward user info as headers to downstream services
			r.Header.Set("X-User-ID", claims.UserID)
			r.Header.Set("X-User-Role", claims.Role)
			r.Header.Set("X-User-Email", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"unauthorized","message":"` + msg + `"}`))
}
