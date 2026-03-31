// ═══════════════════════════════════════════════════════════════
// Puertos de Dominio – Límites hexagonales del servicio de productos
// Puertos primarios (entrada) y secundarios (salida)
// ═══════════════════════════════════════════════════════════════
package port

import (
	"context"

	"github.com/cloudmart/product-service/internal/domain/model"
	"github.com/google/uuid"
)

// ── Puertos Primarios ─────────────────────────────────────────────────────

type ProductService interface {
	List(ctx context.Context, filter model.ProductFilter) (*model.ProductListResult, error)
	GetBySlug(ctx context.Context, slug string) (*model.Product, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	Create(ctx context.Context, req CreateProductRequest) (*model.Product, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id uuid.UUID) error
	ListCategories(ctx context.Context) ([]model.Category, error)
	GetReviews(ctx context.Context, productID uuid.UUID) ([]model.Review, error)
	AddReview(ctx context.Context, req AddReviewRequest) (*model.Review, error)
}

// ── Puertos Secundarios ───────────────────────────────────────────────────

type ProductRepository interface {
	List(ctx context.Context, filter model.ProductFilter) ([]model.Product, int, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Product, error)
	FindBySlug(ctx context.Context, slug string) (*model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryRepository interface {
	List(ctx context.Context) ([]model.Category, error)
	FindByID(ctx context.Context, id uuid.UUID) (*model.Category, error)
}

type ReviewRepository interface {
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]model.Review, error)
	Create(ctx context.Context, review *model.Review) error
	GetAverageRating(ctx context.Context, productID uuid.UUID) (float64, int, error)
}

type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Delete(ctx context.Context, key string) error
}

type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// ── DTOs de Solicitud ──────────────────────────────────────────────────

type CreateProductRequest struct {
	SKU              string            `json:"sku"`
	Name             string            `json:"name"`
	Description      string            `json:"description"`
	ShortDescription string            `json:"short_description"`
	Price            float64           `json:"price"`
	CompareAtPrice   *float64          `json:"compare_at_price,omitempty"`
	CategoryID       *uuid.UUID        `json:"category_id,omitempty"`
	Brand            string            `json:"brand"`
	Tags             []string          `json:"tags"`
	Images           []string          `json:"images"`
	Attributes       map[string]string `json:"attributes,omitempty"`
}

type UpdateProductRequest struct {
	Name             string            `json:"name,omitempty"`
	Description      string            `json:"description,omitempty"`
	ShortDescription string            `json:"short_description,omitempty"`
	Price            *float64          `json:"price,omitempty"`
	CompareAtPrice   *float64          `json:"compare_at_price,omitempty"`
	Brand            string            `json:"brand,omitempty"`
	Tags             []string          `json:"tags,omitempty"`
	Images           []string          `json:"images,omitempty"`
	IsActive         *bool             `json:"is_active,omitempty"`
	IsFeatured       *bool             `json:"is_featured,omitempty"`
	Attributes       map[string]string `json:"attributes,omitempty"`
}

type AddReviewRequest struct {
	ProductID uuid.UUID `json:"product_id"`
	UserID    uuid.UUID `json:"user_id"`
	Rating    int       `json:"rating"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
}
