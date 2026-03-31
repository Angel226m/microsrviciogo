// ═══════════════════════════════════════════════════════════════
// Repositorio PostgreSQL – Adaptador secundario de Transacciones
// Implementa port.TransactionRepository usando PostgreSQL (pgx/v5)
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"encoding/json"

	"github.com/cloudmart/payment-service/internal/domain/model"
	"github.com/cloudmart/payment-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type txPostgresRepo struct{ pool *pgxpool.Pool }

func NewTransactionPostgresRepo(pool *pgxpool.Pool) port.TransactionRepository {
	return &txPostgresRepo{pool: pool}
}

func (r *txPostgresRepo) Create(ctx context.Context, tx *model.Transaction) error {
	meta, _ := json.Marshal(tx.Metadata)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO payments.transactions (id, order_id, user_id, amount, currency, method, status, provider, provider_tx_id, metadata, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		tx.ID, tx.OrderID, tx.UserID, tx.Amount, tx.Currency, string(tx.Method),
		string(tx.Status), tx.Provider, tx.ProviderTxID, meta, tx.CreatedAt, tx.UpdatedAt)
	return err
}

func (r *txPostgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var metaJSON []byte
	err := r.pool.QueryRow(ctx,
		`SELECT id, order_id, user_id, amount, currency, method, status, provider, provider_tx_id, metadata, created_at, updated_at
		FROM payments.transactions WHERE id = $1`, id).Scan(
		&tx.ID, &tx.OrderID, &tx.UserID, &tx.Amount, &tx.Currency, &tx.Method,
		&tx.Status, &tx.Provider, &tx.ProviderTxID, &metaJSON, &tx.CreatedAt, &tx.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(metaJSON, &tx.Metadata)
	return tx, nil
}

func (r *txPostgresRepo) FindByOrderID(ctx context.Context, orderID uuid.UUID) (*model.Transaction, error) {
	tx := &model.Transaction{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, order_id, user_id, amount, currency, method, status, provider, provider_tx_id, created_at
		FROM payments.transactions WHERE order_id = $1 ORDER BY created_at DESC LIMIT 1`, orderID).Scan(
		&tx.ID, &tx.OrderID, &tx.UserID, &tx.Amount, &tx.Currency, &tx.Method,
		&tx.Status, &tx.Provider, &tx.ProviderTxID, &tx.CreatedAt)
	return tx, err
}

func (r *txPostgresRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.PaymentStatus) error {
	_, err := r.pool.Exec(ctx, "UPDATE payments.transactions SET status = $2, updated_at = NOW() WHERE id = $1", id, string(status))
	return err
}
