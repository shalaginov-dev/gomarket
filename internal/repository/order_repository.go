package repository

import (
	"context"
	"fmt"
	"gomarket/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	pool *pgxpool.Pool
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{pool: pool}
}
func (r *OrderRepository) Create(ctx context.Context, userID int, items []domain.OrderItemInput) (*domain.Order, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO orders (user_id, status, total_price)
		VALUES ($1, 'pending', $2)
		RETURNING order_id, user_id, status, total_price, created_at`
	var order domain.Order
	var totalPrice float64
	for _, item := range items {
		totalPrice += item.Price * float64(item.Quantity)
	}
	err = tx.QueryRow(ctx, query, userID, totalPrice).Scan(&order.OrderID, &order.UserID, &order.Status, &order.TotalPrice, &order.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range items {
		query := `
		INSERT INTO order_items (order_id, product_id, quantity, price)
		VALUES ($1, $2, $3, $4)`
		_, err = tx.Exec(ctx, query, order.OrderID, item.ProductID, item.Quantity, item.Price)
		if err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}
		query = `
		UPDATE products 
		SET stock_quantity = stock_quantity - $1
		WHERE product_id = $2`
		_, err = tx.Exec(ctx, query, item.Quantity, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("failed to change product stock quantity: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to finish transaction: %w", err)
	}
	return &order, nil
}
