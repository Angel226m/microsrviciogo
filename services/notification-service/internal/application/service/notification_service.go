// ═══════════════════════════════════════════════════════════════
// Capa de Aplicación – Servicio de Notificaciones (casos de uso)
// Orquesta el envío de notificaciones por email/SMS/push
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"time"

	"github.com/cloudmart/notification-service/internal/domain/model"
	"github.com/cloudmart/notification-service/internal/domain/port"
	"github.com/google/uuid"
)

type notificationService struct {
	repo  port.NotificationRepository
	email port.EmailSender
}

func NewNotificationService(repo port.NotificationRepository, email port.EmailSender) port.NotificationService {
	return &notificationService{repo: repo, email: email}
}

func (s *notificationService) Send(ctx context.Context, req port.SendNotificationRequest) (*model.Notification, error) {
	n := &model.Notification{
		ID:        uuid.New(),
		UserID:    req.UserID,
		Type:      req.Type,
		Channel:   req.Channel,
		Subject:   req.Subject,
		Body:      req.Body,
		Status:    model.StatusPending,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return nil, err
	}

	// Send based on type
	var sendErr error
	switch n.Type {
	case model.NotificationEmail:
		sendErr = s.email.SendEmail(ctx, n.Channel, n.Subject, n.Body)
	default:
		// SMS/Push: log only for now
	}

	if sendErr != nil {
		n.Status = model.StatusFailed
		s.repo.UpdateStatus(ctx, n.ID, model.StatusFailed)
		return n, sendErr
	}

	now := time.Now()
	n.Status = model.StatusSent
	n.SentAt = &now
	s.repo.UpdateStatus(ctx, n.ID, model.StatusSent)
	return n, nil
}

func (s *notificationService) GetByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *notificationService) ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Notification, int, error) {
	return s.repo.FindByUser(ctx, userID, page, limit)
}
