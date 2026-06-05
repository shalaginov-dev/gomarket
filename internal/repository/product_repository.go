package repository

import (
	"context"
	"fmt"
	"gomarket/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
}

func NewProductRepository(pool *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{pool: pool}
}

func (r *ProductRepository) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (product_name, description, price, stock_quantity, created_at)
		VALUES ($1, $2, $3, $4, NOW())
		RETURNING product_id, product_name, description, price, stock_quantity, created_at`

	var product domain.Product
	err := r.pool.QueryRow(ctx, query, p.ProductName, p.Description, p.Price, p.StockQuantity).
		Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return &product, nil
}

func (r *ProductRepository) GetAll(ctx context.Context, limit int) ([]domain.Product, error) {
	query := `
        SELECT product_id, product_name, description, price, stock_quantity, created_at
        FROM products
        ORDER BY created_at ASC
        LIMIT $1`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close() // важно закрывать rows!

	var products []domain.Product

	for rows.Next() {
		var p domain.Product

		err := rows.Scan(
			&p.ProductID,
			&p.ProductName,
			&p.Description,
			&p.Price,
			&p.StockQuantity,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		products = append(products, p)
	}

	return products, nil
}

func (r *ProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	query := `
        SELECT product_id, product_name, description, price, stock_quantity, created_at
        FROM products
        WHERE product_id = $1
        LIMIT 1`

	var product domain.Product
	err := r.pool.QueryRow(ctx, query, id).
		Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return &product, nil
}

func (r *ProductRepository) Update(ctx context.Context, p *domain.Product, id int) (*domain.Product, error) {
	query := `
		UPDATE products
		SET product_name = $1,
		    description = $2,
		    price = $3,
		    stock_quantity = $4
		WHERE product_id = $5
		RETURNING product_id, product_name, description, price, stock_quantity, created_at`

	var product domain.Product

	err := r.pool.QueryRow(ctx, query, p.ProductName, p.Description, p.Price, p.StockQuantity, id).
		Scan(&product.ProductID, &product.ProductName, &product.Description, &product.Price, &product.StockQuantity, &product.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	return &product, nil
}
func (r *ProductRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM products WHERE product_id = $1`

	tag, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	affected := tag.RowsAffected()

	if affected == 0 {
		return fmt.Errorf("product with id %d not found", id)
	}
	return nil
}
