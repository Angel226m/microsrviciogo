// ═══════════════════════════════════════════════════════════════
// Capa de Aplicación – Servicio de Usuario (casos de uso / lógica de negocio)
// Orquesta entidades de dominio, repositorios y eventos
// ═══════════════════════════════════════════════════════════════
package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/cloudmart/user-service/internal/domain/event"
	"github.com/cloudmart/user-service/internal/domain/model"
	"github.com/cloudmart/user-service/internal/domain/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found")
)

// userService implementa port.UserService con arquitectura hexagonal.
type userService struct {
	userRepo    port.UserRepository
	addrRepo    port.AddressRepository
	cache       port.CacheRepository
	events      port.EventPublisher
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewUserService crea un nuevo servicio de aplicación de usuario con todas las dependencias inyectadas.
func NewUserService(
	userRepo port.UserRepository,
	addrRepo port.AddressRepository,
	cache port.CacheRepository,
	events port.EventPublisher,
	jwtSecret string,
) port.UserService {
	return &userService{
		userRepo:    userRepo,
		addrRepo:    addrRepo,
		cache:       cache,
		events:      events,
		jwtSecret:   jwtSecret,
		tokenExpiry: 24 * time.Hour,
	}
}

func (s *userService) Register(ctx context.Context, req port.RegisterRequest) (*model.User, error) {
	// Check if email already exists
	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, ErrEmailTaken
	}

	// Hash password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Phone:        req.Phone,
		Role:         model.RoleCustomer,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// Publish domain event
	_ = s.events.Publish(ctx, event.SubjectUserRegistered, event.UserRegistered{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Timestamp: time.Now(),
	})

	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (*model.AuthTokens, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"role":    string(user.Role),
		"iat":     now.Unix(),
		"exp":     now.Add(s.tokenExpiry).Unix(),
		"iss":     "cloudmart",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, fmt.Errorf("sign token: %w", err)
	}

	// Update last login
	_ = s.userRepo.UpdateLastLogin(ctx, user.ID)

	// Publish login event
	_ = s.events.Publish(ctx, event.SubjectUserLoggedIn, event.UserLoggedIn{
		UserID:    user.ID,
		Email:     user.Email,
		Timestamp: now,
	})

	return &model.AuthTokens{
		AccessToken: tokenString,
		ExpiresIn:   int64(s.tokenExpiry.Seconds()),
		TokenType:   "Bearer",
	}, nil
}

func (s *userService) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("user:%s", id.String())
	_, _ = s.cache.Get(ctx, cacheKey) // simplified — in production, deserialize from cache

	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// Cache the result
	_ = s.cache.Set(ctx, cacheKey, user, 300)

	return user, nil
}

func (s *userService) Update(ctx context.Context, id uuid.UUID, req port.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("update user: %w", err)
	}

	// Invalidate cache
	cacheKey := fmt.Sprintf("user:%s", id.String())
	_ = s.cache.Delete(ctx, cacheKey)

	// Publish event
	_ = s.events.Publish(ctx, event.SubjectUserUpdated, event.UserUpdated{
		UserID:    user.ID,
		Timestamp: time.Now(),
	})

	return user, nil
}

func (s *userService) ListAddresses(ctx context.Context, userID uuid.UUID) ([]model.Address, error) {
	return s.addrRepo.FindByUserID(ctx, userID)
}

func (s *userService) AddAddress(ctx context.Context, userID uuid.UUID, req port.AddAddressRequest) (*model.Address, error) {
	addr := &model.Address{
		ID:        uuid.New(),
		UserID:    userID,
		Label:     req.Label,
		Street:    req.Street,
		City:      req.City,
		State:     req.State,
		ZipCode:   req.ZipCode,
		Country:   req.Country,
		CreatedAt: time.Now(),
	}

	if err := s.addrRepo.Create(ctx, addr); err != nil {
		return nil, fmt.Errorf("create address: %w", err)
	}

	return addr, nil
}
