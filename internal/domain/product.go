package domain

import "time"

type Product struct {
	ProductID     int       `json:"product_id"`
	ProductName   string    `json:"product_name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}
