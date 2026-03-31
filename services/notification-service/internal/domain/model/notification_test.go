// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Notificación
// Valida tipos, estados y construcción de entidades
// ═══════════════════════════════════════════════════════════════
package model

import (
	"testing"

	"github.com/google/uuid"
)

func TestNotificationType_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		tipo     NotificationType
		esperado string
	}{
		{"correo electrónico", NotificationEmail, "email"},
		{"SMS", NotificationSMS, "sms"},
		{"push", NotificationPush, "push"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.tipo) != tt.esperado {
				t.Errorf("Tipo %s = %q, esperado %q", tt.nombre, tt.tipo, tt.esperado)
			}
		})
	}
}

func TestNotificationStatus_Constantes(t *testing.T) {
	tests := []struct {
		nombre   string
		estado   NotificationStatus
		esperado string
	}{
		{"pendiente", StatusPending, "pending"},
		{"enviado", StatusSent, "sent"},
		{"fallido", StatusFailed, "failed"},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			if string(tt.estado) != tt.esperado {
				t.Errorf("Estado %s = %q, esperado %q", tt.nombre, tt.estado, tt.esperado)
			}
		})
	}
}

func TestNotification_Construccion(t *testing.T) {
	notif := Notification{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Type:    NotificationEmail,
		Channel: "usuario@ejemplo.com",
		Subject: "Confirmación de pedido",
		Body:    "Su pedido #CM-12345 ha sido confirmado.",
		Status:  StatusPending,
	}

	if notif.Type != NotificationEmail {
		t.Errorf("Type = %s, esperado email", notif.Type)
	}
	if notif.Status != StatusPending {
		t.Errorf("Status = %s, esperado pending", notif.Status)
	}
	if notif.RetryCount != 0 {
		t.Errorf("RetryCount inicial = %d, esperado 0", notif.RetryCount)
	}
	if notif.SentAt != nil {
		t.Error("SentAt debería ser nil para notificación pendiente")
	}
}

func TestTemplate_Construccion(t *testing.T) {
	plantilla := Template{
		ID:      uuid.New(),
		Name:    "bienvenida",
		Subject: "Bienvenido a CloudMart",
		Body:    "Hola {{.Nombre}}, bienvenido a nuestra plataforma.",
	}

	if plantilla.Name != "bienvenida" {
		t.Errorf("Name = %s, esperado bienvenida", plantilla.Name)
	}
	if plantilla.Subject != "Bienvenido a CloudMart" {
		t.Errorf("Subject = %s, esperado 'Bienvenido a CloudMart'", plantilla.Subject)
	}
}
