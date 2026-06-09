package service

import (
	"context"
	"gomarket/internal/cache"
	"gomarket/internal/domain"
	"gomarket/internal/repository"
)

type CartService struct {
	productRepo *repository.ProductRepository
	cartStore   *cache.CartStore
}

func NewCartService(productRepo *repository.ProductRepository, cartStore *cache.CartStore) *CartService {
	return &CartService{
		productRepo: productRepo,
		cartStore:   cartStore,
	}
}

func (s *CartService) GetItem(ctx context.Context, userID int) (*domain.Cart, error) {
	cart, err := s.cartStore.GetItem(ctx, userID)
	if err != nil {
		return nil, err
	}
	return cart, nil
}

func (s *CartService) AddItem(ctx context.Context, userID int, productID int, quantity int) error {
	_, err := s.productRepo.GetByID(ctx, productID)
	if err != nil {
		return err
	}
	err = s.cartStore.AddItem(ctx, userID, productID, quantity)
	if err != nil {
		return err
	}
	return nil
}

func (s *CartService) RemoveItem(ctx context.Context, userID int, productID int) error {
	err := s.cartStore.RemoveItem(ctx, userID, productID)
	if err != nil {
		return err
	}
	return nil
}
