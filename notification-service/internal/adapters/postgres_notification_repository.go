package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/notification-service/internal/domain"
	"GoCommerceX/notification-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresNotificationRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresNotificationRepository(pool *pgxpool.Pool) *PostgresNotificationRepository {
	return &PostgresNotificationRepository{pool: pool}
}

func (r *PostgresNotificationRepository) Save(ctx context.Context, notification domain.Notification) error {
	query := `
		INSERT INTO notifications (id, order_id, user_id, type, message, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query, notification.ID, notification.OrderID, notification.UserID, notification.Type, notification.Message, notification.Status, notification.CreatedAt, notification.UpdatedAt)
	return err
}

func (r *PostgresNotificationRepository) FindByID(ctx context.Context, id string) (domain.Notification, error) {
	query := `SELECT id, order_id, user_id, type, message, status, created_at, updated_at FROM notifications WHERE id = $1`
	var notification domain.Notification
	err := r.pool.QueryRow(ctx, query, id).Scan(&notification.ID, &notification.OrderID, &notification.UserID, &notification.Type, &notification.Message, &notification.Status, &notification.CreatedAt, &notification.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, ports.ErrNotificationNotFound
		}
		return domain.Notification{}, err
	}
	return notification, nil
}

func (r *PostgresNotificationRepository) FindByOrderID(ctx context.Context, orderID string) (domain.Notification, error) {
	query := `SELECT id, order_id, user_id, type, message, status, created_at, updated_at FROM notifications WHERE order_id = $1`
	var notification domain.Notification
	err := r.pool.QueryRow(ctx, query, orderID).Scan(&notification.ID, &notification.OrderID, &notification.UserID, &notification.Type, &notification.Message, &notification.Status, &notification.CreatedAt, &notification.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Notification{}, ports.ErrNotificationNotFound
		}
		return domain.Notification{}, err
	}
	return notification, nil
}

func (r *PostgresNotificationRepository) touch(notification domain.Notification) domain.Notification {
	notification.UpdatedAt = time.Now()
	return notification
}
