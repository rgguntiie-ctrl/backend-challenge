package ports

import (
	"github.com/kanta/backend-challenge/internal/core/domain"
)

type Service interface {
	Register(name, email, password string) error
	Authenticate(email, password string) (*domain.User, error)
	CreateUser(user *domain.User) error
	GetUserByID(id string) (*domain.User, error)
}
