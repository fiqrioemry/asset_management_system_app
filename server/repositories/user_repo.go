package repositories

import (
	"errors"

	"github.com/fiqrioemry/asset_management_system_app/server/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(data *models.User) error
	Update(data *models.User) error
	Delete(data *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(data *models.User) error {
	return r.db.Create(data).Error
}

func (r *userRepository) Update(data *models.User) error {
	return r.db.Save(data).Error
}

func (r *userRepository) Delete(data *models.User) error {
	if err := r.db.Delete(data).Error; err != nil {
		return err
	}
	return nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}
