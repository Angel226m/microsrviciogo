// ═══════════════════════════════════════════════════════════════
// Puertos de Dominio – Interfaces que definen los límites hexagonales
// Puertos primarios (entrada) y puertos secundarios (salida)
// ═══════════════════════════════════════════════════════════════
package port

import (
	"context"

	"github.com/cloudmart/user-service/internal/domain/model"
	"github.com/google/uuid"
)

// ── Puertos Primarios (llamados por adaptadores como handlers HTTP) ──

// UserService define los casos de uso para gestión de usuarios.
type UserService interface {
	Register(ctx context.Context, req RegisterRequest) (*model.User, error)
	Login(ctx context.Context, email, password string) (*model.AuthTokens, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*model.User, error)
	ListAddresses(ctx context.Context, userID uuid.UUID) ([]model.Address, error)
	AddAddress(ctx context.Context, userID uuid.UUID, req AddAddressRequest) (*model.Address, error)
}

// ── Puertos Secundarios (implementados por adaptadores de infraestructura) ──

// UserRepository define el contrato de persistencia para usuarios.
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id uuid.UUID) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	UpdateLastLogin(ctx context.Context, id uuid.UUID) error
}

// AddressRepository define el contrato de persistencia para direcciones.
type AddressRepository interface {
	Create(ctx context.Context, addr *model.Address) error
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Address, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// CacheRepository define el contrato de almacenamiento en caché.
type CacheRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, ttlSeconds int) error
	Delete(ctx context.Context, key string) error
}

// EventPublisher define el contrato para publicar eventos de dominio.
type EventPublisher interface {
	Publish(ctx context.Context, subject string, data interface{}) error
}

// ── DTOs de Solicitud ────────────────────────────────────────────────────────

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

type UpdateUserRequest struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Phone     string `json:"phone,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

type AddAddressRequest struct {
	Label   string `json:"label"`
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}
