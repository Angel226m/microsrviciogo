// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidades de Inventario
// Stock, movimientos y tipos de operación de almacén
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// MovementType representa el tipo de movimiento de inventario.
type MovementType string

const (
	Inbound     MovementType = "inbound"
	Outbound    MovementType = "outbound"
	Reservation MovementType = "reservation"
	Release     MovementType = "release"
	Adjustment  MovementType = "adjustment"
)

type Stock struct {
	ID              uuid.UUID  `json:"id"`
	ProductID       uuid.UUID  `json:"product_id"`
	SKU             string     `json:"sku"`
	Quantity        int        `json:"quantity"`
	Reserved        int        `json:"reserved"`
	Warehouse       string     `json:"warehouse"`
	ReorderLevel    int        `json:"reorder_level"`
	ReorderQuantity int        `json:"reorder_quantity"`
	LastRestockedAt *time.Time `json:"last_restocked_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Available devuelve la cantidad disponible para venta.
func (s *Stock) Available() int {
	return s.Quantity - s.Reserved
}

// NeedsReorder verifica si el stock está por debajo del nivel de reorden.
func (s *Stock) NeedsReorder() bool {
	return s.Available() <= s.ReorderLevel
}

type Movement struct {
	ID            uuid.UUID    `json:"id"`
	ProductID     uuid.UUID    `json:"product_id"`
	Type          MovementType `json:"type"`
	Quantity      int          `json:"quantity"`
	ReferenceID   *uuid.UUID   `json:"reference_id,omitempty"`
	ReferenceType string       `json:"reference_type,omitempty"`
	Notes         string       `json:"notes,omitempty"`
	CreatedAt     time.Time    `json:"created_at"`
}
