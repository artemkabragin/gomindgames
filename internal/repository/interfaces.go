package repository

import (
	"mindgames/internal/domain"
)

type IUserRepository interface {
	Create(user *domain.User) error
	GetByUsername(username string) (*domain.User, error)
}

type ITokenRepository interface {
	Create(refreshToken *domain.RefreshToken) error
	GetRefreshByValue(refreshTokenValue string) (*domain.RefreshToken, error)
}
