// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Capa de Aplicación del Servicio de Pedidos
// Valida creación, consulta, cancelación y reglas de autorización
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudmart/order-service/internal/domain/model"
	"github.com/cloudmart/order-service/internal/domain/port"
	"github.com/google/uuid"
)

// ── Mocks manuales de puertos secundarios ─────────────────────

type mockOrderRepo struct {
	pedidos       map[uuid.UUID]*model.Order
	statusUpdated model.OrderStatus
	createdItems  []model.OrderItem
	err           error
}

func newMockOrderRepo() *mockOrderRepo {
	return &mockOrderRepo{pedidos: make(map[uuid.UUID]*model.Order)}
}

func (m *mockOrderRepo) Create(_ context.Context, o *model.Order) error {
	if m.err != nil {
		return m.err
	}
	m.pedidos[o.ID] = o
	return nil
}
func (m *mockOrderRepo) FindByID(_ context.Context, id uuid.UUID) (*model.Order, error) {
	if p, ok := m.pedidos[id]; ok {
		return p, nil
	}
	return nil, errors.New("no encontrado")
}
func (m *mockOrderRepo) FindByUserID(_ context.Context, userID uuid.UUID, _, _ int) ([]model.Order, int, error) {
	var resultado []model.Order
	for _, p := range m.pedidos {
		if p.UserID == userID {
			resultado = append(resultado, *p)
		}
	}
	return resultado, len(resultado), m.err
}
func (m *mockOrderRepo) UpdateStatus(_ context.Context, id uuid.UUID, status model.OrderStatus) error {
	if p, ok := m.pedidos[id]; ok {
		p.Status = status
		m.statusUpdated = status
		return m.err
	}
	return errors.New("no encontrado")
}
func (m *mockOrderRepo) CreateItems(_ context.Context, items []model.OrderItem) error {
	m.createdItems = items
	return m.err
}

type mockOrderEvents struct {
	published []string
}

func (m *mockOrderEvents) Publish(_ context.Context, subject string, _ interface{}) error {
	m.published = append(m.published, subject)
	return nil
}

// ── Pruebas ───────────────────────────────────────────────────

func TestOrderService_Create_Exito(t *testing.T) {
	repo := newMockOrderRepo()
	eventos := &mockOrderEvents{}
	svc := NewOrderService(repo, eventos)

	req := port.CreateOrderRequest{
		UserID: uuid.New(),
		Items: []port.OrderItemReq{
			{
				ProductID:   uuid.New(),
				ProductName: "Camiseta Azul",
				ProductSKU:  "CAM-AZL-001",
				Quantity:    2,
				UnitPrice:   299.99,
			},
		},
		ShippingAddress: model.Address{
			Street:  "Av. Reforma 123",
			City:    "Ciudad de México",
			State:   "CDMX",
			ZipCode: "06600",
			Country: "MX",
		},
	}

	pedido, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error inesperado: %v", err)
	}
	if pedido.Status != model.StatusPending {
		t.Errorf("Status = %s, esperado pending", pedido.Status)
	}
	if pedido.Currency != "MXN" {
		t.Errorf("Currency = %s, esperado MXN", pedido.Currency)
	}
	if len(pedido.Items) != 1 {
		t.Fatalf("len(Items) = %d, esperado 1", len(pedido.Items))
	}
	if pedido.Items[0].TotalPrice != 599.98 {
		t.Errorf("Items[0].TotalPrice = %.2f, esperado 599.98", pedido.Items[0].TotalPrice)
	}
	if pedido.Total <= 0 {
		t.Error("Total debería ser mayor a 0 después de CalculateTotal()")
	}
	if len(eventos.published) == 0 || eventos.published[0] != "order.created" {
		t.Error("Se esperaba evento 'order.created'")
	}
}

func TestOrderService_Create_ErrorRepositorio(t *testing.T) {
	repo := newMockOrderRepo()
	repo.err = errors.New("disco lleno")
	svc := NewOrderService(repo, &mockOrderEvents{})

	_, err := svc.Create(context.Background(), port.CreateOrderRequest{
		UserID: uuid.New(),
		Items:  []port.OrderItemReq{{ProductID: uuid.New(), Quantity: 1, UnitPrice: 100}},
	})
	if err == nil {
		t.Fatal("Create() debería retornar error cuando el repositorio falla")
	}
}

