package service

import (
	"context"
	"fmt"
	"gomarket/internal/cache"
	"gomarket/internal/domain"
	"gomarket/internal/repository"
)

type OrderService struct {
	orderRepo   *repository.OrderRepository
	productRepo *repository.ProductRepository
	cartStore   *cache.CartStore
}

func NewOrderService(orderRepo *repository.OrderRepository, productRepo *repository.ProductRepository, cartStore *cache.CartStore) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		cartStore:   cartStore,
	}
}

func (s *OrderService) Checkout(ctx context.Context, userID int) (*domain.Order, error) {
	cart, err := s.cartStore.GetItem(ctx, userID)
	if err != nil {
		return nil, err
	}
	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}
	var orderItem []domain.OrderItemInput
	for _, i := range cart.Items {
		product, err := s.productRepo.GetByID(ctx, i.ProductID)
		if err != nil {
			return nil, err
		}
		if product.StockQuantity >= i.Quantity {
			orderItem = append(orderItem, domain.OrderItemInput{
				ProductID: i.ProductID,
				Quantity:  i.Quantity,
				Price:     product.Price,
			})
		} else {
			return nil, fmt.Errorf("not enough stock for product %d", i.ProductID)
		}

	}
	order, err := s.orderRepo.Create(ctx, userID, orderItem)
	if err != nil {
		return nil, err
	}
	if err = s.cartStore.CleanCart(ctx, userID); err != nil {
		return nil, err
	}
	return order, nil
}

//func (s *OrderService) Create(ctx context.Context, userID, productID int) (*domain.Order, error) {}

//func (s *OrderService) Delete(ctx context.Context, userID int) error {}

func (s *OrderService) GetAll(ctx context.Context, userID int) ([]domain.Order, error) {
	orders, err := s.orderRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *OrderService) GetByID(ctx context.Context, userID, orderID int) (*domain.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, userID, orderID)
	if err != nil {
		return nil, err
	}
	return order, nil
}
