package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudmart/inventory-service/internal/domain/model"
	"github.com/cloudmart/inventory-service/internal/domain/port"
	"github.com/google/uuid"
)

var (
	ErrInsufficientStock = errors.New("insufficient stock available")
	ErrStockNotFound     = errors.New("stock not found for product")
)

type inventoryService struct {
	stockRepo    port.StockRepository
	movementRepo port.MovementRepository
	events       port.EventPublisher
}

func NewInventoryService(stockRepo port.StockRepository, movementRepo port.MovementRepository, events port.EventPublisher) port.InventoryService {
	return &inventoryService{stockRepo: stockRepo, movementRepo: movementRepo, events: events}
}

func (s *inventoryService) GetStock(ctx context.Context, productID uuid.UUID) (*model.Stock, error) {
	stock, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, ErrStockNotFound
	}
	return stock, nil
}

func (s *inventoryService) ListStock(ctx context.Context, page, limit int) ([]model.Stock, int, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 20 }
	return s.stockRepo.List(ctx, page, limit)
}

func (s *inventoryService) Reserve(ctx context.Context, productID uuid.UUID, quantity int, orderID uuid.UUID) error {
	stock, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return ErrStockNotFound
	}
	if stock.Available() < quantity {
		return ErrInsufficientStock
	}

	if err := s.stockRepo.Reserve(ctx, productID, quantity); err != nil {
		return fmt.Errorf("reserve stock: %w", err)
	}

	_ = s.movementRepo.Create(ctx, &model.Movement{
		ID:            uuid.New(),
		ProductID:     productID,
		Type:          model.Reservation,
		Quantity:      quantity,
		ReferenceID:   &orderID,
		ReferenceType: "order",
		CreatedAt:     time.Now(),
	})

	// Check if reorder needed
	stock, _ = s.stockRepo.FindByProductID(ctx, productID)
	if stock != nil && stock.NeedsReorder() {
		_ = s.events.Publish(ctx, "inventory.low_stock", map[string]interface{}{
			"product_id": productID, "available": stock.Available(), "reorder_level": stock.ReorderLevel,
		})
	}

	return nil
}

func (s *inventoryService) ReleaseReservation(ctx context.Context, productID uuid.UUID, quantity int, orderID uuid.UUID) error {
	if err := s.stockRepo.Release(ctx, productID, quantity); err != nil {
		return fmt.Errorf("release reservation: %w", err)
	}
	_ = s.movementRepo.Create(ctx, &model.Movement{
		ID: uuid.New(), ProductID: productID, Type: model.Release,
		Quantity: quantity, ReferenceID: &orderID, ReferenceType: "order", CreatedAt: time.Now(),
	})
	return nil
}

func (s *inventoryService) UpdateStock(ctx context.Context, productID uuid.UUID, req port.UpdateStockRequest) (*model.Stock, error) {
	stock, err := s.stockRepo.FindByProductID(ctx, productID)
	if err != nil {
		return nil, ErrStockNotFound
	}
	if req.ReorderLevel != nil {
		stock.ReorderLevel = *req.ReorderLevel
	}
	if req.ReorderQuantity != nil {
		stock.ReorderQuantity = *req.ReorderQuantity
	}
	stock.UpdatedAt = time.Now()
	if err := s.stockRepo.Update(ctx, stock); err != nil {
		return nil, err
	}
	return stock, nil
}

func (s *inventoryService) Restock(ctx context.Context, productID uuid.UUID, quantity int) error {
	if err := s.stockRepo.AddStock(ctx, productID, quantity); err != nil {
		return fmt.Errorf("restock: %w", err)
	}
	_ = s.movementRepo.Create(ctx, &model.Movement{
		ID: uuid.New(), ProductID: productID, Type: model.Inbound,
		Quantity: quantity, Notes: "restocked", CreatedAt: time.Now(),
	})
	return nil
}

func (s *inventoryService) GetMovements(ctx context.Context, productID uuid.UUID) ([]model.Movement, error) {
	return s.movementRepo.FindByProductID(ctx, productID)
}
