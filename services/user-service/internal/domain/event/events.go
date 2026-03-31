// ═══════════════════════════════════════════════════════════════
// Eventos de Dominio – Eventos emitidos por el contexto acotado de Usuario
// Comunicación asíncrona con otros microservicios vía NATS
// ═══════════════════════════════════════════════════════════════
package event

import (
	"time"

	"github.com/google/uuid"
)

const (
	SubjectUserRegistered = "user.registered"
	SubjectUserUpdated    = "user.updated"
	SubjectUserLoggedIn   = "user.logged_in"
)

type UserRegistered struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Timestamp time.Time `json:"timestamp"`
}

type UserUpdated struct {
	UserID    uuid.UUID `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
}

type UserLoggedIn struct {
	UserID    uuid.UUID `json:"user_id"`
	Email     string    `json:"email"`
	Timestamp time.Time `json:"timestamp"`
}
