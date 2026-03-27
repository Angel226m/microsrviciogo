// ═══════════════════════════════════════════════════════════════
// PostgreSQL Repository – User driven adapter (infrastructure)
// ═══════════════════════════════════════════════════════════════
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudmart/user-service/internal/domain/model"
	"github.com/cloudmart/user-service/internal/domain/port"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userPostgresRepo struct {
	pool *pgxpool.Pool
}

// NewUserPostgresRepo creates a new PostgreSQL user repository adapter.
func NewUserPostgresRepo(pool *pgxpool.Pool) port.UserRepository {
	return &userPostgresRepo{pool: pool}
}

func (r *userPostgresRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users.accounts (id, email, password_hash, first_name, last_name, phone, role, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.PasswordHash, user.FirstName, user.LastName,
		user.Phone, string(user.Role), user.IsActive, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}
	return nil
}

func (r *userPostgresRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, avatar_url,
		       role, is_active, email_verified, last_login_at, created_at, updated_at
		FROM users.accounts WHERE id = $1`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.AvatarURL, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return user, nil
}

func (r *userPostgresRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, email, password_hash, first_name, last_name, phone, avatar_url,
		       role, is_active, email_verified, last_login_at, created_at, updated_at
		FROM users.accounts WHERE email = $1`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FirstName, &user.LastName,
		&user.Phone, &user.AvatarURL, &user.Role, &user.IsActive, &user.EmailVerified,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return user, nil
}

func (r *userPostgresRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users.accounts
		SET first_name = $2, last_name = $3, phone = $4, avatar_url = $5, updated_at = $6
		WHERE id = $1`

	_, err := r.pool.Exec(ctx, query,
		user.ID, user.FirstName, user.LastName, user.Phone, user.AvatarURL, user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (r *userPostgresRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE users.accounts SET last_login_at = $2 WHERE id = $1`
	_, err := r.pool.Exec(ctx, query, id, time.Now())
	return err
}

// ── Address Repository ────────────────────────────────────────────────

type addressPostgresRepo struct {
	pool *pgxpool.Pool
}

func NewAddressPostgresRepo(pool *pgxpool.Pool) port.AddressRepository {
	return &addressPostgresRepo{pool: pool}
}

func (r *addressPostgresRepo) Create(ctx context.Context, addr *model.Address) error {
	query := `
		INSERT INTO users.addresses (id, user_id, label, street, city, state, zip_code, country, is_default, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, query,
		addr.ID, addr.UserID, addr.Label, addr.Street, addr.City,
		addr.State, addr.ZipCode, addr.Country, addr.IsDefault, addr.CreatedAt,
	)
	return err
}

func (r *addressPostgresRepo) FindByUserID(ctx context.Context, userID uuid.UUID) ([]model.Address, error) {
	query := `SELECT id, user_id, label, street, city, state, zip_code, country, is_default, created_at
		FROM users.addresses WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var a model.Address
		if err := rows.Scan(&a.ID, &a.UserID, &a.Label, &a.Street, &a.City,
			&a.State, &a.ZipCode, &a.Country, &a.IsDefault, &a.CreatedAt); err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (r *addressPostgresRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, "DELETE FROM users.addresses WHERE id = $1", id)
	return err
}
