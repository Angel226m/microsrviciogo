// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidades de Notificación
// Tipos de notificación (email, SMS, push) y estados de envío
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// NotificationType representa el canal de envío de la notificación.
type NotificationType string

const (
	NotificationEmail NotificationType = "email"
	NotificationSMS   NotificationType = "sms"
	NotificationPush  NotificationType = "push"
)

type NotificationStatus string

const (
	StatusPending NotificationStatus = "pending"
	StatusSent    NotificationStatus = "sent"
	StatusFailed  NotificationStatus = "failed"
)

type Notification struct {
	ID         uuid.UUID          `json:"id"`
	UserID     uuid.UUID          `json:"user_id"`
	Type       NotificationType   `json:"type"`
	Channel    string             `json:"channel"` // email address, phone, device token
	Subject    string             `json:"subject"`
	Body       string             `json:"body"`
	Status     NotificationStatus `json:"status"`
	RetryCount int                `json:"retry_count"`
	SentAt     *time.Time         `json:"sent_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

type Template struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
