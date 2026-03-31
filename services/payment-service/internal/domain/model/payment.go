// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidades de Pago
// Transacciones, estados y métodos de pago soportados
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// PaymentStatus representa el estado de una transacción de pago.
type PaymentStatus string

// PaymentMethod representa un método de pago soportado.
type PaymentMethod string

const (
	PaymentPending    PaymentStatus = "pending"
	PaymentProcessing PaymentStatus = "processing"
	PaymentCompleted  PaymentStatus = "completed"
	PaymentFailed     PaymentStatus = "failed"
	PaymentRefunded   PaymentStatus = "refunded"

	MethodCreditCard   PaymentMethod = "credit_card"
	MethodDebitCard    PaymentMethod = "debit_card"
	MethodPayPal       PaymentMethod = "paypal"
	MethodBankTransfer PaymentMethod = "bank_transfer"
)

type Transaction struct {
	ID           uuid.UUID              `json:"id"`
	OrderID      uuid.UUID              `json:"order_id"`
	UserID       uuid.UUID              `json:"user_id"`
	Amount       float64                `json:"amount"`
	Currency     string                 `json:"currency"`
	Method       PaymentMethod          `json:"method"`
	Status       PaymentStatus          `json:"status"`
	Provider     string                 `json:"provider"`
	ProviderTxID string                 `json:"provider_tx_id,omitempty"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}
