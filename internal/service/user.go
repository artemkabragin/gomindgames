package service

import (
	"errors"
	"fmt"
	"mindgames/internal/domain"
	"mindgames/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	repo repository.IUserRepository
}

func UserService(repo repository.IUserRepository) *UserServiceImpl {
	return &UserServiceImpl{
		repo,
	}
}

func (s *UserServiceImpl) Create(user *domain.User, password string) error {
	err := s.validateRegister(user, password)
	if err != nil {
		return fmt.Errorf("error validate register: %w", err)
	}

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	user.PasswordHash = string(hashedPassword)

	if err := s.repo.Create(user); err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	return nil
}

func (s *UserServiceImpl) GetByUsername(username string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	return user, err
}

// Private Methods

func (s *UserServiceImpl) validateRegister(user *domain.User, password string) error {
	if user.Username == "" {
		return errors.New("username is required")
	}

	if password == "" {
		return errors.New("password is incorrect")
	}

	return nil
}
