// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidades de Pedido
// Raíz agregada del contexto acotado de pedidos
// ═══════════════════════════════════════════════════════════════
package model

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type OrderStatus string

const (
	StatusPending    OrderStatus = "pending"
	StatusConfirmed  OrderStatus = "confirmed"
	StatusProcessing OrderStatus = "processing"
	StatusShipped    OrderStatus = "shipped"
	StatusDelivered  OrderStatus = "delivered"
	StatusCancelled  OrderStatus = "cancelled"
	StatusRefunded   OrderStatus = "refunded"
)

// Order es la raíz agregada del contexto acotado de pedidos.
type Order struct {
	ID              uuid.UUID   `json:"id"`
	OrderNumber     string      `json:"order_number"`
	UserID          uuid.UUID   `json:"user_id"`
	Status          OrderStatus `json:"status"`
	Subtotal        float64     `json:"subtotal"`
	Tax             float64     `json:"tax"`
	ShippingCost    float64     `json:"shipping_cost"`
	Discount        float64     `json:"discount"`
	Total           float64     `json:"total"`
	Currency        string      `json:"currency"`
	ShippingAddress Address     `json:"shipping_address"`
	BillingAddress  *Address    `json:"billing_address,omitempty"`
	Items           []OrderItem `json:"items"`
	Notes           string      `json:"notes,omitempty"`
	CreatedAt       time.Time   `json:"created_at"`
	UpdatedAt       time.Time   `json:"updated_at"`
}

// OrderItem representa un artículo individual dentro de un pedido.
type OrderItem struct {
	ID          uuid.UUID `json:"id"`
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	ProductSKU  string    `json:"product_sku"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	CreatedAt   time.Time `json:"created_at"`
}

type Address struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

// CanCancel verifica regla de negocio: solo pedidos pendientes/confirmados pueden cancelarse.
func (o *Order) CanCancel() bool {
	return o.Status == StatusPending || o.Status == StatusConfirmed
}

// CalculateTotal calcula el total del pedido: artículos + impuesto + envío - descuentos.
func (o *Order) CalculateTotal() {
	o.Subtotal = 0
	for _, item := range o.Items {
		o.Subtotal += item.TotalPrice
	}
	o.Tax = o.Subtotal * 0.16 // 16% IVA
	o.Total = o.Subtotal + o.Tax + o.ShippingCost - o.Discount
}

// GenerateOrderNumber crea un número de pedido único.
func GenerateOrderNumber() string {
	return fmt.Sprintf("CM-%d%04d", time.Now().Unix()%100000, rand.Intn(10000))
}
