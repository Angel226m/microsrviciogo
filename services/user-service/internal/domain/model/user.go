// ═══════════════════════════════════════════════════════════════
// Modelo de Dominio – Entidad Usuario (núcleo de negocio, cero dependencias)
// Arquitectura Hexagonal: capa de dominio pura sin dependencias externas
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// Role representa los niveles de autorización del usuario.
type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
	RoleSeller   Role = "seller"
)

// User es la entidad principal de dominio que representa un usuario de la plataforma.
type User struct {
	ID            uuid.UUID  `json:"id"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Phone         string     `json:"phone,omitempty"`
	AvatarURL     string     `json:"avatar_url,omitempty"`
	Role          Role       `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// FullName devuelve el nombre completo del usuario.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin verifica si el usuario tiene privilegios de administrador.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Address representa una dirección de envío/facturación del usuario.
type Address struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Label     string    `json:"label"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Country   string    `json:"country"`
	IsDefault bool      `json:"is_default"`
	CreatedAt time.Time `json:"created_at"`
}

// AuthTokens contiene los tokens JWT de acceso y actualización.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
