// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Capa de Aplicación del Servicio de Producto
// Valida casos de uso con mocks manuales de los puertos
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudmart/product-service/internal/domain/model"
	"github.com/cloudmart/product-service/internal/domain/port"
	"github.com/google/uuid"
)

// ── Mocks manuales de puertos secundarios ─────────────────────

type mockProductRepo struct {
	productos []model.Product
	total     int
	err       error
	created   *model.Product
	updated   *model.Product
	deleted   uuid.UUID
}

func (m *mockProductRepo) List(_ context.Context, _ model.ProductFilter) ([]model.Product, int, error) {
	return m.productos, m.total, m.err
}
func (m *mockProductRepo) FindByID(_ context.Context, id uuid.UUID) (*model.Product, error) {
	for i := range m.productos {
		if m.productos[i].ID == id {
			return &m.productos[i], nil
		}
	}
	return nil, errors.New("no encontrado")
}
func (m *mockProductRepo) FindBySlug(_ context.Context, slug string) (*model.Product, error) {
	for i := range m.productos {
		if m.productos[i].Slug == slug {
			return &m.productos[i], nil
		}
	}
	return nil, errors.New("no encontrado")
}
func (m *mockProductRepo) Create(_ context.Context, p *model.Product) error {
	m.created = p
	return m.err
}
func (m *mockProductRepo) Update(_ context.Context, p *model.Product) error {
	m.updated = p
	return m.err
}
func (m *mockProductRepo) Delete(_ context.Context, id uuid.UUID) error {
	m.deleted = id
	return m.err
}

type mockCategoryRepo struct {
	categorias []model.Category
	err        error
}

func (m *mockCategoryRepo) List(_ context.Context) ([]model.Category, error) {
	return m.categorias, m.err
}
func (m *mockCategoryRepo) FindByID(_ context.Context, _ uuid.UUID) (*model.Category, error) {
	if len(m.categorias) > 0 {
		return &m.categorias[0], nil
	}
	return nil, errors.New("no encontrado")
}

type mockReviewRepo struct {
	resenas []model.Review
	err     error
	created *model.Review
}

func (m *mockReviewRepo) FindByProductID(_ context.Context, _ uuid.UUID) ([]model.Review, error) {
	return m.resenas, m.err
}
func (m *mockReviewRepo) Create(_ context.Context, r *model.Review) error {
	m.created = r
	return m.err
}
func (m *mockReviewRepo) GetAverageRating(_ context.Context, _ uuid.UUID) (float64, int, error) {
	return 4.5, 10, nil
}

type mockCache struct{}

func (m *mockCache) Get(_ context.Context, _ string) (string, error)                    { return "", nil }
func (m *mockCache) Set(_ context.Context, _ string, _ interface{}, _ int) error         { return nil }
func (m *mockCache) Delete(_ context.Context, _ string) error                            { return nil }

type mockEvents struct {
	published bool
}

func (m *mockEvents) Publish(_ context.Context, _ string, _ interface{}) error {
	m.published = true
	return nil
}

// ── Pruebas del servicio ──────────────────────────────────────

func TestProductService_List_PaginacionPorDefecto(t *testing.T) {
	repo := &mockProductRepo{
		productos: []model.Product{
			{ID: uuid.New(), Name: "Producto A", Price: 100},
			{ID: uuid.New(), Name: "Producto B", Price: 200},
		},
		total: 2,
	}

	svc := NewProductService(repo, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, &mockEvents{})

	resultado, err := svc.List(context.Background(), model.ProductFilter{})
	if err != nil {
		t.Fatalf("List() error inesperado: %v", err)
	}
	if resultado.Page != 1 {
		t.Errorf("Page = %d, esperado 1 (valor por defecto)", resultado.Page)
	}
	if resultado.Limit != 20 {
		t.Errorf("Limit = %d, esperado 20 (valor por defecto)", resultado.Limit)
	}
	if resultado.Total != 2 {
		t.Errorf("Total = %d, esperado 2", resultado.Total)
	}
	if len(resultado.Products) != 2 {
		t.Errorf("len(Products) = %d, esperado 2", len(resultado.Products))
	}
}

func TestProductService_List_ErrorRepositorio(t *testing.T) {
	repo := &mockProductRepo{err: errors.New("conexión rechazada")}
	svc := NewProductService(repo, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, &mockEvents{})

	_, err := svc.List(context.Background(), model.ProductFilter{Page: 1, Limit: 10})
	if err == nil {
		t.Fatal("List() debería retornar error cuando el repositorio falla")
	}
}

