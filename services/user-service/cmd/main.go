// ═══════════════════════════════════════════════════════════════
// User Service – Entry Point
// Wires all hexagonal layers: domain -> application -> infrastructure
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

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/redis/go-redis/v9"

	appservice "github.com/cloudmart/user-service/internal/application/service"
	"github.com/cloudmart/user-service/internal/infrastructure/adapter/cache"
	handler "github.com/cloudmart/user-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/user-service/internal/infrastructure/adapter/messaging"
	"github.com/cloudmart/user-service/internal/infrastructure/adapter/repository"
	"github.com/cloudmart/user-service/internal/infrastructure/config"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// ── Connect PostgreSQL ──────────────────────────────────
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer pool.Close()

	// ── Connect Redis ───────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	defer rdb.Close()

	// ── Connect NATS ────────────────────────────────────────
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to NATS: %v", err))
	}
	defer nc.Close()

	// ── Wire Hexagonal Architecture ─────────────────────────
	// Driven adapters (infrastructure → domain ports)
	userRepo := repository.NewUserPostgresRepo(pool)
	addrRepo := repository.NewAddressPostgresRepo(pool)
	cacheAdapter := cache.NewRedisCache(rdb)
	eventPub := messaging.NewNATSPublisher(nc)

	// Application service (use cases)
	userService := appservice.NewUserService(userRepo, addrRepo, cacheAdapter, eventPub, cfg.JWTSecret)

	// Driving adapter (HTTP → application)
	userHandler := handler.NewUserHandler(userService)

	// ── Build Router ────────────────────────────────────────
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"user-service"}`))
	})

	userHandler.RegisterRoutes(r)

	// ── Start Server ────────────────────────────────────────
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		fmt.Printf("🚀 User Service starting on port %s\n", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	fmt.Println("✅ User Service stopped")
}
