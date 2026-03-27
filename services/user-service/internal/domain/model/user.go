// ═══════════════════════════════════════════════════════════════
// Domain Model – User entity (business core, zero dependencies)
// ═══════════════════════════════════════════════════════════════
package model

import (
	"time"

	"github.com/google/uuid"
)

// Role represents user authorization levels
type Role string

const (
	RoleCustomer Role = "customer"
	RoleAdmin    Role = "admin"
	RoleSeller   Role = "seller"
)

// User is the core domain entity representing a platform user.
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

// FullName returns the user's full name.
func (u *User) FullName() string {
	return u.FirstName + " " + u.LastName
}

// IsAdmin checks if the user has admin privileges.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// Address represents a shipping/billing address for a user.
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

// AuthTokens holds the JWT access and refresh tokens.
type AuthTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}
