package adapters

import (
	"context"
	"errors"
	"time"

	"GoCommerceX/product-service/internal/domain"
	"GoCommerceX/product-service/internal/ports"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresProductRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresProductRepository(pool *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{pool: pool}
}

func (r *PostgresProductRepository) Save(ctx context.Context, product domain.Product) error {
	query := `
		INSERT INTO products (id, sku, name, description, price, stock_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.pool.Exec(ctx, query,
		product.ID, product.SKU, product.Name, product.Description, product.Price, product.Stock,
		product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id string) (domain.Product, error) {
	query := `SELECT id, sku, name, description, price, stock_quantity, created_at, updated_at FROM products WHERE id = $1`
	var product domain.Product
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Stock,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, ports.ErrProductNotFound
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (r *PostgresProductRepository) FindBySKU(ctx context.Context, sku string) (domain.Product, error) {
	query := `SELECT id, sku, name, description, price, stock_quantity, created_at, updated_at FROM products WHERE sku = $1`
	var product domain.Product
	err := r.pool.QueryRow(ctx, query, sku).Scan(
		&product.ID, &product.SKU, &product.Name, &product.Description, &product.Price, &product.Stock,
		&product.CreatedAt, &product.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, ports.ErrProductNotFound
		}
		return domain.Product{}, err
	}
	return product, nil
}

func (r *PostgresProductRepository) Update(ctx context.Context, product domain.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock_quantity = $4, updated_at = $5
		WHERE id = $6
	`
	product.UpdatedAt = time.Now()
	result, err := r.pool.Exec(ctx, query,
		product.Name, product.Description, product.Price, product.Stock, product.UpdatedAt, product.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrProductNotFound
	}
	return nil
}

func (r *PostgresProductRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM products WHERE id = $1`
	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ports.ErrProductNotFound
	}
	return nil
}
