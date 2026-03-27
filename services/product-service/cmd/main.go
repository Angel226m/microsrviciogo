package main

import (
	"context"
	"encoding/json"
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

	appservice "github.com/cloudmart/product-service/internal/application/service"
	handler "github.com/cloudmart/product-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/product-service/internal/infrastructure/adapter/repository"
)

func main() {
	port := getEnv("PORT", "8082")
	dbURL := getEnv("DATABASE_URL", "postgres://cloudmart:cloudmart_secret@localhost:5432/cloudmart?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "localhost:6379")
	natsURL := getEnv("NATS_URL", "nats://localhost:4222")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		panic(fmt.Sprintf("postgres: %v", err))
	}
	defer pool.Close()

	rdb := redis.NewClient(&redis.Options{Addr: redisURL})
	defer rdb.Close()

	nc, err := nats.Connect(natsURL)
	if err != nil {
		panic(fmt.Sprintf("nats: %v", err))
	}
	defer nc.Close()

	// Wire hexagonal architecture
	productRepo := repository.NewProductPostgresRepo(pool)
	categoryRepo := repository.NewCategoryPostgresRepo(pool)
	reviewRepo := repository.NewReviewPostgresRepo(pool)

	// Simple cache & event adapters (reuse patterns from user-service)
	cacheAdapter := &simpleCache{client: rdb}
	eventPub := &simplePublisher{conn: nc}

	productService := appservice.NewProductService(productRepo, categoryRepo, reviewRepo, cacheAdapter, eventPub)
	productHandler := handler.NewProductHandler(productService)

	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"healthy","service":"product-service"}`))
	})
	productHandler.RegisterRoutes(r)

	srv := &http.Server{Addr: ":" + port, Handler: r, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second}

	go func() {
		fmt.Printf("🚀 Product Service starting on port %s\n", port)
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
}

// ── Simple adapters (inline for brevity) ──

type simpleCache struct{ client *redis.Client }

func (c *simpleCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}
func (c *simpleCache) Set(ctx context.Context, key string, value interface{}, ttl int) error {
	return c.client.Set(ctx, key, value, time.Duration(ttl)*time.Second).Err()
}
func (c *simpleCache) Delete(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

type simplePublisher struct{ conn *nats.Conn }

func (p *simplePublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	payload, _ := json.Marshal(data)
	return p.conn.Publish(subject, payload)
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}
