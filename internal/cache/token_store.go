package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type TokenStore struct {
	client *redis.Client
}

func NewTokenStore(client *redis.Client) *TokenStore {
	return &TokenStore{
		client: client,
	}
}

func (s *TokenStore) SaveRefreshToken(ctx context.Context, userID int, token string, ttl time.Duration) error {
	key := "refresh:" + strconv.Itoa(userID)
	if err := s.client.Set(ctx, key, token, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save refresh token: %w", err)
	}
	return nil
}
func (s *TokenStore) DeleteRefreshToken(ctx context.Context, userID int) error {
	key := "refresh:" + strconv.Itoa(userID)
	if err := s.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}
	return nil
}

func (s *TokenStore) GetRefreshToken(ctx context.Context, userID int) (string, error) {
	key := "refresh:" + strconv.Itoa(userID)
	token, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get refresh token: %w", err)
	}
	return token, nil
}