func TestOrderService_GetByID_NoEncontrado(t *testing.T) {
	repo := newMockOrderRepo()
	svc := NewOrderService(repo, &mockOrderEvents{})

	_, err := svc.GetByID(context.Background(), uuid.New())
	if !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("GetByID() error = %v, esperado ErrOrderNotFound", err)
	}
}

func TestOrderService_ListByUser_PaginacionPorDefecto(t *testing.T) {
	repo := newMockOrderRepo()
	userID := uuid.New()

	// Agregar pedido directamente al mock
	pedidoID := uuid.New()
	repo.pedidos[pedidoID] = &model.Order{ID: pedidoID, UserID: userID, Status: model.StatusPending}

	svc := NewOrderService(repo, &mockOrderEvents{})

	pedidos, total, err := svc.ListByUser(context.Background(), userID, 0, 0) // valores inválidos
	if err != nil {
		t.Fatalf("ListByUser() error inesperado: %v", err)
	}
	if total != 1 {
		t.Errorf("Total = %d, esperado 1", total)
	}
	if len(pedidos) != 1 {
		t.Errorf("len(pedidos) = %d, esperado 1", len(pedidos))
	}
}

func TestOrderService_Cancel_Exito(t *testing.T) {
	repo := newMockOrderRepo()
	eventos := &mockOrderEvents{}
	svc := NewOrderService(repo, eventos)

	userID := uuid.New()
	pedidoID := uuid.New()
	repo.pedidos[pedidoID] = &model.Order{
		ID:     pedidoID,
		UserID: userID,
		Status: model.StatusPending,
	}

	err := svc.Cancel(context.Background(), pedidoID, userID)
	if err != nil {
		t.Fatalf("Cancel() error inesperado: %v", err)
	}
	if repo.statusUpdated != model.StatusCancelled {
		t.Errorf("Estado actualizado = %s, esperado cancelled", repo.statusUpdated)
	}
	if len(eventos.published) == 0 || eventos.published[0] != "order.cancelled" {
		t.Error("Se esperaba evento 'order.cancelled'")
	}
}

func TestOrderService_Cancel_NoAutorizado(t *testing.T) {
	repo := newMockOrderRepo()
	svc := NewOrderService(repo, &mockOrderEvents{})

	dueno := uuid.New()
	intruso := uuid.New()
	pedidoID := uuid.New()
	repo.pedidos[pedidoID] = &model.Order{
		ID:     pedidoID,
		UserID: dueno,
		Status: model.StatusPending,
	}

	err := svc.Cancel(context.Background(), pedidoID, intruso)
	if !errors.Is(err, ErrUnauthorized) {
		t.Errorf("Cancel() error = %v, esperado ErrUnauthorized", err)
	}
}

func TestOrderService_Cancel_EstadoNoCancelable(t *testing.T) {
	repo := newMockOrderRepo()
	svc := NewOrderService(repo, &mockOrderEvents{})

	userID := uuid.New()
	estados := []model.OrderStatus{
		model.StatusShipped,
		model.StatusDelivered,
		model.StatusCancelled,
	}

	for _, estado := range estados {
		pedidoID := uuid.New()
		repo.pedidos[pedidoID] = &model.Order{
			ID:     pedidoID,
			UserID: userID,
			Status: estado,
		}

		err := svc.Cancel(context.Background(), pedidoID, userID)
		if !errors.Is(err, ErrCannotCancel) {
			t.Errorf("Cancel() con estado %s: error = %v, esperado ErrCannotCancel", estado, err)
		}
	}
}

func TestOrderService_Cancel_PedidoNoExiste(t *testing.T) {
	repo := newMockOrderRepo()
	svc := NewOrderService(repo, &mockOrderEvents{})

	err := svc.Cancel(context.Background(), uuid.New(), uuid.New())
	if !errors.Is(err, ErrOrderNotFound) {
		t.Errorf("Cancel() error = %v, esperado ErrOrderNotFound", err)
	}
}
