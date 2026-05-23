package service

import "golang.org/x/crypto/bcrypt"

type PasswordService struct {
	cost int
}

func NewPasswordServise() *PasswordService {
	return &PasswordService{cost: bcrypt.DefaultCost}
}

func (s *PasswordService) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (s *PasswordService) CheckPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
