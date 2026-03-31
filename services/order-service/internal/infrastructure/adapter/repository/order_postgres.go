// ═══════════════════════════════════════════════════════════════
// Repositorio PostgreSQL – Adaptador secundario de Pedidos
// Implementa port.OrderRepository usando PostgreSQL (pgx/v5)
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"encoding/json"

	"github.com/cloudmart/order-service/internal/domain/model"
	"github.com/cloudmart/order-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type orderPostgresRepo struct{ pool *pgxpool.Pool }

func NewOrderPostgresRepo(pool *pgxpool.Pool) port.OrderRepository {
	return &orderPostgresRepo{pool: pool}
}

func (r *orderPostgresRepo) Create(ctx context.Context, order *model.Order) error {
	shippingJSON, _ := json.Marshal(order.ShippingAddress)
	billingJSON, _ := json.Marshal(order.BillingAddress)

	query := `INSERT INTO orders.orders (id, order_number, user_id, status, subtotal, tax, shipping_cost, discount, total, currency, shipping_address, billing_address, notes, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`

	_, err := r.pool.Exec(ctx, query, order.ID, order.OrderNumber, order.UserID,
		string(order.Status), order.Subtotal, order.Tax, order.ShippingCost, order.Discount,
		order.Total, order.Currency, shippingJSON, billingJSON, order.Notes,
		order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *orderPostgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	query := `SELECT id, order_number, user_id, status, subtotal, tax, shipping_cost, discount, total, currency, shipping_address, notes, created_at, updated_at
		FROM orders.orders WHERE id = $1`

	o := &model.Order{}
	var shippingJSON []byte
	err := r.pool.QueryRow(ctx, query, id).Scan(&o.ID, &o.OrderNumber, &o.UserID, &o.Status,
		&o.Subtotal, &o.Tax, &o.ShippingCost, &o.Discount, &o.Total, &o.Currency,
		&shippingJSON, &o.Notes, &o.CreatedAt, &o.UpdatedAt)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(shippingJSON, &o.ShippingAddress)

	// Load items
	itemRows, err := r.pool.Query(ctx, `SELECT id, order_id, product_id, product_name, product_sku, quantity, unit_price, total_price, created_at
		FROM orders.order_items WHERE order_id = $1`, id)
	if err == nil {
		defer itemRows.Close()
		for itemRows.Next() {
			var item model.OrderItem
			itemRows.Scan(&item.ID, &item.OrderID, &item.ProductID, &item.ProductName,
				&item.ProductSKU, &item.Quantity, &item.UnitPrice, &item.TotalPrice, &item.CreatedAt)
			o.Items = append(o.Items, item)
		}
	}

	return o, nil
}

func (r *orderPostgresRepo) FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Order, int, error) {
	var total int
	r.pool.QueryRow(ctx, "SELECT COUNT(*) FROM orders.orders WHERE user_id = $1", userID).Scan(&total)

	offset := (page - 1) * limit
	query := `SELECT id, order_number, user_id, status, subtotal, tax, shipping_cost, discount, total, currency, notes, created_at, updated_at
		FROM orders.orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []model.Order
	for rows.Next() {
		var o model.Order
		rows.Scan(&o.ID, &o.OrderNumber, &o.UserID, &o.Status, &o.Subtotal, &o.Tax,
			&o.ShippingCost, &o.Discount, &o.Total, &o.Currency, &o.Notes, &o.CreatedAt, &o.UpdatedAt)
		orders = append(orders, o)
	}
	return orders, total, nil
}

func (r *orderPostgresRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error {
	_, err := r.pool.Exec(ctx, "UPDATE orders.orders SET status = $2, updated_at = NOW() WHERE id = $1", id, string(status))
	return err
}

func (r *orderPostgresRepo) CreateItems(ctx context.Context, items []model.OrderItem) error {
	for _, item := range items {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO orders.order_items (id, order_id, product_id, product_name, product_sku, quantity, unit_price, total_price, created_at)
			VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
			item.ID, item.OrderID, item.ProductID, item.ProductName, item.ProductSKU,
			item.Quantity, item.UnitPrice, item.TotalPrice, item.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}
