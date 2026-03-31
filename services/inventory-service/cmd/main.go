// ═══════════════════════════════════════════════════════════════
// Servicio de Inventario – Punto de Entrada
// Gestiona stock, reservas, reabastecimiento y movimientos de inventario
// Escucha eventos de pedidos vía NATS para reservar/liberar stock
// ═══════════════════════════════════════════════════════════════
package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	appService "github.com/cloudmart/inventory-service/internal/application/service"
	httpHandler "github.com/cloudmart/inventory-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/inventory-service/internal/infrastructure/adapter/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()

	dbURL := envOrDefault("DATABASE_URL", "postgres://cloudmart:cloudmart@localhost:5432/cloudmart?sslmode=disable")
	pool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	nc, err := nats.Connect(envOrDefault("NATS_URL", "nats://localhost:4222"))
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()

	stockRepo := repository.NewStockPostgresRepo(pool)
	movementRepo := repository.NewMovementPostgresRepo(pool)
	publisher := &natsPublisher{nc: nc}

	svc := appService.NewInventoryService(stockRepo, movementRepo, publisher)
	handler := httpHandler.NewHandler(svc)

	// Subscribe to order events
	nc.Subscribe("orders.created", func(msg *nats.Msg) {
		var evt struct {
			OrderID string `json:"order_id"`
			Items   []struct {
				ProductID string `json:"product_id"`
				Quantity  int    `json:"quantity"`
			} `json:"items"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			log.Printf("failed to parse order event: %v", err)
			return
		}
		for _, item := range evt.Items {
			pid, _ := uuid.Parse(item.ProductID)
			oid, _ := uuid.Parse(evt.OrderID)
			svc.Reserve(ctx, pid, item.Quantity, oid)
		}
	})

	nc.Subscribe("orders.cancelled", func(msg *nats.Msg) {
		var evt struct {
			OrderID string `json:"order_id"`
			Items   []struct {
				ProductID string `json:"product_id"`
				Quantity  int    `json:"quantity"`
			} `json:"items"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			return
		}
		for _, item := range evt.Items {
			pid, _ := uuid.Parse(item.ProductID)
			oid, _ := uuid.Parse(evt.OrderID)
			svc.ReleaseReservation(ctx, pid, item.Quantity, oid)
		}
	})

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})
	r.Handle("/metrics", promhttp.Handler())
	r.Route("/api/v1", handler.RegisterRoutes)

	port := envOrDefault("PORT", "8085")
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Printf("inventory-service listening on :%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	srv.Shutdown(shutdownCtx)
	log.Println("inventory-service stopped")
}

type natsPublisher struct{ nc *nats.Conn }

func (p *natsPublisher) Publish(ctx context.Context, subject string, event interface{}) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	return p.nc.Publish(subject, data)
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
