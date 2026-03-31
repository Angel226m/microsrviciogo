// ═══════════════════════════════════════════════════════════════
// Puertos de Dominio – Límites hexagonales del servicio de pagos
// Define interfaces para procesador, repositorio y publicador de eventos
// ═══════════════════════════════════════════════════════════════
package port

import (
	"context"

	"github.com/cloudmart/payment-service/internal/domain/model"
	"github.com/google/uuid"
)

// ── Puertos Primarios ─────────────────────────────────────────────────────

// PaymentService define los casos de uso del contexto acotado de pagos.
type PaymentService interface {
	ProcessPayment(ctx context.Context, req ProcessPaymentRequest) (*model.Transaction, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	RefundPayment(ctx context.Context, id uuid.UUID) error
}

type TransactionRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error)
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Transaction, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.PaymentStatus) error
}

type PaymentGateway interface {
	Charge(ctx context.Context, amount float64, currency string, method model.PaymentMethod, metadata map[string]interface{}) (providerTxID string, err error)
	Refund(ctx context.Context, providerTxID string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

type ProcessPaymentRequest struct {
	OrderID  uuid.UUID              `json:"order_id"`
	UserID   uuid.UUID              `json:"user_id"`
	Amount   float64                `json:"amount"`
	Currency string                 `json:"currency"`
	Method   model.PaymentMethod    `json:"method"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
