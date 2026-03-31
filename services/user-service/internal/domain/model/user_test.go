// ═══════════════════════════════════════════════════════════════
// Pruebas Unitarias – Modelo de Dominio de Usuario
// Valida reglas de negocio y comportamiento de entidades
// ═══════════════════════════════════════════════════════════════
package model

import "testing"

func TestUser_FullName(t *testing.T) {
	tests := []struct {
		nombre   string
		usuario  User
		esperado string
	}{
		{
			nombre:   "nombre completo con ambos campos",
			usuario:  User{FirstName: "Juan", LastName: "Pérez"},
			esperado: "Juan Pérez",
		},
		{
			nombre:   "solo nombre",
			usuario:  User{FirstName: "María", LastName: ""},
			esperado: "María ",
		},
		{
			nombre:   "campos vacíos",
			usuario:  User{},
			esperado: " ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			resultado := tt.usuario.FullName()
			if resultado != tt.esperado {
				t.Errorf("FullName() = %q, esperado %q", resultado, tt.esperado)
			}
		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tests := []struct {
		nombre   string
		rol      Role
		esperado bool
	}{
		{"rol admin es admin", RoleAdmin, true},
		{"rol customer no es admin", RoleCustomer, false},
		{"rol seller no es admin", RoleSeller, false},
	}

	for _, tt := range tests {
		t.Run(tt.nombre, func(t *testing.T) {
			u := User{Role: tt.rol}
			if resultado := u.IsAdmin(); resultado != tt.esperado {
				t.Errorf("IsAdmin() = %v, esperado %v", resultado, tt.esperado)
			}
		})
	}
}

func TestRole_Constantes(t *testing.T) {
	// Verificar que las constantes de rol tienen los valores esperados
	if RoleCustomer != "customer" {
		t.Errorf("RoleCustomer = %q, esperado %q", RoleCustomer, "customer")
	}
	if RoleAdmin != "admin" {
		t.Errorf("RoleAdmin = %q, esperado %q", RoleAdmin, "admin")
	}
	if RoleSeller != "seller" {
		t.Errorf("RoleSeller = %q, esperado %q", RoleSeller, "seller")
	}
}
