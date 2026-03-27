package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudmart/order-service/internal/domain/model"
	"github.com/cloudmart/order-service/internal/domain/port"
	"github.com/google/uuid"
)

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrCannotCancel     = errors.New("order cannot be cancelled in current status")
	ErrUnauthorized     = errors.New("not authorized to access this order")
)

type orderService struct {
	repo   port.OrderRepository
	events port.EventPublisher
}

func NewOrderService(repo port.OrderRepository, events port.EventPublisher) port.OrderService {
	return &orderService{repo: repo, events: events}
}

func (s *orderService) Create(ctx context.Context, req port.CreateOrderRequest) (*model.Order, error) {
	order := &model.Order{
		ID:              uuid.New(),
		OrderNumber:     model.GenerateOrderNumber(),
		UserID:          req.UserID,
		Status:          model.StatusPending,
		Currency:        "MXN",
		ShippingAddress: req.ShippingAddress,
		ShippingCost:    99.00,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	for _, item := range req.Items {
		order.Items = append(order.Items, model.OrderItem{
			ID:          uuid.New(),
			OrderID:     order.ID,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			ProductSKU:  item.ProductSKU,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			TotalPrice:  float64(item.Quantity) * item.UnitPrice,
			CreatedAt:   time.Now(),
		})
	}

	order.CalculateTotal()

	if err := s.repo.Create(ctx, order); err != nil {
		return nil, fmt.Errorf("create order: %w", err)
	}
	if err := s.repo.CreateItems(ctx, order.Items); err != nil {
		return nil, fmt.Errorf("create order items: %w", err)
	}

	_ = s.events.Publish(ctx, "order.created", map[string]interface{}{
		"order_id":     order.ID,
		"order_number": order.OrderNumber,
		"user_id":      order.UserID,
		"total":        order.Total,
		"items_count":  len(order.Items),
	})

	return order, nil
}

func (s *orderService) GetByID(ctx context.Context, id uuid.UUID) (*model.Order, error) {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}
	return order, nil
}

func (s *orderService) ListByUser(ctx context.Context, userID uuid.UUID, page, limit int) ([]model.Order, int, error) {
	if page < 1 { page = 1 }
	if limit < 1 || limit > 50 { limit = 10 }
	return s.repo.FindByUserID(ctx, userID, page, limit)
}

func (s *orderService) Cancel(ctx context.Context, id uuid.UUID, userID uuid.UUID) error {
	order, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.UserID != userID {
		return ErrUnauthorized
	}
	if !order.CanCancel() {
		return ErrCannotCancel
	}

	if err := s.repo.UpdateStatus(ctx, id, model.StatusCancelled); err != nil {
		return fmt.Errorf("cancel order: %w", err)
	}

	_ = s.events.Publish(ctx, "order.cancelled", map[string]interface{}{
		"order_id": order.ID,
		"user_id":  order.UserID,
	})

	return nil
}
