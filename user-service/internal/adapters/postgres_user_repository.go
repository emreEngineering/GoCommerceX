package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/user-service/internal/domain"
	"GoCommerceX/user-service/internal/ports"

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
		INSERT INTO user_profiles (id, email, first_name, last_name, phone, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query,
		user.ID, user.Email, user.FirstName, user.LastName, user.Phone,
		user.CreatedAt, user.UpdatedAt,
	)
	return err
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	query := `SELECT id, email, first_name, last_name, phone, created_at, updated_at FROM user_profiles WHERE id = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	query := `SELECT id, email, first_name, last_name, phone, created_at, updated_at FROM user_profiles WHERE email = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.Phone,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, ports.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user domain.User) error {
	query := `
		UPDATE user_profiles
		SET first_name = $1, last_name = $2, phone = $3, updated_at = $4
		WHERE id = $5
	`
	user.UpdatedAt = time.Now()
	result, err := r.pool.Exec(ctx, query,
		user.FirstName, user.LastName, user.Phone, user.UpdatedAt, user.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM user_profiles WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrUserNotFound
	}
	return nil
}
