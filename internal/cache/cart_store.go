package cache

import (
	"context"
	"fmt"
	"gomarket/internal/domain"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type CartStore struct {
	client *redis.Client
}

func NewCartStore(client *redis.Client) *CartStore {
	return &CartStore{
		client: client,
	}
}

func (s *CartStore) GetItem(ctx context.Context, userID int) (*domain.Cart, error) {
	key := "cart:" + strconv.Itoa(userID)
	cart, err := s.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var items []domain.CartItem
	for k, c := range cart {
		productID, _ := strconv.Atoi(k)
		quantity, _ := strconv.Atoi(c)
		items = append(items, domain.CartItem{ProductID: productID, Quantity: quantity})
	}
	return &domain.Cart{UserID: userID, Items: items}, nil
}

func (s *CartStore) AddItem(ctx context.Context, userID, productID, quantity int) error {
	key := "cart:" + strconv.Itoa(userID)
	if err := s.client.HSet(ctx, key, productID, quantity).Err(); err != nil {
		return err
	}
	if err := s.client.Expire(ctx, key, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("failed to add expire time: %w", err)
	}
	return nil
}

func (s *CartStore) RemoveItem(ctx context.Context, userID, productID int) error {
	key := "cart:" + strconv.Itoa(userID)
	result, err := s.client.HDel(ctx, key, strconv.Itoa(productID)).Result()
	if err != nil {
		return err
	}
	if result == 0 {
		return fmt.Errorf("product %d not found in cart", productID)
	}
	return nil
}
