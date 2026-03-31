// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Pago
// Valida constantes de estado y métodos de pago
// ═══════════════════════════════════════════════════════════════
package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestPaymentStatus_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		estado   PaymentStatus
		esperado string
	}{
		{"pendiente", PaymentPending, "pending"},
		{"procesando", PaymentProcessing, "processing"},
		{"completado", PaymentCompleted, "completed"},
		{"fallido", PaymentFailed, "failed"},
		{"reembolsado", PaymentRefunded, "refunded"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.estado) != tt.esperado {
				t.Errorf("Estado %s = %q, esperado %q", tt.nombre, tt.estado, tt.esperado)
			}
		})
	}
}

func TestPaymentMethod_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		metodo   PaymentMethod
		esperado string
	}{
		{"tarjeta de crédito", MethodCreditCard, "credit_card"},
		{"tarjeta de débito", MethodDebitCard, "debit_card"},
		{"PayPal", MethodPayPal, "paypal"},
		{"transferencia bancaria", MethodBankTransfer, "bank_transfer"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.metodo) != tt.esperado {
				t.Errorf("Método %s = %q, esperado %q", tt.nombre, tt.metodo, tt.esperado)
			}
		})
	}
}

func TestTransaction_CamposRequeridos(t *testing.T) {
	// Verifica que una transacción puede construirse con campos mínimos
	tx := Transaction{
		ID:       uuid.New(),
		OrderID:  uuid.New(),
		UserID:   uuid.New(),
		Amount:   99.99,
		Currency: "MXN",
		Method:   MethodCreditCard,
		Status:   PaymentPending,
		Provider: "stripe",
	}

	if tx.Amount != 99.99 {
		t.Errorf("Amount = %.2f, esperado 99.99", tx.Amount)
	}
	if tx.Currency != "MXN" {
		t.Errorf("Currency = %s, esperado MXN", tx.Currency)
	}
	if tx.Status != PaymentPending {
		t.Errorf("Status = %s, esperado pending", tx.Status)
	}
	if tx.Method != MethodCreditCard {
		t.Errorf("Method = %s, esperado credit_card", tx.Method)
	}
}
