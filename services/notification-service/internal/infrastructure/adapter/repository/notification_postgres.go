// ═══════════════════════════════════════════════════════════════
// Repositorio PostgreSQL – Adaptador secundario de Notificaciones
// Implementa port.NotificationRepository usando PostgreSQL (pgx/v5)
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"time"

	"github.com/cloudmart/notification-service/internal/domain/model"
	"github.com/cloudmart/notification-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type notificationRepo struct{ pool *pgxpool.Pool }

func NewNotificationPostgresRepo(pool *pgxpool.Pool) port.NotificationRepository {
	return &notificationRepo{pool: pool}
}

func (r *notificationRepo) Create(ctx context.Context, n *model.Notification) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO notifications.notifications (id, user_id, type, channel, subject, body, status, retry_count, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		n.ID, n.UserID, string(n.Type), n.Channel, n.Subject, n.Body, string(n.Status), n.RetryCount, n.CreatedAt)
	return err
}

func (r *notificationRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.NotificationStatus) error {
	var sentAt *time.Time
	if status == model.StatusSent {
		now := time.Now()
		sentAt = &now
	}
	_, err := r.pool.Exec(ctx,
		"UPDATE notifications.notifications SET status = $2, sent_at = $3 WHERE id = $1",
		id, string(status), sentAt)
	return err
}

func (r *notificationRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Notification, error) {
	n := &model.Notification{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, type, channel, subject, body, status, retry_count, sent_at, created_at
		FROM notifications.notifications WHERE id = $1`, id).Scan(
		&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Subject, &n.Body, &n.Status, &n.RetryCount, &n.SentAt, &n.CreatedAt)
	return n, err
}

func (r *notificationRepo) FindByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Notification, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM notifications.notifications WHERE user_id = $1", userID).Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, type, channel, subject, body, status, retry_count, sent_at, created_at
		FROM notifications.notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var notifs []model.Notification
	for rows.Next() {
		var n model.Notification
		rows.Scan(&n.ID, &n.UserID, &n.Type, &n.Channel, &n.Subject, &n.Body, &n.Status, &n.RetryCount, &n.SentAt, &n.CreatedAt)
		notifs = append(notifs, n)
	}
	return notifs, total, nil
}
