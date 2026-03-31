// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Inventario
// Valida disponibilidad, nivel de reorden y movimientos
// ═══════════════════════════════════════════════════════════════
package model

import (
	"testing"
)

func TestStock_Available(t *testing.T) {
	tests := []struct {
		nombre    string
		cantidad  int
		reservado int
		esperado  int
	}{
		{"sin reservas", 100, 0, 100},
		{"con reservas parciales", 100, 30, 70},
		{"todo reservado", 50, 50, 0},
		{"sobre-reservado (negativo)", 10, 20, -10},
		{"stock vacío", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			stock := &Stock{Quantity: tt.cantidad, Reserved: tt.reservado}
			if resultado := stock.Available(); resultado != tt.esperado {
				t.Errorf("Available() = %d, esperado %d", resultado, tt.esperado)
			}
		})
	}
}

func TestStock_NeedsReorder(t *testing.T) {
	tests := []struct {
		nombre       string
		cantidad     int
		reservado    int
		nivelReorden int
		esperado     bool
	}{
		{"stock suficiente", 100, 0, 10, false},
		{"disponible justo en el nivel", 10, 0, 10, true},
		{"disponible bajo el nivel", 15, 10, 10, true},
		{"sin stock necesita reorden", 0, 0, 5, true},
		{"nivel de reorden cero con stock", 10, 0, 0, false},
		{"nivel de reorden cero sin stock", 0, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			stock := &Stock{
				Quantity:     tt.cantidad,
				Reserved:     tt.reservado,
				ReorderLevel: tt.nivelReorden,
			}
			if resultado := stock.NeedsReorder(); resultado != tt.esperado {
				t.Errorf("NeedsReorder() = %v, esperado %v (disponible: %d, nivel: %d)",
					resultado, tt.esperado, stock.Available(), tt.nivelReorden)
			}
		})
	}
}

func TestMovementType_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		tipo     MovementType
		esperado string
	}{
		{"entrada", Inbound, "inbound"},
		{"salida", Outbound, "outbound"},
		{"reserva", Reservation, "reservation"},
		{"liberación", Release, "release"},
		{"ajuste", Adjustment, "adjustment"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.tipo) != tt.esperado {
				t.Errorf("Tipo %s = %q, esperado %q", tt.nombre, tt.tipo, tt.esperado)
			}
		})
	}
}
