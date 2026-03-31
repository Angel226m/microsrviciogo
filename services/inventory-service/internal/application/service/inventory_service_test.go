// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Capa de Aplicación del Servicio de Inventario
// Valida reservas, liberaciones, reabastecimiento y stock bajo
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/cloudmart/inventory-service/internal/domain/model"
	"github.com/cloudmart/inventory-service/internal/domain/port"
	"github.com/google/uuid"
)

// ── Mocks manuales de puertos secundarios ─────────────────────

type mockStockRepo struct {
	stocks map[uuid.UUID]*model.Stock
	err    error
}

func newMockStockRepo() *mockStockRepo {
	return &mockStockRepo{stocks: make(map[uuid.UUID]*model.Stock)}
}

func (m *mockStockRepo) FindByProductID(_ context.Context, productID uuid.UUID) (*model.Stock, error) {
	if s, ok := m.stocks[productID]; ok {
		return s, nil
	}
	return nil, errors.New("no encontrado")
}
func (m *mockStockRepo) List(_ context.Context, _, _ int) ([]model.Stock, int, error) {
	var lista []model.Stock
	for _, s := range m.stocks {
		lista = append(lista, *s)
	}
	return lista, len(lista), m.err
}
func (m *mockStockRepo) Reserve(_ context.Context, productID uuid.UUID, qty int) error {
	if s, ok := m.stocks[productID]; ok {
		s.Reserved += qty
		return nil
	}
	return errors.New("no encontrado")
}
func (m *mockStockRepo) Release(_ context.Context, productID uuid.UUID, qty int) error {
	if s, ok := m.stocks[productID]; ok {
		s.Reserved -= qty
		return nil
	}
	return errors.New("no encontrado")
}
func (m *mockStockRepo) Deduct(_ context.Context, productID uuid.UUID, qty int) error {
	if s, ok := m.stocks[productID]; ok {
		s.Quantity -= qty
		return nil
	}
	return errors.New("no encontrado")
}
func (m *mockStockRepo) AddStock(_ context.Context, productID uuid.UUID, qty int) error {
	if s, ok := m.stocks[productID]; ok {
		s.Quantity += qty
		return m.err
	}
	return errors.New("no encontrado")
}
func (m *mockStockRepo) Update(_ context.Context, stock *model.Stock) error {
	m.stocks[stock.ProductID] = stock
	return m.err
}

type mockMovementRepo struct {
	movimientos []model.Movement
}

func (m *mockMovementRepo) Create(_ context.Context, mov *model.Movement) error {
	m.movimientos = append(m.movimientos, *mov)
	return nil
}
func (m *mockMovementRepo) FindByProductID(_ context.Context, productID uuid.UUID) ([]model.Movement, error) {
	var resultado []model.Movement
	for _, mov := range m.movimientos {
		if mov.ProductID == productID {
			resultado = append(resultado, mov)
		}
	}
	return resultado, nil
}

type mockInventoryEvents struct {
	published []string
}

func (m *mockInventoryEvents) Publish(_ context.Context, subject string, _ interface{}) error {
	m.published = append(m.published, subject)
	return nil
}

// ── Helpers ───────────────────────────────────────────────────

func crearStock(productID uuid.UUID, cantidad, reservado, nivelReorden int) *model.Stock {
	return &model.Stock{
		ID:           uuid.New(),
		ProductID:    productID,
		SKU:          "TEST-001",
		Quantity:     cantidad,
		Reserved:     reservado,
		Warehouse:    "principal",
		ReorderLevel: nivelReorden,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// ── Pruebas ───────────────────────────────────────────────────

func TestInventoryService_GetStock_Exito(t *testing.T) {
	repo := newMockStockRepo()
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 100, 10, 5)

	svc := NewInventoryService(repo, &mockMovementRepo{}, &mockInventoryEvents{})

	stock, err := svc.GetStock(context.Background(), productID)
	if err != nil {
		t.Fatalf("GetStock() error inesperado: %v", err)
	}
	if stock.Quantity != 100 {
		t.Errorf("Quantity = %d, esperado 100", stock.Quantity)
	}
}

func TestInventoryService_GetStock_NoEncontrado(t *testing.T) {
	svc := NewInventoryService(newMockStockRepo(), &mockMovementRepo{}, &mockInventoryEvents{})

	_, err := svc.GetStock(context.Background(), uuid.New())
	if !errors.Is(err, ErrStockNotFound) {
		t.Errorf("GetStock() error = %v, esperado ErrStockNotFound", err)
	}
}

func TestInventoryService_Reserve_Exito(t *testing.T) {
	repo := newMockStockRepo()
	movRepo := &mockMovementRepo{}
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 50, 0, 5)

	svc := NewInventoryService(repo, movRepo, &mockInventoryEvents{})

	err := svc.Reserve(context.Background(), productID, 10, uuid.New())
	if err != nil {
		t.Fatalf("Reserve() error inesperado: %v", err)
	}
	if repo.stocks[productID].Reserved != 10 {
		t.Errorf("Reserved = %d, esperado 10", repo.stocks[productID].Reserved)
	}
	if len(movRepo.movimientos) != 1 {
		t.Fatalf("Se esperaba 1 movimiento, obtuve %d", len(movRepo.movimientos))
	}
	if movRepo.movimientos[0].Type != model.Reservation {
		t.Errorf("Tipo de movimiento = %s, esperado reservation", movRepo.movimientos[0].Type)
	}
}

