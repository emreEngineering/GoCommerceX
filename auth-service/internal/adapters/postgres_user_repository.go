package adapters

import (
	"context"
	"errors"

	"GoCommerceX/auth-service/internal/domain"
	"GoCommerceX/auth-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) Save(ctx context.Context, user domain.User) error {
	query := `
		INSERT INTO "users" (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `
		SELECT id, email, password_hash, created_at, updated_at
		FROM "users"
		WHERE email = $1
	`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM "users"
		WHERE id = $1
	`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrUserNotFound
	}
	return nil
}
