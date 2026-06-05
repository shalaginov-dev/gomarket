package service

import (
	"context"
	"gomarket/internal/domain"
	"gomarket/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) Create(ctx context.Context, p *domain.Product) (*domain.Product, error) {
	product, err := s.productRepo.Create(ctx, p)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) GetAll(ctx context.Context, limit int) ([]domain.Product, error) {
	products, err := s.productRepo.GetAll(ctx, limit)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s *ProductService) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Update(ctx context.Context, p *domain.Product, id int) (*domain.Product, error) {
	product, err := s.productRepo.Update(ctx, p, id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Delete(ctx context.Context, id int) error {
	if err := s.productRepo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
