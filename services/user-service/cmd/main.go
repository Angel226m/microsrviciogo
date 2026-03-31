// ═══════════════════════════════════════════════════════════════
// Servicio de Usuarios – Punto de Entrada
// Conecta todas las capas hexagonales: dominio → aplicación → infraestructura
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

	// ── Conectar PostgreSQL ──────────────────────────────────
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}
	defer pool.Close()

	// ── Conectar Redis ─────────────────────────────────────
	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisURL})
	defer rdb.Close()

	// ── Conectar NATS ──────────────────────────────────────
	nc, err := nats.Connect(cfg.NatsURL)
	if err != nil {
		panic(fmt.Sprintf("failed to connect to NATS: %v", err))
	}
	defer nc.Close()

	// ── Cableado de Arquitectura Hexagonal ─────────────────────
	// Adaptadores secundarios (infraestructura → puertos de dominio)
	userRepo := repository.NewUserPostgresRepo(pool)
	addrRepo := repository.NewAddressPostgresRepo(pool)
	cacheAdapter := cache.NewRedisCache(rdb)
	eventPub := messaging.NewNATSPublisher(nc)

	// Servicio de aplicación (casos de uso)
	userService := appservice.NewUserService(userRepo, addrRepo, cacheAdapter, eventPub, cfg.JWTSecret)

	// Adaptador primario (HTTP → aplicación)
	userHandler := handler.NewUserHandler(userService)

	// ── Construir Enrutador ──────────────────────────────────
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

	// ── Iniciar Servidor ────────────────────────────────────
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
