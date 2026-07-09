package repositories

import (
	"errors"

	"gorm.io/gorm"

	"ticket-system/internal/models"
)

// UserRepository is the data-access layer for User records.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository builds a UserRepository backed by the given DB
// connection.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new user row.
func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByEmail returns the user with the given email, or ErrNotFound.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

// FindByID returns the user with the given primary key, or ErrNotFound.
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
