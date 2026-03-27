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

	appservice "github.com/cloudmart/payment-service/internal/application/service"
	"github.com/cloudmart/payment-service/internal/infrastructure/adapter/gateway"
	handler "github.com/cloudmart/payment-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/payment-service/internal/infrastructure/adapter/repository"
)

func main() {
	port := env("PORT", "8084")
	ctx := context.Background()

	pool, _ := pgxpool.New(ctx, env("DATABASE_URL", "postgres://cloudmart:cloudmart_secret@localhost:5432/cloudmart?sslmode=disable"))
	defer pool.Close()
	nc, _ := nats.Connect(env("NATS_URL", "nats://localhost:4222"))
	defer nc.Close()

	txRepo := repository.NewTransactionPostgresRepo(pool)
	gw := gateway.NewMockPaymentGateway()
	pub := &natsPub{nc}
	svc := appservice.NewPaymentService(txRepo, gw, pub)
	h := handler.NewPaymentHandler(svc)

	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.Logger, chimw.Recoverer)
	r.Handle("/metrics", promhttp.Handler())
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"healthy","service":"payment-service"}`))
	})
	h.RegisterRoutes(r)

	srv := &http.Server{Addr: ":" + port, Handler: r}
	go func() { fmt.Printf("🚀 Payment Service on :%s\n", port); srv.ListenAndServe() }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	srv.Shutdown(c)
}

type natsPub struct{ c *nats.Conn }

func (p *natsPub) Publish(_ context.Context, subj string, data interface{}) error {
	b, _ := json.Marshal(data)
	return p.c.Publish(subj, b)
}
func env(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
