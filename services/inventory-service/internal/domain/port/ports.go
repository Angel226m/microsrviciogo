package port

import (
	"context"

	"github.com/cloudmart/inventory-service/internal/domain/model"
	"github.com/google/uuid"
)

type InventoryService interface {
	GetStock(ctx context.Context, productID uuid.UUID) (*model.Stock, error)
	ListStock(ctx context.Context, page, limit int) ([]model.Stock, int, error)
	Reserve(ctx context.Context, productID uuid.UUID, quantity int, orderID uuid.UUID) error
	ReleaseReservation(ctx context.Context, productID uuid.UUID, quantity int, orderID uuid.UUID) error
	UpdateStock(ctx context.Context, productID uuid.UUID, req UpdateStockRequest) (*model.Stock, error)
	Restock(ctx context.Context, productID uuid.UUID, quantity int) error
	GetMovements(ctx context.Context, productID uuid.UUID) ([]model.Movement, error)
}

type StockRepository interface {
	FindByProductID(ctx context.Context, productID uuid.UUID) (*model.Stock, error)
	List(ctx context.Context, page, limit int) ([]model.Stock, int, error)
	Reserve(ctx context.Context, productID uuid.UUID, quantity int) error
	Release(ctx context.Context, productID uuid.UUID, quantity int) error
	Deduct(ctx context.Context, productID uuid.UUID, quantity int) error
	AddStock(ctx context.Context, productID uuid.UUID, quantity int) error
	Update(ctx context.Context, stock *model.Stock) error
}

type MovementRepository interface {
	Create(ctx context.Context, movement *model.Movement) error
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]model.Movement, error)
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

type UpdateStockRequest struct {
	ReorderLevel    *int `json:"reorder_level,omitempty"`
	ReorderQuantity *int `json:"reorder_quantity,omitempty"`
}
