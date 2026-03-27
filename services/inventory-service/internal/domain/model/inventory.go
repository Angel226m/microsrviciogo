package model

import (
	"time"

	"github.com/google/uuid"
)

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

// Available returns the quantity available for sale.
func (s *Stock) Available() int {
	return s.Quantity - s.Reserved
}

// NeedsReorder checks if stock is below reorder level.
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
