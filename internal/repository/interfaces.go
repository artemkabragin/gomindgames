package repository

import (
	"mindgames/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
}

type TokenRepository interface {
	Create(refreshToken *domain.RefreshToken) error
	GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error)
}
