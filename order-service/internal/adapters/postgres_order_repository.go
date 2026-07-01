package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/order-service/internal/domain"
	"GoCommerceX/order-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresOrderRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{pool: pool}
}

func (r *PostgresOrderRepository) Save(ctx context.Context, order domain.Order) error {
	query := `
		INSERT INTO orders (id, user_id, payment_id, status, total_amount, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	var paymentID any
	if order.PaymentID != "" {
		paymentID = order.PaymentID
	}
	_, err := r.pool.Exec(ctx, query, order.ID, order.UserID, paymentID, order.Status, order.TotalAmount, order.CreatedAt, order.UpdatedAt)
	return err
}

func (r *PostgresOrderRepository) FindByID(ctx context.Context, id string) (domain.Order, error) {
	query := `SELECT id, user_id, payment_id, status, total_amount, created_at, updated_at FROM orders WHERE id = $1`
	var order domain.Order
	err := r.pool.QueryRow(ctx, query, id).Scan(&order.ID, &order.UserID, &order.PaymentID, &order.Status, &order.TotalAmount, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Order{}, ports.ErrOrderNotFound
		}
		return domain.Order{}, err
	}
	return order, nil
}

func (r *PostgresOrderRepository) Update(ctx context.Context, order domain.Order) error {
	query := `
		UPDATE orders
		SET payment_id = $1, status = $2, total_amount = $3, updated_at = $4
		WHERE id = $5
	`
	order.UpdatedAt = time.Now()
	var paymentID any
	if order.PaymentID != "" {
		paymentID = order.PaymentID
	}
	result, err := r.pool.Exec(ctx, query, paymentID, order.Status, order.TotalAmount, order.UpdatedAt, order.ID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrOrderNotFound
	}
	return nil
}

func (r *PostgresOrderRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM orders WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrOrderNotFound
	}
	return nil
}
