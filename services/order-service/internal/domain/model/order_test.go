// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Pedido
// Valida cancelación, cálculo de totales y generación de número
// ═══════════════════════════════════════════════════════════════
package model

import (
	"strings"
	"testing"
)

func TestOrder_CanCancel(t *testing.T) {
	tests := []struct {
		nombre   string
		estado   OrderStatus
		esperado bool
	}{
		{"pendiente puede cancelarse", StatusPending, true},
		{"confirmado puede cancelarse", StatusConfirmed, true},
		{"en proceso no puede cancelarse", StatusProcessing, false},
		{"enviado no puede cancelarse", StatusShipped, false},
		{"entregado no puede cancelarse", StatusDelivered, false},
		{"ya cancelado no puede re-cancelarse", StatusCancelled, false},
		{"reembolsado no puede cancelarse", StatusRefunded, false},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			pedido := &Order{Status: tt.estado}
			if resultado := pedido.CanCancel(); resultado != tt.esperado {
				t.Errorf("CanCancel() con estado %q = %v, esperado %v", tt.estado, resultado, tt.esperado)
			}
		})
	}
}

func TestOrder_CalculateTotal(t *testing.T) {
	tests := []struct {
		nombre      string
		articulos   []OrderItem
		envio       float64
		descuento   float64
		subEsperado float64
		ivaEsperado float64
		totalMinimo float64
		totalMaximo float64
	}{
		{
			nombre: "un artículo sin envío ni descuento",
			articulos: []OrderItem{
				{Quantity: 2, UnitPrice: 50.0, TotalPrice: 100.0},
			},
			envio:       0,
			descuento:   0,
			subEsperado: 100.0,
			ivaEsperado: 16.0,
			totalMinimo: 115.99,
			totalMaximo: 116.01,
		},
		{
			nombre: "múltiples artículos con envío y descuento",
			articulos: []OrderItem{
				{Quantity: 1, UnitPrice: 200.0, TotalPrice: 200.0},
				{Quantity: 3, UnitPrice: 50.0, TotalPrice: 150.0},
			},
			envio:       25.0,
			descuento:   10.0,
			subEsperado: 350.0,
			ivaEsperado: 56.0,
			totalMinimo: 420.99,
			totalMaximo: 421.01,
		},
		{
			nombre:      "pedido vacío",
			articulos:   []OrderItem{},
			envio:       0,
			descuento:   0,
			subEsperado: 0,
			ivaEsperado: 0,
			totalMinimo: -0.01,
			totalMaximo: 0.01,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			pedido := &Order{
				Items:        tt.articulos,
				ShippingCost: tt.envio,
				Discount:     tt.descuento,
			}
			pedido.CalculateTotal()

			if pedido.Subtotal != tt.subEsperado {
				t.Errorf("Subtotal = %.2f, esperado %.2f", pedido.Subtotal, tt.subEsperado)
			}
			if pedido.Tax != tt.ivaEsperado {
				t.Errorf("Tax (IVA 16%%) = %.2f, esperado %.2f", pedido.Tax, tt.ivaEsperado)
			}
			if pedido.Total < tt.totalMinimo || pedido.Total > tt.totalMaximo {
				t.Errorf("Total = %.2f, esperado entre %.2f y %.2f", pedido.Total, tt.totalMinimo, tt.totalMaximo)
			}
		})
	}
}

func TestGenerateOrderNumber(t *testing.T) {
	numero := GenerateOrderNumber()

	if !strings.HasPrefix(numero, "CM-") {
		t.Errorf("Número de pedido debería iniciar con 'CM-', obtuve: %s", numero)
	}

	// Verificar unicidad básica generando varios
	numeros := make(map[string]bool)
	for i := 0; i < 100; i++ {
		n := GenerateOrderNumber()
		numeros[n] = true
	}
	// Al menos 50 de 100 deberían ser únicos (componente aleatorio)
	if len(numeros) < 50 {
		t.Errorf("Poca unicidad: solo %d números únicos de 100 generados", len(numeros))
	}
}

func TestOrderStatus_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		estado   OrderStatus
		esperado string
	}{
		{"pendiente", StatusPending, "pending"},
		{"confirmado", StatusConfirmed, "confirmed"},
		{"procesando", StatusProcessing, "processing"},
		{"enviado", StatusShipped, "shipped"},
		{"entregado", StatusDelivered, "delivered"},
		{"cancelado", StatusCancelled, "cancelled"},
		{"reembolsado", StatusRefunded, "refunded"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.estado) != tt.esperado {
				t.Errorf("Estado %s = %q, esperado %q", tt.nombre, tt.estado, tt.esperado)
			}
		})
	}
}
