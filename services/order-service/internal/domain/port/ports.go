package port

import (
	"context"

	"github.com/cloudmart/order-service/internal/domain/model"
	"github.com/google/uuid"
)

// ── Driving Ports ─────────────────────────────────────────────────────

type OrderService interface {
	Create(ctx context.Context, req CreateOrderRequest) (*model.Order, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Order, int, error)
	Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error
}

// ── Driven Ports ──────────────────────────────────────────────────────

type OrderRepository interface {
	Create(ctx context.Context, order *model.Order) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.Order, error)
	FindByUserID(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Order, int, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.OrderStatus) error
	CreateItems(ctx context.Context, items []model.OrderItem) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// ── Request DTOs ──────────────────────────────────────────────────────

type CreateOrderRequest struct {
	UserID          uuid.UUID      `json:"user_id"`
	Items           []OrderItemReq `json:"items"`
	ShippingAddress model.Address  `json:"shipping_address"`
	Notes           string         `json:"notes,omitempty"`
}

type OrderItemReq struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductSKU  string    `json:"product_sku"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
}
