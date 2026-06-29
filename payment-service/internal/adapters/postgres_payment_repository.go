package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/payment-service/internal/domain"
	"GoCommerceX/payment-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresPaymentRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresPaymentRepository(pool *pgxpool.Pool) *PostgresPaymentRepository {
	return &PostgresPaymentRepository{pool: pool}
}

func (r *PostgresPaymentRepository) Save(ctx context.Context, payment domain.Payment) error {
	query := `
		INSERT INTO payments (id, order_id, amount, method, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.pool.Exec(ctx, query, payment.ID, payment.OrderID, payment.Amount, payment.Method, payment.Status, payment.CreatedAt, payment.UpdatedAt)
	return err
}

func (r *PostgresPaymentRepository) FindByID(ctx context.Context, id string) (domain.Payment, error) {
	query := `SELECT id, order_id, amount, method, status, created_at, updated_at FROM payments WHERE id = $1`
	var payment domain.Payment
	err := r.pool.QueryRow(ctx, query, id).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &payment.Method, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Payment{}, ports.ErrPaymentNotFound
		}
		return domain.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (domain.Payment, error) {
	query := `SELECT id, order_id, amount, method, status, created_at, updated_at FROM payments WHERE order_id = $1`
	var payment domain.Payment
	err := r.pool.QueryRow(ctx, query, orderID).Scan(&payment.ID, &payment.OrderID, &payment.Amount, &payment.Method, &payment.Status, &payment.CreatedAt, &payment.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Payment{}, ports.ErrPaymentNotFound
		}
		return domain.Payment{}, err
	}
	return payment, nil
}

func (r *PostgresPaymentRepository) Update(ctx context.Context, payment domain.Payment) error {
	query := `
		UPDATE payments
		SET status = $1, updated_at = $2
		WHERE id = $3
	`
	payment.UpdatedAt = time.Now()
	result, err := r.pool.Exec(ctx, query, payment.Status, payment.UpdatedAt, payment.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrPaymentNotFound
	}
	return nil
}