func TestProductService_GetBySlug_Exito(t *testing.T) {
	productoID := uuid.New()
	repo := &mockProductRepo{
		productos: []model.Product{
			{ID: productoID, Slug: "camiseta-azul", Name: "Camiseta Azul"},
		},
	}

	svc := NewProductService(repo, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, &mockEvents{})

	producto, err := svc.GetBySlug(context.Background(), "camiseta-azul")
	if err != nil {
		t.Fatalf("GetBySlug() error inesperado: %v", err)
	}
	if producto.ID != productoID {
		t.Errorf("ID = %v, esperado %v", producto.ID, productoID)
	}
}

func TestProductService_GetBySlug_NoEncontrado(t *testing.T) {
	repo := &mockProductRepo{productos: []model.Product{}}
	svc := NewProductService(repo, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, &mockEvents{})

	_, err := svc.GetBySlug(context.Background(), "slug-inexistente")
	if err == nil {
		t.Fatal("GetBySlug() debería retornar error para slug inexistente")
	}
	if !errors.Is(err, ErrProductNotFound) {
		t.Errorf("error = %v, esperado ErrProductNotFound", err)
	}
}

func TestProductService_Create_Exito(t *testing.T) {
	repo := &mockProductRepo{}
	eventos := &mockEvents{}
	svc := NewProductService(repo, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, eventos)

	req := port.CreateProductRequest{
		SKU:         "CAM-001",
		Name:        "Camiseta Negra",
		Description: "Camiseta de algodón premium",
		Price:       299.99,
		Brand:       "CloudMart",
		Tags:        []string{"ropa", "camiseta"},
		Images:      []string{"https://img.ejemplo.com/cam-001.jpg"},
	}

	producto, err := svc.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create() error inesperado: %v", err)
	}
	if producto.Name != "Camiseta Negra" {
		t.Errorf("Name = %s, esperado 'Camiseta Negra'", producto.Name)
	}
	if producto.Slug != "camiseta-negra" {
		t.Errorf("Slug = %s, esperado 'camiseta-negra'", producto.Slug)
	}
	if producto.ThumbnailURL != "https://img.ejemplo.com/cam-001.jpg" {
		t.Errorf("ThumbnailURL no se asignó correctamente desde Images[0]")
	}
	if !producto.IsActive {
		t.Error("IsActive debería ser true para productos nuevos")
	}
	if !eventos.published {
		t.Error("Se esperaba que se publicara el evento product.created")
	}
	if repo.created == nil {
		t.Error("Se esperaba que el producto se guardara en el repositorio")
	}
}

func TestProductService_AddReview_CalificacionInvalida(t *testing.T) {
	svc := NewProductService(&mockProductRepo{}, &mockCategoryRepo{}, &mockReviewRepo{}, &mockCache{}, &mockEvents{})

	tests := []struct {
		nombre      string
		calificacion int
	}{
		{"calificación cero", 0},
		{"calificación negativa", -1},
		{"calificación mayor a 5", 6},
		{"calificación muy alta", 100},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			req := port.AddReviewRequest{
				ProductID: uuid.New(),
				UserID:    uuid.New(),
				Rating:    tt.calificacion,
				Title:     "Test",
				Body:      "Test body",
			}
			_, err := svc.AddReview(context.Background(), req)
			if !errors.Is(err, ErrInvalidRating) {
				t.Errorf("AddReview() con rating %d: error = %v, esperado ErrInvalidRating", tt.calificacion, err)
			}
		})
	}
}

func TestProductService_AddReview_CalificacionValida(t *testing.T) {
	reviewRepo := &mockReviewRepo{}
	svc := NewProductService(&mockProductRepo{}, &mockCategoryRepo{}, reviewRepo, &mockCache{}, &mockEvents{})

	for rating := 1; rating <= 5; rating++ {
		req := port.AddReviewRequest{
			ProductID: uuid.New(),
			UserID:    uuid.New(),
			Rating:    rating,
			Title:     "Excelente",
			Body:      "Muy buen producto",
		}
		resena, err := svc.AddReview(context.Background(), req)
		if err != nil {
			t.Errorf("AddReview() con rating %d: error inesperado %v", rating, err)
		}
		if resena.Rating != rating {
			t.Errorf("Rating = %d, esperado %d", resena.Rating, rating)
		}
	}
}

func TestGenerateSlug(t *testing.T) {
	tests := []struct {
		nombre   string
		entrada  string
		esperado string
	}{
		{"texto simple", "Camiseta Azul", "camiseta-azul"},
		{"con apóstrofe", "Men's Jacket", "mens-jacket"},
		{"con comillas", "Producto \"Premium\"", "producto-premium"},
		{"todo minúsculas", "already lowercase", "already-lowercase"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			resultado := generateSlug(tt.entrada)
			if resultado != tt.esperado {
				t.Errorf("generateSlug(%q) = %q, esperado %q", tt.entrada, resultado, tt.esperado)
			}
		})
	}
}
