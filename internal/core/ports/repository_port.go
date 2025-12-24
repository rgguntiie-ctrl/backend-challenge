package ports

import (
	"github.com/kanta/backend-challenge/internal/core/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindOne(filter map[string]interface{}) (*domain.User, error)
	FindByID(id string) (*domain.User, error)
}
