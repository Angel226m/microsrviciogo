package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	appService "github.com/cloudmart/notification-service/internal/application/service"
	"github.com/cloudmart/notification-service/internal/domain/model"
	"github.com/cloudmart/notification-service/internal/domain/port"
	emailAdapter "github.com/cloudmart/notification-service/internal/infrastructure/adapter/email"
	httpHandler "github.com/cloudmart/notification-service/internal/infrastructure/adapter/http"
	"github.com/cloudmart/notification-service/internal/infrastructure/adapter/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, envOrDefault("DATABASE_URL", "postgres://cloudmart:cloudmart@localhost:5432/cloudmart?sslmode=disable"))
	if err != nil {
		log.Fatalf("db connect: %v", err)
	}
	defer pool.Close()

	nc, err := nats.Connect(envOrDefault("NATS_URL", "nats://localhost:4222"))
	if err != nil {
		log.Fatalf("nats connect: %v", err)
	}
	defer nc.Close()

	smtpPort, _ := strconv.Atoi(envOrDefault("SMTP_PORT", "1025"))
	emailSender := emailAdapter.NewSMTPSender(
		envOrDefault("SMTP_HOST", "localhost"),
		smtpPort,
		envOrDefault("SMTP_USER", ""),
		envOrDefault("SMTP_PASS", ""),
		envOrDefault("SMTP_FROM", "noreply@cloudmart.dev"),
	)

	repo := repository.NewNotificationPostgresRepo(pool)
	svc := appService.NewNotificationService(repo, emailSender)
	handler := httpHandler.NewHandler(svc)

	// Subscribe to domain events from other services
	nc.Subscribe("users.registered", func(msg *nats.Msg) {
		var evt struct {
			UserID string `json:"user_id"`
			Email  string `json:"email"`
			Name   string `json:"name"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			return
		}
		uid, _ := uuid.Parse(evt.UserID)
		svc.Send(ctx, port.SendNotificationRequest{
			UserID:  uid,
			Type:    model.NotificationEmail,
			Channel: evt.Email,
			Subject: "Welcome to CloudMart!",
			Body:    "<h1>Welcome " + evt.Name + "!</h1><p>Thank you for joining CloudMart.</p>",
		})
	})

	nc.Subscribe("orders.created", func(msg *nats.Msg) {
		var evt struct {
			UserID      string  `json:"user_id"`
			Email       string  `json:"email"`
			OrderNumber string  `json:"order_number"`
			Total       float64 `json:"total"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			return
		}
		uid, _ := uuid.Parse(evt.UserID)
		svc.Send(ctx, port.SendNotificationRequest{
			UserID:  uid,
			Type:    model.NotificationEmail,
			Channel: evt.Email,
			Subject: "Order Confirmation - " + evt.OrderNumber,
			Body:    "<h1>Order Confirmed</h1><p>Your order <strong>" + evt.OrderNumber + "</strong> has been received.</p>",
		})
	})

	nc.Subscribe("payments.completed", func(msg *nats.Msg) {
		var evt struct {
			UserID        string  `json:"user_id"`
			Email         string  `json:"email"`
			TransactionID string  `json:"transaction_id"`
			Amount        float64 `json:"amount"`
		}
		if err := json.Unmarshal(msg.Data, &evt); err != nil {
			return
		}
		uid, _ := uuid.Parse(evt.UserID)
		svc.Send(ctx, port.SendNotificationRequest{
			UserID:  uid,
			Type:    model.NotificationEmail,
			Channel: evt.Email,
			Subject: "Payment Received",
			Body:    "<h1>Payment Confirmed</h1><p>We received your payment for transaction " + evt.TransactionID + ".</p>",
		})
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

	port := envOrDefault("PORT", "8086")
	srv := &http.Server{Addr: ":" + port, Handler: r}

	go func() {
		log.Printf("notification-service listening on :%s", port)
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
	log.Println("notification-service stopped")
}

func envOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
