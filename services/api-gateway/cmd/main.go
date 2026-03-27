// ═══════════════════════════════════════════════════════════════
// CloudMart – API Gateway · Entry Point
// Reverse proxy with auth, rate limiting, circuit breaker & tracing
// ═══════════════════════════════════════════════════════════════
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/cloudmart/api-gateway/internal/infrastructure/config"
	gateway "github.com/cloudmart/api-gateway/internal/infrastructure/adapter/http"
	"github.com/cloudmart/api-gateway/pkg/logger"
	"github.com/cloudmart/api-gateway/pkg/telemetry"
)

func main() {
	// ── Load configuration ──────────────────────────────────
	cfg := config.Load()

	// ── Initialize logger ───────────────────────────────────
	log := logger.New(cfg.LogLevel)
	defer log.Sync()

	// ── Initialize OpenTelemetry ────────────────────────────
	tp, err := telemetry.InitTracer("api-gateway", cfg.JaegerEndpoint)
	if err != nil {
		log.Fatal("failed to initialize tracer", "error", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tp.Shutdown(ctx)
	}()

	// ── Build router ────────────────────────────────────────
	router := gateway.NewRouter(cfg, log)

	// ── Create HTTP server ──────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// ── Graceful shutdown ───────────────────────────────────
	go func() {
		log.Info("🚀 API Gateway starting", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("⏳ Shutting down API Gateway...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("forced shutdown", "error", err)
	}

	log.Info("✅ API Gateway stopped gracefully")
}
