// ═══════════════════════════════════════════════════════════════
// Servicio de Pedidos – Punto de Entrada
// Gestiona el ciclo de vida completo de los pedidos del marketplace
// ═══════════════════════════════════════════════════════════════
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

	appservice "github.com/cloudmart/order-service/internal/application/service"
	handler "github.com/cloudmart/order-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/order-service/internal/infrastructure/adapter/repository"
)

func main() {
	port := getEnv("PORT", "8083")
	ctx := context.Background()

	pool, _ := pgxpool.New(ctx, getEnv("DATABASE_URL", "postgres://cloudmart:cloudmart_secret@localhost:5432/cloudmart?sslmode=disable"))
	defer pool.Close()

	nc, _ := nats.Connect(getEnv("NATS_URL", "nats://localhost:4222"))
	defer nc.Close()

	orderRepo := repository.NewOrderPostgresRepo(pool)
	eventPub := &natsPublisher{conn: nc}
	orderService := appservice.NewOrderService(orderRepo, eventPub)
	orderHandler := handler.NewOrderHandler(orderService)

	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"healthy","service":"order-service"}`))
	})
	orderHandler.RegisterRoutes(r)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() {
		fmt.Printf("🚀 Order Service starting on port %s\n", port)
		srv.ListenAndServe()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
}

type natsPublisher struct{ conn *nats.Conn }

func (p *natsPublisher) Publish(ctx context.Context, subject string, data interface{}) error {
	payload, _ := json.Marshal(data)
	return p.conn.Publish(subject, payload)
}

func getEnv(k, d string) string {
	if v, ok := os.LookupEnv(k); ok {
		return v
	}
	return d
}
