package service

import (
	"context"
	"fmt"
	"gomarket/internal/domain"
	"gomarket/internal/repository"
)

type UserService struct {
	userRepo        *repository.UserRepository
	passwordService *PasswordService
}

func NewUserService(userRepo *repository.UserRepository, pwService *PasswordService) *UserService {
	return &UserService{
		userRepo:        userRepo,
		passwordService: pwService,
	}
}
func (s *UserService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Хэшируем пароль
	hashedPassword, err := s.passwordService.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Создаём пользователя с уже хэшированным паролем
	user, err := s.userRepo.Create(ctx, email, hashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (*domain.User, error) {
	// Ищем пользователя по почте
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}
	// Проверям пароль
	if s.passwordService.CheckPassword(password, user.PasswordHash) {
		return user, nil
	}
	return nil, fmt.Errorf("invalid credentials")
}
