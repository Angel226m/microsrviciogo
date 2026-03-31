// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Producto
// Valida descuentos, porcentajes y reglas de negocio
// ═══════════════════════════════════════════════════════════════
package model

import (
	"testing"
)

func TestProduct_HasDiscount(t *testing.T) {
	tests := []struct {
		nombre         string
		precio         float64
		precioAnterior *float64
		esperado       bool
	}{
		{"sin precio de comparación", 100.0, nil, false},
		{"con descuento", 80.0, punteroFloat(120.0), true},
		{"precio igual", 100.0, punteroFloat(100.0), false},
		{"precio mayor que comparación", 150.0, punteroFloat(100.0), false},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p := Product{Price: tt.precio, CompareAtPrice: tt.precioAnterior}
			if resultado := p.HasDiscount(); resultado != tt.esperado {
				t.Errorf("HasDiscount() = %v, esperado %v", resultado, tt.esperado)
			}
		})
	}
}

func TestProduct_DiscountPercentage(t *testing.T) {
	tests := []struct {
		nombre         string
		precio         float64
		precioAnterior *float64
		esperado       float64
	}{
		{"sin descuento", 100.0, nil, 0},
		{"50% de descuento", 50.0, punteroFloat(100.0), 50.0},
		{"25% de descuento", 75.0, punteroFloat(100.0), 25.0},
		{"precio igual sin descuento", 100.0, punteroFloat(100.0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			p := Product{Price: tt.precio, CompareAtPrice: tt.precioAnterior}
			resultado := p.DiscountPercentage()
			if resultado != tt.esperado {
				t.Errorf("DiscountPercentage() = %.2f, esperado %.2f", resultado, tt.esperado)
			}
		})
	}
}

func TestProductFilter_Defaults(t *testing.T) {
	// Verificar que un filtro vacío tiene valores cero
	filtro := ProductFilter{}
	if filtro.Page != 0 {
		t.Errorf("Page por defecto debería ser 0, obtuve %d", filtro.Page)
	}
	if filtro.Limit != 0 {
		t.Errorf("Limit por defecto debería ser 0, obtuve %d", filtro.Limit)
	}
}

func punteroFloat(v float64) *float64 {
	return &v
}
