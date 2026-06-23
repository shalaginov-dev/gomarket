package domain

import "time"

type Order struct {
	OrderID    int       `json:"order_id"`
	UserID     int       `json:"user_id"`
	Status     string    `json:"status"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderItem struct {
	OrderItemID int     `json:"order_item_id"`
	OrderID     int     `json:"order_id"`
	ProductID   int     `json:"product_id"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
}
