package repository

import (
	"fmt"
	"mindgames/internal/domain"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func UserRepo(db *gorm.DB) IUserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (r UserRepositoryImpl) Create(user *domain.User) error {
	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	if err := r.db.Create(&user).Error; err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (r UserRepositoryImpl) GetByUsername(username string) (*domain.User, error) {
	var user domain.User

	if err := r.db.First(&user, "username = ?", username).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}
