// ═══════════════════════════════════════════════════════════════
// Puertos de Dominio – Límites hexagonales del servicio de notificaciones
// ═══════════════════════════════════════════════════════════════
package port

import (
	"context"

	"github.com/cloudmart/notification-service/internal/domain/model"
	"github.com/google/uuid"
)

// ── Puertos Primarios ─────────────────────────────────────────────────────

// NotificationService define los casos de uso del contexto acotado de notificaciones.
type NotificationService interface {
	Send(ctx context.Context, req SendNotificationRequest) (*model.Notification, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Notification, int, error)
}

type SendNotificationRequest struct {
	UserID  uuid.UUID              `json:"user_id"`
	Type    model.NotificationType `json:"type"`
	Channel string                 `json:"channel"`
	Subject string                 `json:"subject"`
	Body    string                 `json:"body"`
}

// Puertos secundarios
type NotificationRepository interface {
	Create(ctx context.Context, n *model.Notification) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.NotificationStatus) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Notification, error)
	FindByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Notification, int, error)
}

type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type TemplateRepository interface {
	FindByName(ctx context.Context, name string) (*model.Template, error)
}
