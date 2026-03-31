// ═══════════════════════════════════════════════════════════════
// Repositorio PostgreSQL – Adaptador secundario de Inventario
// Implementa port.StockRepository y port.MovementRepository
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudmart/inventory-service/internal/domain/model"
	"github.com/cloudmart/inventory-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type stockRepo struct{ pool *pgxpool.Pool }

func NewStockPostgresRepo(pool *pgxpool.Pool) port.StockRepository {
	return &stockRepo{pool: pool}
}

func (r *stockRepo) FindByProductID(ctx context.Context, productID uuid.UUID) (*model.Stock, error) {
	s := &model.Stock{}
	err := r.pool.QueryRow(ctx,
		`SELECT id, product_id, sku, quantity, reserved, warehouse, reorder_level, reorder_quantity, last_restocked_at, created_at, updated_at
		FROM inventory.stock WHERE product_id = $1`, productID).Scan(
		&s.ID, &s.ProductID, &s.SKU, &s.Quantity, &s.Reserved, &s.Warehouse,
		&s.ReorderLevel, &s.ReorderQuantity, &s.LastRestockedAt, &s.CreatedAt, &s.UpdatedAt)
	return s, err
}

func (r *stockRepo) List(ctx context.Context, page, limit int) ([]model.Stock, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM inventory.stock").Scan(&total)

	rows, err := r.pool.Query(ctx,
		`SELECT id, product_id, sku, quantity, reserved, warehouse, reorder_level, reorder_quantity, created_at, updated_at
		FROM inventory.stock ORDER BY sku LIMIT $1 OFFSET $2`, limit, (page-1)*limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stocks []model.Stock
	for rows.Next() {
		var s model.Stock
		rows.Scan(&s.ID, &s.ProductID, &s.SKU, &s.Quantity, &s.Reserved, &s.Warehouse,
			&s.ReorderLevel, &s.ReorderQuantity, &s.CreatedAt, &s.UpdatedAt)
		stocks = append(stocks, s)
	}
	return stocks, total, nil
}

func (r *stockRepo) Reserve(ctx context.Context, productID uuid.UUID, qty int) error {
	result, err := r.pool.Exec(ctx,
		"UPDATE inventory.stock SET reserved = reserved + $2, updated_at = NOW() WHERE product_id = $1 AND (quantity - reserved) >= $2",
		productID, qty)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("insufficient stock")
	}
	return nil
}

func (r *stockRepo) Release(ctx context.Context, productID uuid.UUID, qty int) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory.stock SET reserved = GREATEST(reserved - $2, 0), updated_at = NOW() WHERE product_id = $1",
		productID, qty)
	return err
}

func (r *stockRepo) Deduct(ctx context.Context, productID uuid.UUID, qty int) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory.stock SET quantity = quantity - $2, reserved = GREATEST(reserved - $2, 0), updated_at = NOW() WHERE product_id = $1",
		productID, qty)
	return err
}

func (r *stockRepo) AddStock(ctx context.Context, productID uuid.UUID, qty int) error {
	now := time.Now()
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory.stock SET quantity = quantity + $2, last_restocked_at = $3, updated_at = $3 WHERE product_id = $1",
		productID, qty, now)
	return err
}

func (r *stockRepo) Update(ctx context.Context, stock *model.Stock) error {
	_, err := r.pool.Exec(ctx,
		"UPDATE inventory.stock SET reorder_level = $2, reorder_quantity = $3, updated_at = $4 WHERE id = $1",
		stock.ID, stock.ReorderLevel, stock.ReorderQuantity, stock.UpdatedAt)
	return err
}

// Movement repository
type movementRepo struct{ pool *pgxpool.Pool }

func NewMovementPostgresRepo(pool *pgxpool.Pool) port.MovementRepository {
	return &movementRepo{pool: pool}
}

func (r *movementRepo) Create(ctx context.Context, m *model.Movement) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO inventory.movements (id, product_id, type, quantity, reference_id, reference_type, notes, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		m.ID, m.ProductID, string(m.Type), m.Quantity, m.ReferenceID, m.ReferenceType, m.Notes, m.CreatedAt)
	return err
}

func (r *movementRepo) FindByProductID(ctx context.Context, productID uuid.UUID) ([]model.Movement, error) {
	rows, err := r.pool.Query(ctx,
		"SELECT id, product_id, type, quantity, reference_id, reference_type, notes, created_at FROM inventory.movements WHERE product_id = $1 ORDER BY created_at DESC",
		productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movements []model.Movement
	for rows.Next() {
		var m model.Movement
		rows.Scan(&m.ID, &m.ProductID, &m.Type, &m.Quantity, &m.ReferenceID, &m.ReferenceType, &m.Notes, &m.CreatedAt)
		movements = append(movements, m)
	}
	return movements, nil
}
