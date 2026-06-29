package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/inventory-service/internal/domain"
	"GoCommerceX/inventory-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresInventoryRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresInventoryRepository(pool *pgxpool.Pool) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{pool: pool}
}

func (r *PostgresInventoryRepository) Save(ctx context.Context, inventory domain.Inventory) error {
	query := `
		INSERT INTO inventories (id, product_id, available_quantity, reserved_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.pool.Exec(ctx, query,
		inventory.ID, inventory.ProductID, inventory.AvailableQuantity, inventory.ReservedQuantity,
		inventory.CreatedAt, inventory.UpdatedAt,
	)
	return err
}

func (r *PostgresInventoryRepository) FindByID(ctx context.Context, id string) (domain.Inventory, error) {
	query := `SELECT id, product_id, available_quantity, reserved_quantity, created_at, updated_at FROM inventories WHERE id = $1`
	var inventory domain.Inventory
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&inventory.ID, &inventory.ProductID, &inventory.AvailableQuantity, &inventory.ReservedQuantity,
		&inventory.CreatedAt, &inventory.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Inventory{}, ports.ErrInventoryNotFound
		}
		return domain.Inventory{}, err
	}
	return inventory, nil
}

func (r *PostgresInventoryRepository) FindByProductID(ctx context.Context, productID string) (domain.Inventory, error) {
	query := `SELECT id, product_id, available_quantity, reserved_quantity, created_at, updated_at FROM inventories WHERE product_id = $1`
	var inventory domain.Inventory
	err := r.pool.QueryRow(ctx, query, productID).Scan(
		&inventory.ID, &inventory.ProductID, &inventory.AvailableQuantity, &inventory.ReservedQuantity,
		&inventory.CreatedAt, &inventory.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Inventory{}, ports.ErrInventoryNotFound
		}
		return domain.Inventory{}, err
	}
	return inventory, nil
}

func (r *PostgresInventoryRepository) Update(ctx context.Context, inventory domain.Inventory) error {
	query := `
		UPDATE inventories
		SET available_quantity = $1, reserved_quantity = $2, updated_at = $3
		WHERE id = $4
	`
	inventory.UpdatedAt = time.Now()
	result, err := r.pool.Exec(ctx, query,
		inventory.AvailableQuantity, inventory.ReservedQuantity, inventory.UpdatedAt, inventory.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrInventoryNotFound
	}
	return nil
}

func (r *PostgresInventoryRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM inventories WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrInventoryNotFound
	}
	return nil
}
