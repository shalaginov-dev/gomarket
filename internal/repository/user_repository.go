package repository

import (
	"context"
	"fmt"
	"gomarket/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, email string, passwordHash string) (*domain.User, error) {
	query := `
		INSERT INTO users (email, password_hash, created_at)
		VALUES ($1, $2, NOW())
		RETURNING user_id, email, role, created_at`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, email, passwordHash).
		Scan(&user.UserID, &user.Email, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
        SELECT user_id, email, password_hash, role, created_at
        FROM users
        WHERE email = $1
        LIMIT 1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(&user.UserID, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	return &user, nil
}
