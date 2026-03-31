// ═══════════════════════════════════════════════════════════════
// Capa de Aplicación – Servicio de Pagos (casos de uso)
// Orquesta procesamiento de pagos, reembolsos y eventos
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudmart/payment-service/internal/domain/model"
	"github.com/cloudmart/payment-service/internal/domain/port"
	"github.com/google/uuid"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrPaymentFailed       = errors.New("payment processing failed")
)

type paymentService struct {
	repo    port.TransactionRepository
	gateway port.PaymentGateway
	events  port.EventPublisher
}

func NewPaymentService(repo port.TransactionRepository, gateway port.PaymentGateway, events port.EventPublisher) port.PaymentService {
	return &paymentService{repo: repo, gateway: gateway, events: events}
}

func (s *paymentService) ProcessPayment(ctx context.Context, req port.ProcessPaymentRequest) (*model.Transaction, error) {
	tx := &model.Transaction{
		ID:        uuid.New(),
		OrderID:   req.OrderID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Method:    req.Method,
		Status:    model.PaymentProcessing,
		Provider:  "stripe",
		Metadata:  req.Metadata,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transaction: %w", err)
	}

	providerTxID, err := s.gateway.Charge(ctx, req.Amount, req.Currency, req.Method, req.Metadata)
	if err != nil {
		tx.Status = model.PaymentFailed
		s.repo.UpdateStatus(ctx, tx.ID, model.PaymentFailed)
		_ = s.events.Publish(ctx, "payment.failed", map[string]interface{}{
			"transaction_id": tx.ID, "order_id": tx.OrderID, "error": err.Error(),
		})
		return tx, ErrPaymentFailed
	}

	tx.ProviderTxID = providerTxID
	tx.Status = model.PaymentCompleted
	s.repo.UpdateStatus(ctx, tx.ID, model.PaymentCompleted)

	_ = s.events.Publish(ctx, "payment.completed", map[string]interface{}{
		"transaction_id": tx.ID, "order_id": tx.OrderID, "amount": tx.Amount,
	})

	return tx, nil
}

func (s *paymentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *paymentService) RefundPayment(ctx context.Context, id uuid.UUID) error {
	tx, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrTransactionNotFound
	}

	if err := s.gateway.Refund(ctx, tx.ProviderTxID); err != nil {
		return fmt.Errorf("refund: %w", err)
	}

	s.repo.UpdateStatus(ctx, id, model.PaymentRefunded)
	_ = s.events.Publish(ctx, "payment.refunded", map[string]interface{}{
		"transaction_id": id, "order_id": tx.OrderID,
	})
	return nil
}
