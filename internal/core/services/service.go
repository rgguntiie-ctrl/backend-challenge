package services

import (
	"errors"
	"log"
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

func (s *service) ListUsers(page, limit int) ([]domain.User, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	skip := (page - 1) * limit
	return s.userRepo.ListPaginated(skip, limit)
}

func (s *service) UpdateUser(id, name, email string) error {
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	user.Name = name
	user.Email = email
	return s.userRepo.Update(user)
}

func (s *service) DeleteUser(id string) error {
	return s.userRepo.Delete(id)
}

func (s *service) RunUserCountLogger(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				count, err := s.userRepo.Count()
				if err != nil {
					log.Println("Error counting users:", err)
					continue
				}
				log.Println("Total users in DB:", count)
			}
		}
	}()
}
