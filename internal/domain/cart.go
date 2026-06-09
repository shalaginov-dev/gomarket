package domain

type CartItem struct {
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
}

type Cart struct {
	UserID int        `json:"user_id"`
	Items  []CartItem `json:"items"`
}
