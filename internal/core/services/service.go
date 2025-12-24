package services

import (
	"errors"
	"time"

	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo ports.UserRepository
}

func NewBackEndService(
	userRepo ports.UserRepository,
) ports.Service {
	return &service{
		userRepo,
	}
}

func (s *service) Register(name, email, password string) error {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := &domain.User{
		Name:     name,
		Email:    email,
		Password: string(hashed),
	}
	return s.userRepo.Create(user)
}

func (s *service) Authenticate(email, password string) (*domain.User, error) {
	filter := map[string]interface{}{"email": email}
	user, err := s.userRepo.FindOne(filter)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *service) CreateUser(user *domain.User) error {
	user.CreatedAt = time.Now()
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	return s.userRepo.Create(user)
}

func (s *service) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}
