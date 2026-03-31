// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidades de Producto (cero dependencias externas)
// Arquitectura Hexagonal: capa de dominio pura
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// Product representa un artículo vendible en el catálogo.
type Product struct {
	ID               uuid.UUID          `json:"id"`
	SKU              string             `json:"sku"`
	Name             string             `json:"name"`
	Slug             string             `json:"slug"`
	Description      string             `json:"description"`
	ShortDescription string             `json:"short_description"`
	Price            float64            `json:"price"`
	CompareAtPrice   *float64           `json:"compare_at_price,omitempty"`
	Cost             *float64           `json:"cost,omitempty"`
	CategoryID       *uuid.UUID         `json:"category_id,omitempty"`
	Brand            string             `json:"brand"`
	Tags             []string           `json:"tags"`
	Images           []string           `json:"images"`
	ThumbnailURL     string             `json:"thumbnail_url"`
	Weight           *float64           `json:"weight,omitempty"`
	Dimensions       map[string]float64 `json:"dimensions,omitempty"`
	Attributes       map[string]string  `json:"attributes,omitempty"`
	IsActive         bool               `json:"is_active"`
	IsFeatured       bool               `json:"is_featured"`
	RatingAvg        float64            `json:"rating_avg"`
	RatingCount      int                `json:"rating_count"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

// HasDiscount verifica si este producto tiene un precio de comparación mayor al actual.
func (p *Product) HasDiscount() bool {
	return p.CompareAtPrice != nil && *p.CompareAtPrice > p.Price
}

// DiscountPercentage devuelve el porcentaje de descuento si aplica.
func (p *Product) DiscountPercentage() float64 {
	if !p.HasDiscount() {
		return 0
	}
	return ((*p.CompareAtPrice - p.Price) / *p.CompareAtPrice) * 100
}

// Category representa una agrupación de productos.
type Category struct {
	ID          uuid.UUID  `json:"id"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	ImageURL    string     `json:"image_url,omitempty"`
	ParentID    *uuid.UUID `json:"parent_id,omitempty"`
	SortOrder   int        `json:"sort_order"`
	IsActive    bool       `json:"is_active"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Review representa una reseña de producto hecha por un cliente.
type Review struct {
	ID         uuid.UUID `json:"id"`
	ProductID  uuid.UUID `json:"product_id"`
	UserID     uuid.UUID `json:"user_id"`
	Rating     int       `json:"rating"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}

// ProductFilter contiene los parámetros de consulta para listar productos.
type ProductFilter struct {
	CategoryID *uuid.UUID
	Search     string
	MinPrice   *float64
	MaxPrice   *float64
	Brand      string
	Tags       []string
	IsFeatured *bool
	SortBy     string
	SortDir    string
	Page       int
	Limit      int
}

// ProductListResult contiene los resultados paginados de productos.
type ProductListResult struct {
	Products   []Product `json:"products"`
	Total      int       `json:"total"`
	Page       int       `json:"page"`
	Limit      int       `json:"limit"`
	TotalPages int       `json:"total_pages"`
}