func TestInventoryService_Reserve_StockInsuficiente(t *testing.T) {
	repo := newMockStockRepo()
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 10, 8, 5) // solo 2 disponibles

	svc := NewInventoryService(repo, &mockMovementRepo{}, &mockInventoryEvents{})

	err := svc.Reserve(context.Background(), productID, 5, uuid.New())
	if !errors.Is(err, ErrInsufficientStock) {
		t.Errorf("Reserve() error = %v, esperado ErrInsufficientStock", err)
	}
}

func TestInventoryService_Reserve_AlertaStockBajo(t *testing.T) {
	repo := newMockStockRepo()
	eventos := &mockInventoryEvents{}
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 15, 0, 10) // 15 disp, nivel 10

	svc := NewInventoryService(repo, &mockMovementRepo{}, eventos)

	// Reservar 6 → quedan 9 disponibles (< nivel 10) → alerta
	err := svc.Reserve(context.Background(), productID, 6, uuid.New())
	if err != nil {
		t.Fatalf("Reserve() error inesperado: %v", err)
	}
	encontrado := false
	for _, e := range eventos.published {
		if e == "inventory.low_stock" {
			encontrado = true
			break
		}
	}
	if !encontrado {
		t.Error("Se esperaba evento 'inventory.low_stock' al bajar del nivel de reorden")
	}
}

func TestInventoryService_ReleaseReservation(t *testing.T) {
	repo := newMockStockRepo()
	movRepo := &mockMovementRepo{}
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 50, 20, 5)

	svc := NewInventoryService(repo, movRepo, &mockInventoryEvents{})

	err := svc.ReleaseReservation(context.Background(), productID, 10, uuid.New())
	if err != nil {
		t.Fatalf("ReleaseReservation() error inesperado: %v", err)
	}
	if repo.stocks[productID].Reserved != 10 {
		t.Errorf("Reserved = %d, esperado 10 (de 20 liberamos 10)", repo.stocks[productID].Reserved)
	}
	if len(movRepo.movimientos) != 1 || movRepo.movimientos[0].Type != model.Release {
		t.Error("Se esperaba un movimiento de tipo 'release'")
	}
}

func TestInventoryService_Restock(t *testing.T) {
	repo := newMockStockRepo()
	movRepo := &mockMovementRepo{}
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 30, 0, 10)

	svc := NewInventoryService(repo, movRepo, &mockInventoryEvents{})

	err := svc.Restock(context.Background(), productID, 50)
	if err != nil {
		t.Fatalf("Restock() error inesperado: %v", err)
	}
	if repo.stocks[productID].Quantity != 80 {
		t.Errorf("Quantity = %d, esperado 80 (30 + 50)", repo.stocks[productID].Quantity)
	}
	if len(movRepo.movimientos) != 1 || movRepo.movimientos[0].Type != model.Inbound {
		t.Error("Se esperaba un movimiento de tipo 'inbound'")
	}
}

func TestInventoryService_UpdateStock(t *testing.T) {
	repo := newMockStockRepo()
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 100, 0, 5)

	svc := NewInventoryService(repo, &mockMovementRepo{}, &mockInventoryEvents{})

	nuevoNivel := 20
	nuevaCantidad := 50
	stock, err := svc.UpdateStock(context.Background(), productID, port.UpdateStockRequest{
		ReorderLevel:    &nuevoNivel,
		ReorderQuantity: &nuevaCantidad,
	})
	if err != nil {
		t.Fatalf("UpdateStock() error inesperado: %v", err)
	}
	if stock.ReorderLevel != 20 {
		t.Errorf("ReorderLevel = %d, esperado 20", stock.ReorderLevel)
	}
	if stock.ReorderQuantity != 50 {
		t.Errorf("ReorderQuantity = %d, esperado 50", stock.ReorderQuantity)
	}
}

func TestInventoryService_ListStock_PaginacionPorDefecto(t *testing.T) {
	repo := newMockStockRepo()
	productID := uuid.New()
	repo.stocks[productID] = crearStock(productID, 100, 0, 5)

	svc := NewInventoryService(repo, &mockMovementRepo{}, &mockInventoryEvents{})

	stocks, total, err := svc.ListStock(context.Background(), 0, 0) // valores inválidos
	if err != nil {
		t.Fatalf("ListStock() error inesperado: %v", err)
	}
	if total != 1 {
		t.Errorf("Total = %d, esperado 1", total)
	}
	if len(stocks) != 1 {
		t.Errorf("len(stocks) = %d, esperado 1", len(stocks))
	}
}
