package repositories

import (
	"errors"
	"time"

	"github.com/kanta/backend-challenge/internal/adapters/repositories/models"
	"github.com/kanta/backend-challenge/internal/core/domain"
	"github.com/kanta/backend-challenge/internal/core/ports"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) ports.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(user *domain.User) error {
	user.CreatedAt = time.Now()

	m := models.ToUserModels(user)

	result := r.db.Create(m)
	return result.Error
}

func (r *userRepository) FindOne(filter map[string]interface{}) (*domain.User, error) {
	var m models.User

	result := r.db.Where(filter).First(&m)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return models.ToUserDomain(&m), nil
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	var m models.User

	result := r.db.First(&m, "id = ?", id)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New("user not found")
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return models.ToUserDomain(&m), nil
}
